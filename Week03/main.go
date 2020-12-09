package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	ADDR = ("127.0.0.1:8080")
)

func httpServer() error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Path))
	}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/api", handler)
	server := http.Server{
		Addr:    ADDR,
		Handler: serverMux,
	}
	defer func() {
		if err := recover(); err != nil {
			server.Close()
			fmt.Errorf("http server recover info:%v\n", err)
		}
	}()
	return server.ListenAndServe()
}

//监听SIGNAL信号
func signalHandle(ctx context.Context) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigCh)
	select {
	case <-ctx.Done():
		return fmt.Errorf("http quit")
	case err := <-sigCh:
		return fmt.Errorf("syscall signal info: %v\n", err)
	}
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return httpServer()
	})

	eg.Go(func() error {
		return signalHandle(ctx)
	})

	if err := eg.Wait(); err != nil {
		fmt.Errorf("server exception, info: %v\n", err)
	}
}
