// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cursor

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestDecoder(t *testing.T) {
	suite.Run(t, &decoderSuite{})
}

type decoderSuite struct {
	suite.Suite
}

/* decode */

func (s *decoderSuite) TestDecodeKeyNotMatchedModel() {
	_, err := NewDecoder([]DecoderField{{Key: "Key"}}).Decode("cursor", struct{ ID string }{})
	s.Equal(ErrInvalidModel, err)
}

func (s *decoderSuite) TestDecodeNonStructModel() {
	_, err := NewDecoder([]DecoderField{{Key: "Key"}}).Decode("cursor", 123)
	s.Equal(ErrInvalidModel, err)
}

func (s *decoderSuite) TestDecodeInvalidCursorFormat() {
	type model struct {
		Value string
	}
	d := NewDecoder([]DecoderField{{Key: "Value"}})

	// cursor must be a base64 encoded string
	_, err := d.Decode("123", model{})
	s.Equal(ErrInvalidCursor, err)

	// cursor must be a valid json
	c := base64.StdEncoding.EncodeToString([]byte(`["123"}`))
	_, err = d.Decode(c, model{})
	s.Equal(ErrInvalidCursor, err)

	// cursor must be a json array
	c = base64.StdEncoding.EncodeToString([]byte(`{"value": "123"}`))
	_, err = d.Decode(c, model{})
	s.Equal(ErrInvalidCursor, err)
}

func (s *decoderSuite) TestDecodeInvalidCursorType() {
	c, _ := NewEncoder([]EncoderField{{Key: "Value"}}).Encode(struct{ Value int }{123})
	_, err := NewDecoder([]DecoderField{{Key: "Value"}}).Decode(c, struct{ Value string }{})
	s.Equal(ErrInvalidCursor, err)
}

/* decode struct */

func (s *decoderSuite) TestDecodeStructInvalidModel() {
	err := NewDecoder([]DecoderField{{Key: "Value"}}).DecodeStruct("123", struct{ ID string }{})
	s.Equal(ErrInvalidModel, err)
}

func (s *decoderSuite) TestDecodeStructInvalidCursor() {
	err := NewDecoder([]DecoderField{{Key: "Value"}}).DecodeStruct("123", struct{ Value string }{})
	s.Equal(ErrInvalidCursor, err)
}
