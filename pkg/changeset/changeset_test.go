package changeset

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateChangesetNameHasThreeParts(t *testing.T) {
	filename := generateChangesetName()

	assert.Equal(t, 3, len(strings.Split(filename, "-")))
}
