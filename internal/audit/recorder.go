package audit

import (
	"errors"
	"kbt/internal/model"
	"kbt/pkg/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/json"
)

type Entry struct {
	Operator     string
	Action       string
	Resource     string
	ResourceName string
	Namespace    string
	Result       string
	ErrorMsg     string
	Before       interface{}
	After        interface{}
	ClientIP     string
	UserAgent    string
	TraceID      string
	KeyVersion   string
}

type Recorder struct {
	DB        *gorm.DB
	Ch        chan *Entry
	SecretKey []byte
	LastHash  string
}

func NewRecorder(db *gorm.DB, secretKey []byte, bufferSize int) *Recorder {
	r := &Recorder{
		DB:        db,
		Ch:        make(chan *Entry, bufferSize),
		SecretKey: secretKey,
	}
	go r.loop()
	return r
}

func (r *Recorder) Record(e *Entry) {
	select {
	case r.Ch <- e:
	default:
		utils.L().Warn("record channel is full")
	}
}

func (r *Recorder) loop() {
	if err := r.loadLastHash(); err != nil {
		utils.L().Error("load last audit hash failed", zap.Error(err))
	}
	for e := range r.Ch {
		if err := r.write(e); err != nil {
			utils.L().Error("write audit log failed", zap.Error(err))
		}
	}
}

func (r *Recorder) loadLastHash() error {
	var last model.AuditLog
	err := r.DB.Order("id DESC").First(&last).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.LastHash = GenesisHash
		return nil
	}
	if err != nil {
		return err
	}
	r.LastHash = last.CurrHash
	return nil
}

func (r *Recorder) write(e *Entry) error {
	beforeJson := marshalPtr(e.Before)
	afterJson := marshalPtr(e.After)
	change, err := ComputeDiff(beforeJson, afterJson)
	if err != nil {
		utils.L().Warn("compute audit diff failed", zap.Error(err))
	}
	diffJson := marshalPtr(change)

	rec := &model.AuditLog{
		// 截断到毫秒，与 DB 列 datetime(3) 精度及 ComputeHash 的格式化保持一致，
		// 否则纳秒部分会被 MySQL 四舍五入，导致读出时重算的 hash 与写入时不一致。
		CreatedAt:    time.Now().UTC().Truncate(time.Millisecond),
		Operator:     e.Operator,
		Action:       e.Action,
		Resource:     e.Resource,
		ResourceName: e.ResourceName,
		Namespace:    e.Namespace,
		Result:       e.Result,
		ErrorMsg:     e.ErrorMsg,
		BeforeJSON:   beforeJson,
		AfterJSON:    afterJson,
		DiffJSON:     diffJson,
		ClientIP:     e.ClientIP,
		UserAgent:    e.UserAgent,
		TraceID:      e.TraceID,
		PrevHash:     r.LastHash,
		KeyVersion:   e.KeyVersion,
	}

	currHash, err := ComputeHash(rec, r.SecretKey)
	if err != nil {
		return err
	}
	rec.CurrHash = currHash

	if err := r.DB.Create(rec).Error; err != nil {
		return err
	}
	r.LastHash = rec.CurrHash
	return nil
}

func marshalPtr(v interface{}) *string {
	if v == nil {
		return nil
	}
	marshal, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	str := string(marshal)
	return &str
}
