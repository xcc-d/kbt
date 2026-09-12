package model

import (
	"kbt/internal/client"

	"gorm.io/gorm"
)

type Client struct {
	K8sClient *client.K8sClient
	DB        *gorm.DB
}
