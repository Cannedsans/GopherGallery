package handler

import (
	"net/http"

	"github.com/Cannedsans/GopherGallery/internal/config"
	"github.com/Cannedsans/GopherGallery/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func DeleteFileHandler(c *gin.Context) {
	fileKey := c.Param("id")
	file := models.ImageModel{}
	if fileKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Erro": "Parametro id faltando",
		})
		return
	}

	if err := config.DataBase.Where("storage_key = ?", fileKey).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Erro": "Arquivo não encontrado",
		})
		return
	}

	if err := config.DataBase.Delete(&file, file.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Erro": "Arquivo não encontrado",
		})
		return
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(config.Bockete),
		Key:    aws.String(fileKey),
	}

	_, err := config.S3Client.DeleteObject(c, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Erro": "erro ao deletar arquivo do bucket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "arquivo deletado com sucesso",
	})
}
