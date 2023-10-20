// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package s3

import (
	"context"
	"io"
)

var _ Client = &MemoryClient{}

type MemoryClient struct {
	Files map[string][]byte
}

func NewMemoryClient() *MemoryClient {
	return &MemoryClient{
		Files: make(map[string][]byte),
	}
}

func (m MemoryClient) Upload(_ context.Context, key string, r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	m.Files[key] = data
	return key, nil
}

func (m MemoryClient) Download(_ context.Context, key string, w io.WriterAt) (int64, error) {
	data := m.Files[key]
	n, err := w.WriteAt(data, 0)
	if err != nil {
		return 0, err
	}

	return int64(n), nil
}

func (m MemoryClient) ListObjects(_ context.Context, _ ListOptions) (*ListOutput, error) {
	output := &ListOutput{}
	for key := range m.Files {
		output.Items = append(output.Items, ListItem{
			Key:  key,
			Size: int64(len(m.Files[key])),
		})
	}

	return output, nil
}
