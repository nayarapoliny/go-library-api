package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/nayarapoliny/go-library-api/internal/book"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var repo book.Repository
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL != "" {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Erro ao abrir ligação PostgreSQL: %v", err)
		}
		defer db.Close()

		pgRepo, err := book.NewPostgresRepository(db)
		if err != nil {
			log.Fatalf("Erro ao inicializar tabela no PostgreSQL: %v", err)
		}
		repo = pgRepo
		log.Println("🗄️ Conectado ao PostgreSQL com sucesso!")
	} else {
		repo = book.NewMemoryRepository()
		log.Println("⚡ A executar com Repositório Thread-Safe em Memória (defina DATABASE_URL para usar PostgreSQL)")
	}

	svc := book.NewService(repo)
	handler := book.NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Printf("📚 Servidor Library Management API ativo em http://localhost:%s/api/v1/books\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor HTTP: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("⏳ A encerrar o servidor de forma segura (Graceful Shutdown)...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
