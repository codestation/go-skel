// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package s3

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

type FakeWriterAt struct {
	w io.Writer
}

func (fw FakeWriterAt) WriteAt(p []byte, _ int64) (n int, err error) {
	return fw.w.Write(p)
}

func TestMemoryClient(t *testing.T) {
	client := NewMemoryClient()

	ctx := context.Background()
	content := "Hello, World!"

	key, err := client.Upload(ctx, "foo.txt", strings.NewReader(content))
	if err != nil {
		t.Fatal(err)
	}

	// create a writerAt that holds a buffer

	buf := bytes.Buffer{}
	_, err = client.Download(ctx, key, FakeWriterAt{&buf})
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != content {
		t.Fatalf("expected %q, got %q", content, buf.String())
	}
}
