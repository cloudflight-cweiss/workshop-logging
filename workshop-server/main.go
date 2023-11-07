package main

import (
	"fmt"

	"log"
	"log/slog"

	gpmiddleware "github.com/carousell/gin-prometheus-middleware"
	"github.com/gin-gonic/gin"
)

const (
	ApplicationName = "workshop-server"
)

func main() {
	slog.Info(fmt.Sprintf("Running %s ...", ApplicationName))

	g := gin.Default()
	p := gpmiddleware.NewPrometheus("gin")
	p.Use(g)
	if err := addEndpointHandlers(g); err != nil {
		log.Fatal(err)
	}

	if err := g.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
