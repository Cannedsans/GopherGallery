package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

var (
	Bockete  string
	S3Client *s3.Client
)

// InitAWS configura o cliente S3 e garante que o bucket exista.
func InitAWS(ctx context.Context) error {
	_ = godotenv.Load()

	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	Bockete = os.Getenv("BUCKET_NAME")

	// 1. Carrega a configuração básica
	conf, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("falha ao carregar a configuração AWS: %w", err)
	}

	// 2. Cria o cliente S3 usando a sintaxe moderna para injetar o LocalStack
	S3Client = s3.NewFromConfig(conf, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true // Substitui o antigo HostnameImmutable
		}
	})

	// 3. Garante que o bucket exista
	return createBucket(ctx)
}

// GetS3Client retorna o cliente S3 globalmente inicializado.
func GetS3Client() (*s3.Client, error) {
	if S3Client == nil {
		return nil, fmt.Errorf("cliente S3 não foi inicializado")
	}
	return S3Client, nil
}

func createBucket(ctx context.Context) error {
	// Verifica se o bucket existe (Forma simplificada usando HeadBucket)
	_, err := S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(Bockete),
	})

	// Se der erro, assumimos que o bucket não existe e tentamos criar
	if err != nil {
		_, err := S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(Bockete),
		})
		if err != nil {
			return fmt.Errorf("falha ao criar bucket '%s': %w", Bockete, err)
		}
		fmt.Printf("✅ Bucket '%s' criado com sucesso!\n", Bockete)
	}
	
	return nil
}