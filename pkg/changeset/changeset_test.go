package changeset

import (
	"strings"
	"testing"

	"github.com/alex-way/changesets/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestGenerateChangeNameHasThreeParts(t *testing.T) {
	filename := generateChangeName()

	assert.Equal(t, 3, len(strings.Split(filename, "-")))
}

func TestDetermineFinalBumpType(t *testing.T) {
	changeset := Changeset{
		Changes: []Change{
			{BumpType: version.Major, Message: ""},
			{BumpType: version.Minor, Message: ""},
			{BumpType: version.Patch, Message: ""},
			{BumpType: version.None, Message: ""},
		},
		CurrentVersion: version.Version{Major: 0, Minor: 0, Patch: 0},
	}
	assert.Equal(t, version.Major, changeset.DetermineFinalBumpType())

	changeset = Changeset{
		Changes: []Change{
			{BumpType: version.Minor, Message: ""},
			{BumpType: version.Patch, Message: ""},
			{BumpType: version.None, Message: ""},
		},
		CurrentVersion: version.Version{Major: 0, Minor: 0, Patch: 0},
	}
	assert.Equal(t, version.Minor, changeset.DetermineFinalBumpType())

	changeset = Changeset{
		Changes: []Change{
			{BumpType: version.Patch, Message: ""},
			{BumpType: version.None, Message: ""},
		},
		CurrentVersion: version.Version{Major: 0, Minor: 0, Patch: 0},
	}
	assert.Equal(t, version.Patch, changeset.DetermineFinalBumpType())

	changeset = Changeset{
		Changes: []Change{
			{BumpType: version.None, Message: ""},
		},
		CurrentVersion: version.Version{Major: 0, Minor: 0, Patch: 0},
	}
	assert.Equal(t, version.None, changeset.DetermineFinalBumpType())

	changeset = Changeset{
		Changes:        []Change{},
		CurrentVersion: version.Version{Major: 0, Minor: 0, Patch: 0},
	}
	assert.Equal(t, version.Undetermined, changeset.DetermineFinalBumpType())
}
