package audit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"kbt/internal/model"

	"k8s.io/apimachinery/pkg/util/json"
)

const GenesisHash = "0000000000000000000000000000000000000000000000000000000000000000"

type hashPayload struct {
	CreatedAt    string `json:"created_at"`
	Operator     string `json:"operator"`
	Action       string `json:"action"`
	Resource     string `json:"resource"`
	ResourceName string `json:"resource_name"`
	Namespace    string `json:"namespace"`
	Result       string `json:"result"`
	ErrorMsg     string `json:"error_msg"`
	BeforeJSON   string `json:"before_json"`
	AfterJSON    string `json:"after_json"`
	DiffJSON     string `json:"diff_json"`
	ClientIP     string `json:"client_ip"`
	UserAgent    string `json:"user_agent"`
	TraceID      string `json:"trace_id"`
	PrevHash     string `json:"prev_hash"`
	KeyVersion   string `json:"key_version"`
}

func ComputeHash(log *model.AuditLog, secretKey []byte) (string, error) {
	value := hashPayload{
		CreatedAt:    log.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		Operator:     log.Operator,
		Action:       log.Action,
		Resource:     log.Resource,
		ResourceName: log.ResourceName,
		Namespace:    log.Namespace,
		Result:       log.Result,
		ErrorMsg:     log.ErrorMsg,
		BeforeJSON:   deref(log.BeforeJSON),
		AfterJSON:    deref(log.AfterJSON),
		DiffJSON:     deref(log.DiffJSON),
		ClientIP:     log.ClientIP,
		UserAgent:    log.UserAgent,
		TraceID:      log.TraceID,
		PrevHash:     log.PrevHash,
		KeyVersion:   log.KeyVersion,
	}

	marshal, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write(marshal)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
