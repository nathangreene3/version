package version_test

import (
	"testing"

	. "github.com/nathangreene3/version"
	. "github.com/onsi/gomega"
)

// TestVersion ensures New and Parsed produce consistent values. New
// should produce a valid version and Parse should parse its string
// returning an equivalent copy.
func TestVersion(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		major      uint
		minor      uint
		patch      uint
		preRelease string
		version    string
	}{
		{
			version: "v0.0.0",
		},
		{
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "Abc.1.2.3",
			version:    "v1.2.3-Abc.1.2.3",
		},
	}

	for _, test := range tests {
		v0, err := New(test.major, test.minor, test.patch, PreRelease(test.preRelease))
		g.Expect(err).To(BeNil())
		g.Expect(v0.IsValid()).To(BeTrue())
		g.Expect(v0.Major()).To(Equal(test.major))
		g.Expect(v0.Minor()).To(Equal(test.minor))
		g.Expect(v0.Patch()).To(Equal(test.patch))
		g.Expect(v0.PreRelease()).To(Equal(test.preRelease))

		vs0 := v0.String()
		g.Expect(IsValid(vs0)).To(BeTrue())
		g.Expect(vs0).To(Equal(test.version))

		v1, err := Parse(vs0)
		g.Expect(err).To(BeNil())
		g.Expect(v0.Compare(v1)).To(BeZero())
		g.Expect(v1.IsValid()).To(BeTrue())
		g.Expect(v1.Major()).To(Equal(test.major))
		g.Expect(v1.Minor()).To(Equal(test.minor))
		g.Expect(v1.Patch()).To(Equal(test.patch))
		g.Expect(v1.PreRelease()).To(Equal(test.preRelease))
	}
}
