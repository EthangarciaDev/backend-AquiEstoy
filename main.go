package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Servidor con Gin funcionando 🚀",
        })
    })

    r.Run(":8080") // Escucha en http://localhost:8080
}