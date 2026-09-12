package audit

import (
	"github.com/r3labs/diff/v3"
	"k8s.io/apimachinery/pkg/util/json"
)

type Change struct {
	Path string      `json:"path"`
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

func ComputeDiff(before, after *string) ([]Change, error) {
	if before == nil || after == nil {
		return nil, nil
	}

	var a, b interface{}
	if err := json.Unmarshal([]byte(*before), &b); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(*after), &a); err != nil {
		return nil, err
	}

	changeLog, err := diff.Diff(b, a)
	if err != nil {
		return nil, err
	}

	changes := make([]Change, 0, len(changeLog))
	for _, change := range changeLog {
		changes = append(changes, Change{
			Path: joinPath(change.Path),
			From: change.From,
			To:   change.To,
		})
	}
	return changes, nil
}

func joinPath(s []string) string {
	out := ""
	for i, v := range s {
		if i > 0 {
			out += "."
		}
		out += v
	}
	return out
}
