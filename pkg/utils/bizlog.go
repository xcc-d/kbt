package utils

import (
	"context"

	"go.uber.org/zap"
)

const (
	UserContextKey      = "username"
	ClientIPContextKey  = "client_ip"
	UserAgentContextKey = "user_agent"
)

type BizLogger struct {
	logger       *zap.Logger
	user         string
	action       string
	resource     string
	resourceName string
	extra        []zap.Field
}

func Biz(ctx context.Context, action, resource, resourceName string) *BizLogger {
	user, _ := ctx.Value(UserContextKey).(string)
	return &BizLogger{
		logger:       L(),
		user:         user,
		action:       action,
		resource:     resource,
		resourceName: resourceName,
	}
}

func (b *BizLogger) WithExtra(extra ...zap.Field) *BizLogger {
	b.extra = append(b.extra, extra...)
	return b
}

func (b *BizLogger) fields() []zap.Field {
	fields := make([]zap.Field, 0, 4+len(b.extra))
	fields = append(fields,
		zap.String("user", b.user),
		zap.String("action", b.action),
		zap.String("resource", b.resource),
		zap.String("resource_name", b.resourceName),
	)
	fields = append(fields, b.extra...)
	return fields
}

func (b *BizLogger) Success() {
	b.logger.Info("biz", b.fields()...)
}

func (b *BizLogger) Fail(err error) {
	fields := append(b.fields(), zap.Error(err))
	b.logger.Error("biz", fields...)
}

func UserFromCtx(ctx context.Context) string {
	user, _ := ctx.Value(UserContextKey).(string)
	return user
}

func ClientIPFromCtx(ctx context.Context) string {
	ip, _ := ctx.Value(ClientIPContextKey).(string)
	return ip
}

func UserAgentFromCtx(ctx context.Context) string {
	ua, _ := ctx.Value(UserAgentContextKey).(string)
	return ua
}
