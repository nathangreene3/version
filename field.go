package version

// Field corresponds to the fields major, minor, and patch.
type Field int

const (
	Patch Field = 1 + iota
	Minor
	Major
)
