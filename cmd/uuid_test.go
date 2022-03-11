package cmd

import (
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"
)

func TestCheckParam(t *testing.T) {
	count = -1
	assert.Equal(t, false, checkParam())

	count = 100
	assert.Equal(t, false, checkParam())

	count = 33
	*isTrim = true
	assert.Equal(t, false, checkParam())

	*isTrim = false
	count = 30
	assert.Equal(t, true, checkParam())

	*isTrim = true
	count = 30
	assert.Equal(t, true, checkParam())
}

func TestHandleUid(t *testing.T) {
	uid, _ := uuid.GenerateUUID()

	*isTrim = true
	count = 30
	assert.Equal(t, true, checkParam())
	assert.Equal(t, 30, len(handleUID(uid)))

	*isTrim = false
	count = 10
	assert.Equal(t, true, checkParam())
	assert.Equal(t, 10, len(handleUID(uid)))
}
