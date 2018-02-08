package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecToJson(t *testing.T) {

	assert.True(t, IncludedIn([]string{"a", "b"}, []string{"b", "a", "c"}))
	assert.False(t, IncludedIn([]string{"a", "b"}, []string{"a", "c"}))
	assert.False(t, IncludedIn([]string{"a", "b"}, []string{"b", "c"}))
	assert.True(t, IncludedIn([]string{}, []string{"b", "c"}))
	assert.True(t, IncludedIn([]string{}, []string{}))
}
