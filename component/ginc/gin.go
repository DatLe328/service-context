package ginc

import (
	"flag"
	"fmt"

	sctx "github.com/DatLe328/service-context"
	"github.com/DatLe328/service-context/logger"
	"github.com/gin-gonic/gin"
)

const (
	defaultPort = 3000
	defaultMode = "debug"
)

type Config struct {
	port    int
	ginMode string
}

type ginEngine struct {
	*Config
	id     string
	logger logger.Logger
	router *gin.Engine
}

func NewGin(id string) *ginEngine {
	return &ginEngine{
		Config: new(Config),
		id:     id,
	}
}

func (g *ginEngine) ID() string {
	return g.id
}

func (g *ginEngine) Activate(serviceContext sctx.ServiceContext) error {
	// gin-mode > app-env
	env := serviceContext.EnvName()
	mode := gin.ReleaseMode

	if env == sctx.DevEnv {
		mode = gin.DebugMode
	}

	if g.ginMode != "" {
		switch g.ginMode {
		case gin.DebugMode, gin.ReleaseMode:
			mode = g.ginMode

		default:
			return fmt.Errorf("invalid gin mode: %s (allowed: debug | release)", g.ginMode)
		}
	}

	gin.SetMode(mode)

	g.logger = serviceContext.Logger(g.id)
	g.logger.Info("init engine...")
	g.router = gin.New()

	return nil
}

func (g *ginEngine) Stop() error {
	return nil
}

func (g *ginEngine) InitFlags() {
	flag.IntVar(&g.port, "gin-port", defaultPort, "gin server port. Default 3000")
	flag.StringVar(&g.ginMode, "gin-mode", defaultMode, "gin server (debug | release). Default debug")
}

func (g *ginEngine) GetPort() int {
	return g.port
}

func (g *ginEngine) GetRouter() *gin.Engine {
	return g.router
}
