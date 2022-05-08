package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	arg := []string{"1", "2"}
	p := arg[2:]
	assert.Equal(t, 0, len(p))
}
