package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codeberg.org/MrTomate/web/handler"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.HandleHello())
	mux.HandleFunc("/me/{Name}", handler.HandleMe())
	mux.HandleFunc("/", handler.HandleNotFound())

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.ico")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: Logger(mux),
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Servidor iniciador na porta :8080...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro Fatal no servidor: %v ", err)
		}
	}()
	<-stop

	log.Println("Sinal de desligamento.... Iniciando ShutDown!!!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro durante o desligamento forçado : %v", err)
	}
	log.Println("GoodBye!!")
}

// Um Logger simples
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("METHOD: %s | URL: %s | TIME: %s",
			r.Method,
			r.URL.Path,
			time.Since(start))
	})
}
