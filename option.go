package version

// An Option modifies a version. The version should be validated
// immediately after modification.
type Option func(v *Version)

// PreRelease returns an option that sets the pre-release field. The
// version should be validated immediately after setting the pre-release
// field.
func PreRelease(s string) Option {
	return func(v *Version) {
		v.preRelease = s
	}
}
