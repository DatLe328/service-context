package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	sctx "github.com/DatLe328/service-context"
	"github.com/DatLe328/service-context/component/gormc"
)

func main() {
	serviceCtx := sctx.NewServiceContext(
		sctx.WithName("demo-gorm"),
		sctx.WithComponent(gormc.NewGormDB("sqlite", "sqlite")),
	)

	if err := serviceCtx.Load(); err != nil {
		log.Fatal("Failed to load service context:", err)
	}

	log.Print("Init complete")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutting down...")
	if err := serviceCtx.Stop(); err != nil {
		log.Fatal("Failed to stop service context:", err)
	}
}
