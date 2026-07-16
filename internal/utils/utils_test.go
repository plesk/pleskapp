// Copyright 1999-2026. WebPros International GmbH.

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	password := GeneratePassword(10)
	assert.Len(t, password, 10)

	password2 := GeneratePassword(10)
	assert.NotEqual(t, password, password2)
}

func TestGenerateUsername(t *testing.T) {
	username := GenerateUsername(8)
	assert.Len(t, username, 8)

	username2 := GenerateUsername(8)
	assert.NotEqual(t, username, username2)
}
