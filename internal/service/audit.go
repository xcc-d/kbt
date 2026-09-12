package service

import (
	"context"
	"kbt/internal/model"
	"time"

	"gorm.io/gorm"
)

type AuditService struct {
	DB *gorm.DB
}

type AuditQuery struct {
	Operator     string
	Resource     string
	ResourceName string
	Namespace    string
	Action       string
	Result       string
	StartTime    *time.Time
	EndTime      *time.Time
	Page         int
	PageSize     int
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{DB: db}
}

func (a *AuditService) Audit(ctx context.Context, q AuditQuery) ([]model.AuditLog, int64, error) {
	tx := a.DB.WithContext(ctx).Model(&model.AuditLog{})

	if q.Operator != "" {
		tx = tx.Where("operator = ?", q.Operator)
	}
	if q.Resource != "" {
		tx = tx.Where("resource = ?", q.Resource)
	}
	if q.ResourceName != "" {
		tx = tx.Where("resource_name = ?", q.ResourceName)
	}
	if q.Namespace != "" {
		tx = tx.Where("namespace = ?", q.Namespace)
	}
	if q.Action != "" {
		tx = tx.Where("action = ?", q.Action)
	}
	if q.Result != "" {
		tx = tx.Where("result = ?", q.Result)
	}
	if q.StartTime != nil {
		tx = tx.Where("created_at >= ?", *q.StartTime)
	}
	if q.EndTime != nil {
		tx = tx.Where("created_at <= ?", *q.EndTime)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 20
	}

	var logs []model.AuditLog
	err := tx.Order("id DESC").
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
