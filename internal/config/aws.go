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

	// Configurações do cliente AWS
	conf, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					if endpoint != "" {
						return aws.Endpoint{
							URL:               endpoint,
							HostnameImmutable: true,
						}, nil
					}
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				}),
		),
	)
	if err != nil {
		return fmt.Errorf("falha ao carregar a configuração AWS: %w", err)
	}

	// Cria o cliente S3 e o armazena na variável global
	S3Client = s3.NewFromConfig(conf)

	createBucket(ctx)
	
	return nil
}

// GetS3Client retorna o cliente S3 globalmente inicializado.
func GetS3Client() (*s3.Client, error) {
	if S3Client == nil {
		return nil, fmt.Errorf("cliente S3 não foi inicializado. Chame InitAWS() na inicialização da aplicação")
	}
	return S3Client, nil
}

func createBucket(ctx context.Context) error {
	// Primeiro, listar todos os buckets para verificar se o bucket desejado já existe.
	listBucketsInput := &s3.ListBucketsInput{}
	listBucketsOutput, err := S3Client.ListBuckets(ctx, listBucketsInput)
	if err != nil {
		return fmt.Errorf("falha ao listar buckets: %w", err)
	}

	bucketExists := false
	for _, bucket := range listBucketsOutput.Buckets {
		if aws.ToString(bucket.Name) == Bockete {
			bucketExists = true
			break
		}
	}

	// Se o bucket não existir, crie-o.
	if !bucketExists {
		createBucketInput := &s3.CreateBucketInput{
			Bucket: aws.String(Bockete),
		}
		_, err := S3Client.CreateBucket(ctx, createBucketInput)
		if err != nil {
			return fmt.Errorf("falha ao criar bucket '%s': %w", Bockete, err)
		}
	} 
	return nil
}
