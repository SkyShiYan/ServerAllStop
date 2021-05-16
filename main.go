package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// 启动服务
	go startServer(ctx)

	sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// signal.Notify(sigs)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		sig := <-sigs
		fmt.Println("收到CTRL+C signal信号，开始通知服务关闭", sig)
		cancel()
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	//程序将在此处等待，直到它预期信号（如Goroutine所示）
	//在“done”上发送一个值，然后退出。
	fmt.Println("等待取消信号")
	wg.Wait()
	fmt.Println("关闭")
}

func startServer(ctx context.Context) {
	fmt.Println("开始启动服务，监听8080端口")
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})
	go func(ctx context.Context, srv *http.Server) {
		select {
		case <-ctx.Done():
			fmt.Println("收到cancel信号准备关闭服务")
			if err := srv.Close(); err != nil {
				fmt.Println("服务关闭有异常", err.Error())
				panic(err)
			} else {
				fmt.Println("服务已关闭")
			}
		}
	}(ctx, srv)
	// 启动http server
	fmt.Println("启动服务")
	srv.ListenAndServe()
}
