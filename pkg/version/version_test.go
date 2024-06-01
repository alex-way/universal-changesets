package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	cs, _ := ParseVersion("1.2.3")

	assert.Equal(t, 1, cs.Major)
	assert.Equal(t, 2, cs.Minor)
	assert.Equal(t, 3, cs.Patch)
}

func TestBumpMajor(t *testing.T) {
	cs := Version{0, 0, 0}
	cs.BumpMajor()

	assert.Equal(t, 1, cs.Major)
	assert.Equal(t, 0, cs.Minor)
	assert.Equal(t, 0, cs.Patch)
}

func TestBumpMinor(t *testing.T) {
	cs := Version{0, 0, 0}
	cs.BumpMinor()

	assert.Equal(t, 0, cs.Major)
	assert.Equal(t, 1, cs.Minor)
	assert.Equal(t, 0, cs.Patch)
}

func TestBumpPatch(t *testing.T) {
	cs := Version{0, 0, 0}
	cs.BumpPatch()

	assert.Equal(t, 0, cs.Major)
	assert.Equal(t, 0, cs.Minor)
	assert.Equal(t, 1, cs.Patch)
}

func TestParseBumpType(t *testing.T) {
	bump_type, err := ParseBumpType("major")
	assert.NoError(t, err)
	assert.Equal(t, Major, bump_type)

	bump_type, err = ParseBumpType("minor")
	assert.NoError(t, err)
	assert.Equal(t, Minor, bump_type)

	bump_type, err = ParseBumpType("patch")
	assert.NoError(t, err)
	assert.Equal(t, Patch, bump_type)

	bump_type, err = ParseBumpType("none")
	assert.NoError(t, err)
	assert.Equal(t, None, bump_type)

	_, err = ParseBumpType("invalid")
	assert.Error(t, err)
}

func TestVersionString(t *testing.T) {
	cs := Version{1, 2, 3}
	assert.Equal(t, "1.2.3", cs.String())
}
