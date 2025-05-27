package main

import (
	"busca-cpf/database"
	"busca-cpf/handlers"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var appVersion = "N/A"

func loadAppVersion() {
	versionBytes, err := os.ReadFile("VERSION") // Alterado de ioutil.ReadFile para os.ReadFile
	if err != nil {
		log.Printf("Aviso: Arquivo VERSION não encontrado. Usando '%s'. Certifique-se que o script de deploy o baixe.", appVersion)
		return
	}
	appVersion = strings.TrimSpace(string(versionBytes)) // Corrigido para strings.TrimSpace
	log.Printf("Versão da Aplicação: %s", appVersion)
}

func main() {
	godotenv.Load()  // Carrega .env
	loadAppVersion() // Carrega a versão da aplicação

	// Conexão com PostgreSQL
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatal("Falha ao conectar ao PostgreSQL:", err)
	}

	// Conexão com Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "",
		DB:       0,
	})

	// Configuração do Gin
	r := gin.Default()

	// Conexão com Elasticsearch
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_HOST")},
	})
	if err != nil {
		log.Fatal("Falha ao conectar ao Elasticsearch:", err)
	}

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://187.62.247.143", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rotas xxxxxxxxxxxxx
	r.GET("/api/search", handlers.SearchUsers(db, rdb, es))
	r.GET("/api/user/:id", handlers.GetUserDetails(db))
	r.GET("/api/surname-suggestions", handlers.GetSurnameSuggestions(db, rdb))
	r.GET("/api/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": appVersion})
	})

	// Inicia o servidor na porta especificada
	r.Run(":" + os.Getenv("PORT"))
}
