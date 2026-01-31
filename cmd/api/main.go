package main

import (
	"context"
	"log"

	"github.com/Cannedsans/GopherGallery/internal/config"
	"github.com/Cannedsans/GopherGallery/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	config.InitAWS(context.TODO())

	godotenv.Load()

	// ---- INICIALIZAÇÃO DE RECURSOS GLOBAIS ----
	// 1. Inicializa o Banco de Dados
	config.StartBd()
	if config.DataBase == nil {
		log.Fatalf("Falha crítica ao inicializar o banco de dados.")
		return
	}
	//	log.Println("Banco de dados inicializado com sucesso.")

	// ------------------------------------------

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/files", handler.ListFilesHandler)
		api.POST("/Upfile", handler.SaveFileHandler)
		api.DELETE("/Delefile/:id", handler.DeleteFileHandler)

		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK"})
		})
	}

	log.Println("Servidor Gin iniciando em :8080")
	r.Run()
}
