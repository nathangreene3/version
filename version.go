package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	// versionPattern is the regular expression that defines a semantic
	// version within Go's module system. Note Go's definition is based
	// on Semantic Versioning, but deviates in several aspects. Build
	// tags are not supported and versions are prefixed with v.
	//
	// Sources:
	// 	- https://go.dev/doc/modules/version-numbers
	// 	- https://semver.org/
	versionPattern = "^v(0|[1-9][0-9]*)(.(0|[1-9][0-9]*)){2}(-([a-zA-Z0-9]+)(.[a-zA-Z0-9]+)*)?$"
)

// A Version holds semantic version information as defined by Go's
// module system. See the following link for details.
//
// Source: https://go.dev/doc/modules/version-numbers
type Version struct {
	major, minor, patch uint
	preRelease          string
}

// New returns a version with any given options.
func New(major, minor, patch uint, opts ...Option) (Version, error) {
	v := Version{
		major: major,
		minor: minor,
		patch: patch,
	}

	for i := 0; i < len(opts); i++ {
		opts[i](&v)
	}

	if !IsValid(v.String()) {
		return Version{}, ErrInvalidVersion
	}

	return v, nil
}

// Parse returns a version parsed from a given string.
//
// Examples:
//   - In-development:   v0.y.z (major version is zero)
//   - Pre-release:      vx.y.z-alpha.a (major version is non-zero)
//   - Official release: vx.y.z (major version is non-zero)
func Parse(s string) (Version, error) {
	if !IsValid(s) {
		return Version{}, ErrInvalidVersion
	}

	var fields, preRelease string

	if i := strings.Index(s, "-"); i < 0 {
		fields = s[1:]
	} else {
		fields = s[1:i]
		preRelease = s[i+1:]
	}

	ss := strings.Split(fields, ".")

	major, err := strconv.ParseUint(ss[0], 10, strconv.IntSize)
	if err != nil {
		return Version{}, fmt.Errorf("strconv.ParseUint: %w", err)
	}

	minor, err := strconv.ParseUint(ss[1], 10, strconv.IntSize)
	if err != nil {
		return Version{}, fmt.Errorf("strconv.ParseUint: %w", err)
	}

	patch, err := strconv.ParseUint(ss[2], 10, strconv.IntSize)
	if err != nil {
		return Version{}, fmt.Errorf("strconv.ParseUint: %w", err)
	}

	v, err := New(uint(major), uint(minor), uint(patch), PreRelease(preRelease))
	if err != nil {
		return Version{}, fmt.Errorf("New: %w", err)
	}

	return v, nil
}

// IsValid determines if a string is a valid representation of a
// version.
func IsValid(version string) bool {
	ok, err := regexp.MatchString(versionPattern, version)
	return ok && err == nil
}

// Bump returns the next version given a field.
func (v Version) Bump(f Field, opts ...Option) (Version, error) {
	switch f {
	case Major:
		v1, err := New(v.major+1, 0, 0, opts...)
		if err != nil {
			return Version{}, fmt.Errorf("New: %w", err)
		}

		return v1, nil
	case Minor:
		v1, err := New(v.major, v.minor+1, 0, opts...)
		if err != nil {
			return Version{}, fmt.Errorf("New: %w", err)
		}

		return v1, nil
	case Patch:
		v1, err := New(v.major, v.minor, v.patch+1, opts...)
		if err != nil {
			return Version{}, fmt.Errorf("New: %w", err)
		}

		return v1, nil
	default:
		return Version{}, ErrInvalidField
	}
}

// Compare returns the comparison of two versions.
func (v Version) Compare(v1 Version) int {
	switch {
	case v.major < v1.major:
		return -1
	case v1.major < v.major:
		return 1
	case v.minor < v1.minor:
		return -1
	case v1.minor < v.minor:
		return 1
	case v.patch < v1.patch:
		return -1
	case v1.patch < v.patch:
		return 1
	case v.preRelease < v1.preRelease:
		return -1
	case v1.preRelease < v.preRelease:
		return 1
	default:
		return 0
	}
}

// IsValid determines if a version is valid.
func (v Version) IsValid() bool {
	v1, err := Parse(v.String())
	return err == nil && v.Compare(v1) == 0
}

// Major returns the major version.
func (v Version) Major() uint {
	return v.major
}

// Minor returns the minor version.
func (v Version) Minor() uint {
	return v.minor
}

// Patch returns the patch version.
func (v Version) Patch() uint {
	return v.patch
}

// PreRelease returns the pre-release.
func (v Version) PreRelease() string {
	return v.preRelease
}

// String returns a representation of a version.
func (v Version) String() string {
	if v.preRelease == "" {
		return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
	}

	return fmt.Sprintf("v%d.%d.%d-%s", v.major, v.minor, v.patch, v.preRelease)
}
