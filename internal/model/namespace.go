package model

type Namespace struct {
	Name string `json:"name" form:"name" binding:"required"`
}
