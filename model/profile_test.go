package model

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProfile(t *testing.T) {
	request := ProfileRequest{
		FirstName: "John",
		LastName:  "Doe",
	}
	id := "a28435c0-0007-426c-8a22-acfa89bad6f1"
	profile := request.Profile(WithUUID(uuid.Must(uuid.FromString(id))))
	assert.Equal(t, "John", profile.FirstName)
	assert.Equal(t, "Doe", profile.LastName)
	assert.Equal(t, id, profile.ExternalID.String())
}
