package models

import "gorm.io/gorm"

type ImageModel struct{
	gorm.Model
	
	// sobre o arquivo
	FileName    string `gorm:"not null"` // Nome original (ex: ferias.jpg)
    StorageKey  string `gorm:"uniqueIndex;not null"` // Caminho no bucket 

    // meta dados
    ContentType string
    FileSize int64
    
    // acessibilidade
    AltText string `gorm:"type:varchar(255)"`
	
}

