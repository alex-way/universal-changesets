package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseChangeset(t *testing.T) {
	cs, _ := ParseVersion("1.2.3")

	assert.Equal(t, 1, cs.Major)
	assert.Equal(t, 2, cs.Minor)
	assert.Equal(t, 3, cs.Patch)
}

func TestIncrementMajor(t *testing.T) {
	cs := Version{0, 0, 0}
	cs.IncrementMajor()

	assert.Equal(t, 1, cs.Major)
	assert.Equal(t, 0, cs.Minor)
	assert.Equal(t, 0, cs.Patch)
}

func TestIncrementMinor(t *testing.T) {
	cs := Version{0, 0, 0}
	cs.IncrementMinor()

	assert.Equal(t, 0, cs.Major)
	assert.Equal(t, 1, cs.Minor)
	assert.Equal(t, 0, cs.Patch)
}

func TestIncrementPatch(t *testing.T) {
	cs := Version{0, 0, 0}
	cs.IncrementPatch()

	assert.Equal(t, 0, cs.Major)
	assert.Equal(t, 0, cs.Minor)
	assert.Equal(t, 1, cs.Patch)
}
