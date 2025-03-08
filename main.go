package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mailbox-api/api/router"
	"mailbox-api/config"
	"mailbox-api/db"
	"mailbox-api/logger"
	"mailbox-api/repository"
	"mailbox-api/service"
)

func main() {
	l := logger.NewLogger()
	l.Info("Starting mailbox API service")

	cfg, err := config.Load()
	if err != nil {
		l.Fatal("Failed to load configuration", "error", err)
	}

	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		l.Fatal("Failed to connect to database", "error", err)
	}
	defer dbConn.Close()

	mailboxRepo := repository.NewMailboxRepository(dbConn)
	departmentRepo := repository.NewDepartmentRepository(dbConn)

	mailboxService := service.NewMailboxService(mailboxRepo, departmentRepo)

	r := router.SetupRouter(cfg, l, mailboxService)

	srv := r.Start(cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("Server forced to shutdown", "error", err)
	}

	l.Info("Server exiting")
}
