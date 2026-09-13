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
	"codeberg.org/MrTomate/web/internal/adapters"
	"codeberg.org/MrTomate/web/internal/core"
)

func main() {
	adapter := adapters.NewInMemoryUserRepository()
	service := core.NewUserService(adapter)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.HandleHello())
	mux.HandleFunc("/me/{Name}", handler.HandleMe())
	mux.HandleFunc("/", handler.HandleNotFound())

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.ico")
	})
	mux.HandleFunc("POST /users/create", handler.HandleCreateUser(service))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      Logger(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		log.Println("Servidor iniciado na porta :8080...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		log.Fatalf("Erro fatal no servidor: %v", err)
	case <-ctx.Done():
		log.Println("Sinal de desligamento recebido. Iniciando graceful shutdown...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Erro durante o desligamento forçado do servidor HTTP: %v", err)
		if err := srv.Close(); err != nil {
			log.Fatalf("Erro ao fechar conexões ativas: %v", err)
		}
	}

	// Espaço reservado para liberação de outros recursos (ex: adapter.Close())
	log.Println("Servidor finalizado com sucesso.")
}

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
