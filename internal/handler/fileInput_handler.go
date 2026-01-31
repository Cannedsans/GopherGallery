package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Cannedsans/GopherGallery/internal/config"
	"github.com/Cannedsans/GopherGallery/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func SaveFileHandler(c *gin.Context) {
	foto, header, err := c.Request.FormFile("file")
	client, err := config.GetS3Client()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao obter o cliente de armazenamento.",
		})
		return
	}
	defer foto.Close()

	ext := filepath.Ext(header.Filename)
	randS := make([]byte, 8)

	rand.Read(randS)
	hash := hex.EncodeToString(randS)

	randStr := hash

	storageKey := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), randStr, ext)

	upInput := &s3.PutObjectInput{
		Bucket:      aws.String(config.Bockete),
		Key:         aws.String(storageKey),
		Body:        foto,
		ContentType: aws.String(header.Header.Get("Content-type")),
	}

	_, err = client.PutObject(c.Request.Context(), upInput)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao salvar o arquivo no armazenamento.",
		})
		return
	}

	var fileSize int64 = 0
	if header.Size > 0 {
		fileSize = header.Size
	}

	imagem := models.ImageModel{
		FileName:    header.Filename,
		StorageKey:  storageKey,
		ContentType: header.Header.Get("Content-Type"),
		FileSize:    fileSize,
		AltText:     c.PostForm("alt_text"),
	}

	if err := config.DataBase.Create(&imagem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao registrar o arquivo no banco de dados.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Arquivo enviado e registrado com sucesso!",
		"fileName":   imagem.FileName,
		"storageKey": imagem.StorageKey,
		"url":        fmt.Sprintf("%s/%s/%s", os.Getenv("AWS_ENDPOINT"), config.Bockete, storageKey), // Exemplo de URL pública, pode precisar de ajuste
	})
}
