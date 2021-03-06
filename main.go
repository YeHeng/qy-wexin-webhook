package main

import (
	"github.com/YeHeng/qy-wexin-webhook/handler"
	. "github.com/YeHeng/qy-wexin-webhook/util"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := app()
	configRoute(r)

	err := r.Run(":9091")
	if err != nil {
		Logger.Fatalf("Gin start fail. %v", err)
	}
}

func configRoute(engine *gin.Engine) {
	engine.POST("/alertmanager", handler.AlertManagerHandler())
	engine.POST("/grafana", handler.GrafanaManagerHandler())
}

func app() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		Logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}, gin.Recovery())
	return r
}
