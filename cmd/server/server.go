package main

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(ware)
	router.StaticFile("/", "./frontend/index.html")
	router.StaticFile("/wasm_exec.js", "./frontend/wasm_exec.js")
	router.StaticFile("/main.wasm", "./frontend/main.wasm")
	router.StaticFile("/app.js", "./frontend/app.js")
	router.StaticFile("/style.css", "./frontend/style.css")
	router.POST("/dev/print", func(g *gin.Context) {
		b, err := io.ReadAll(g.Request.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	})
	err := router.Run()
	fmt.Println(err)
}

func ware(g *gin.Context) {
	// TODO client side caching
	g.Header("Access-Control-Allow-Origin", "*")
	g.Header("Access-Control-Allow-Methods", "*")
	g.Header("Access-Control-Allow-Headers", "*")
	switch g.Request.Method {
	case "HEAD":
		g.AbortWithStatus(204)
		return
	case "OPTIONS":
		g.AbortWithStatus(204)
		return
	}
	g.Next()
}
