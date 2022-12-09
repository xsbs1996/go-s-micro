package scli

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type CliServer struct {
	mode        string
	port        string
	engine      *gin.Engine
	ginRegister CliRegisterFn
}

type CliRegisterFn func(*gin.Engine)

func NewCliServer(mode string, port string, fn CliRegisterFn) *CliServer {
	gin.SetMode(mode)
	cli := &CliServer{
		mode:        mode,
		port:        port,
		engine:      gin.Default(),
		ginRegister: fn,
	}

	return cli
}

func (c *CliServer) Start() error {
	err := c.engine.Run(fmt.Sprintf(":%s", c.port))
	if err != nil {
		return err
	}
	return nil
}
