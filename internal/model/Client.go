package model

import "kbt/internal/client"

type Client struct {
	K8sClient *client.K8sClient
}
