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

func httpServer(ctx context.Context) error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Path))
	}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/api", handler)
	server := http.Server{
		Addr:    ADDR,
		Handler: serverMux,
	}
	go func() {
		<-ctx.Done()
		server.Shutdown(context.TODO())
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
	case <-sigCh:
		return nil
	}
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return httpServer(ctx)
	})

	eg.Go(func() error {
		return signalHandle(ctx)
	})

	if err := eg.Wait(); err != nil {
		fmt.Errorf("server exception, info: %v\n", err)
	}
}
