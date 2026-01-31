package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Cannedsans/GopherGallery/internal/config"
	"github.com/Cannedsans/GopherGallery/internal/models"
	"github.com/gin-gonic/gin"
)

// DTO para enviar uma resposta limpa para o front-end
type FileResponse struct {
	Name string `json:"name"`
	Size int64  `json:"size_bytes"`
	Url  string `json:"url"`
	Key string `json:"fileKey"`
}

func ListFilesHandler(c *gin.Context) {
	var imagens []models.ImageModel

	config.DataBase.Find(&imagens)

	out := make([]FileResponse, len(imagens))

	for i, img := range imagens {
		out[i] = FileResponse{
			Name: img.FileName,
			Size: img.FileSize,
			Url:  fmt.Sprintf("%s/%s/%s", os.Getenv("AWS_ENDPOINT"), config.Bockete, img.StorageKey),
			Key: img.StorageKey,
		}
	}
	// 5. Retorna o JSON para o cliente (Site/API)
	c.JSON(http.StatusOK, out)
}
