package version

import (
	"errors"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

type IncrementType int

func (increment_type IncrementType) String() string {
	switch increment_type {
	case Major:
		return "major"
	case Minor:
		return "minor"
	case Patch:
		return "patch"
	}
	return "none"
}

func ParseIncrementType(s string) (IncrementType, error) {
	switch s {
	case "major":
		return Major, nil
	case "minor":
		return Minor, nil
	case "patch":
		return Patch, nil
	default:
		return 0, errors.New("invalid increment type. Must be one of: major, minor, patch")
	}
}

const (
	Major IncrementType = iota
	Minor
	Patch
	None
)

func (cs *Version) Increment(t IncrementType) {
	switch t {
	case Major:
		cs.IncrementMajor()
	case Minor:
		cs.IncrementMinor()
	case Patch:
		cs.IncrementPatch()
	}
}

func (cs *Version) IncrementPatch() {
	cs.Patch += 1
}

func (cs *Version) IncrementMinor() {
	cs.Minor += 1
	cs.Patch = 0
}

func (cs *Version) IncrementMajor() {
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
