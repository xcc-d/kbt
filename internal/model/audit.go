package model

import "time"

type AuditLog struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt    time.Time `gorm:"index;not null" json:"created_at"`
	Operator     string    `gorm:"size:128;not null;index:idx_operator_time,priority:1" json:"operator"`
	Action       string    `gorm:"size:32;not null" json:"action"`
	Resource     string    `gorm:"size:64;not null;index:idx_resource,priority:1" json:"resource"`
	ResourceName string    `gorm:"size:255;not null;index:idx_resource,priority:2" json:"resource_name"`
	Namespace    string    `gorm:"size:255;not null;default:''" json:"namespace"`
	Result       string    `gorm:"size:16;not null" json:"result"`
	ErrorMsg     string    `gorm:"size:1024;not null;default:''" json:"error_msg"`
	BeforeJSON   *string   `gorm:"type:longtext" json:"before_json,omitempty"`
	AfterJSON    *string   `gorm:"type:longtext" json:"after_json,omitempty"`
	DiffJSON     *string   `gorm:"type:longtext" json:"diff_json,omitempty"`
	ClientIP     string    `gorm:"size:64;not null;default:''" json:"client_ip"`
	UserAgent    string    `gorm:"size:512;not null;default:''" json:"user_agent"`
	TraceID      string    `gorm:"size:64;not null;default:''" json:"trace_id"`
	PrevHash     string    `gorm:"size:64;not null" json:"prev_hash"`
	CurrHash     string    `gorm:"size:64;not null" json:"curr_hash"`
	KeyVersion   string    `gorm:"size:64;not null;default:''" json:"key_version"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}
