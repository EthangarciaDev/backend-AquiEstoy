package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase establece la conexión con PostgreSQL de AWS RDS
func ConnectDatabase() {
	var err error

	log.Println("🔧 Iniciando configuración de base de datos...")

	// Obtener variables de entorno para la conexión a RDS
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Debug: mostrar variables (sin password)
	log.Printf("📋 Variables de entorno cargadas:")
	log.Printf("   DB_HOST: %s", host)
	log.Printf("   DB_PORT: %s", port)
	log.Printf("   DB_USER: %s", user)
	log.Printf("   DB_NAME: %s", dbname)
	log.Printf("   DB_SSLMODE: %s", sslmode)
	log.Printf("   DB_PASSWORD: %s", func() string {
		if password == "" {
			return "(vacío)"
		}
		return "(configurado - " + string(password[0]) + "***)"
	}())

	// Validar que las variables de entorno estén configuradas
	if host == "" {
		log.Fatal("❌ DB_HOST no está configurado")
	}
	if user == "" {
		log.Fatal("❌ DB_USER no está configurado")
	}
	if password == "" {
		log.Fatal("❌ DB_PASSWORD no está configurado")
	}
	if dbname == "" {
		log.Fatal("❌ DB_NAME no está configurado")
	}

	// Valores por defecto
	if port == "" {
		port = "5432"
		log.Println("🔧 Usando puerto por defecto: 5432")
	}
	if sslmode == "" {
		sslmode = "require"
		log.Println("🔧 Usando SSL mode por defecto: require")
	}

	// Crear el DSN (Data Source Name) con configuraciones adicionales
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC connect_timeout=30",
		host, port, user, password, dbname, sslmode)

	log.Printf("🔗 Intentando conectar a PostgreSQL RDS: %s:%s/%s", host, port, dbname)
	log.Println("⏳ Esto puede tomar hasta 30 segundos...")

	// Configurar GORM con timeout y logging
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Solo warnings para no saturar logs
	}

	// Conectar a la base de datos
	log.Println("📡 Estableciendo conexión...")
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Printf("❌ Error conectando a la base de datos: %v", err)
		log.Println("")
		log.Println("💡 Guía de troubleshooting:")
		log.Println("   1. Security Group debe permitir puerto 5432 desde tu IP")
		log.Println("   2. La instancia RDS debe ser 'Publicly accessible'")
		log.Println("   3. Verificar credenciales de usuario/password")
		log.Println("   4. Confirmar que la base de datos existe")
		log.Println("")
		log.Println("🧪 Para probar conectividad básica, ejecuta:")
		log.Printf("   nc -zv %s %s", host, port)
		log.Println("")
		log.Fatal("🚫 No se pudo establecer conexión con PostgreSQL RDS")
	}

	// Configurar el pool de conexiones
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Error obteniendo instancia de base de datos:", err)
	}

	// Configurar parámetros del pool de conexiones
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Probar la conexión
	log.Println("🏓 Realizando ping a la base de datos...")
	if err = sqlDB.Ping(); err != nil {
		log.Printf("❌ Error haciendo ping a la base de datos: %v", err)
		log.Fatal("🚫 No se pudo verificar la conexión a PostgreSQL RDS")
	}

	log.Println("✅ Conexión exitosa a PostgreSQL RDS")
	log.Println("🎉 Base de datos configurada correctamente")
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}