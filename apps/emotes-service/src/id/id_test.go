package id

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdGeneration(t *testing.T) {
	generatedId := New("generated")
	idRegex := regexp.MustCompile(`^generated_\w{24}$`)

	assert.True(t, idRegex.MatchString(generatedId))
}
