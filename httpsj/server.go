package httpsj

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xsbs1996/go-s-micro/httpsj/handler"
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
	// 添加全局中间件
	c.engine.Use(
		handler.TracingHandler(),
		handler.TracingLog(),
		handler.BreakerHandler(),
	)

	c.ginRegister(c.engine)
	err := c.engine.Run(fmt.Sprintf(":%s", c.port))
	if err != nil {
		return err
	}
	return nil
}
