// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ctxlog

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNewContextHandler(t *testing.T) {
	type want struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields []slog.Attr
		want   []want
	}{
		{
			name:   "no attributes",
			fields: []slog.Attr{},
			want:   []want{{"request_id", ""}},
		},
		{
			name:   "single attribute",
			fields: []slog.Attr{slog.String("request_id", "123")},
			want:   []want{{"request_id", "123"}},
		},
		{
			name:   "multiple attribute",
			fields: []slog.Attr{slog.String("request_id", "123"), slog.Int64("other", 1)},
			want:   []want{{"request_id", "123"}, {"other", float64(1)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			handler := NewHandler(slog.NewJSONHandler(&buf, nil))
			log := slog.New(handler)
			ctx := AddContextAttrs(context.Background(), tt.fields...)
			log.InfoContext(ctx, "message", "someKey", "someValue")
			actual := make(map[string]any)
			if err := json.Unmarshal(buf.Bytes(), &actual); err != nil {
				t.Errorf("failed to unmarshal json: %v", err)
			}
			for _, field := range tt.want {
				if actual[field.key] != field.value {
					t.Errorf("expected %s to be %s, got %v", field.key, field.value, actual[field.key])
				}
			}
		})
	}
}
