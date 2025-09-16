package main

import (
	"log"

	"github.com/EthangarciaDev/backend-AquiEstoy/internal/config"
	"github.com/EthangarciaDev/backend-AquiEstoy/internal/handlers"
	"github.com/EthangarciaDev/backend-AquiEstoy/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables del sistema")
	}

	// Conectar a la base de datos
	config.ConnectDatabase()

	// Inicializar servicios
	userService := services.NewUserService()
	authService := services.NewAuthService(userService)

	// Ejecutar migraciones
	if err := userService.MigrateDB(); err != nil {
		log.Fatal("Error ejecutando migraciones:", err)
	}

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService, userService)

	// Configurar Gin
	r := gin.Default()

	// Ruta de health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Servidor AquiEstoy Backend funcionando 🚀",
			"status":  "ok",
		})
	})

	// Rutas de autenticación (públicas)
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Rutas protegidas (requieren autenticación)
	protected := r.Group("/api")
	protected.Use(authHandler.AuthMiddleware())
	{
		protected.GET("/profile", authHandler.GetProfile)
		// Aquí puedes agregar más rutas protegidas según necesites
	}

	// Iniciar servidor
	log.Println("Servidor iniciado en puerto 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}