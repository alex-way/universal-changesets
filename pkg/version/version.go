package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

type BumpType int8

func (bump_type BumpType) String() string {
	switch bump_type {
	case Major:
		return "major"
	case Minor:
		return "minor"
	case Patch:
		return "patch"
	}
	return "none"
}

func ParseBumpType(s string) (BumpType, error) {
	switch s {
	case "major":
		return Major, nil
	case "minor":
		return Minor, nil
	case "patch":
		return Patch, nil
	case "none":
		return None, nil
	default:
		return 0, errors.New("invalid bump type. Must be one of: major, minor, patch, none")
	}
}

const (
	Undetermined BumpType = -1
	None         BumpType = 0
	Patch        BumpType = 1
	Minor        BumpType = 2
	Major        BumpType = 3
)

func (cs *Version) Bump(t BumpType) {
	switch t {
	case Major:
		cs.BumpMajor()
	case Minor:
		cs.BumpMinor()
	case Patch:
		cs.BumpPatch()
	}
}

func (cs *Version) BumpPatch() {
	cs.Patch += 1
}

func (cs *Version) BumpMinor() {
	cs.Minor += 1
	cs.Patch = 0
}

func (cs *Version) BumpMajor() {
	cs.Major += 1
	cs.Minor = 0
	cs.Patch = 0
}

func ParseVersion(s string) (Version, error) {
	cs := Version{}
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return cs, errors.New("invalid change set")
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return cs, errors.New("invalid change set")
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return cs, errors.New("invalid change set")
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return cs, errors.New("invalid change set")
	}
	return Version{major, minor, patch}, nil
}
