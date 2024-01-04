package config

import (
	"fmt"
	"strings"
)

// the final version string
var version string

// -ldflags "-X github.com/maqdev/go-be-template/config/version.branch=master"
var branch string

// -ldflags "-X github.com/maqdev/go-be-template/config/version.build=alpha"
var build string

const (
	VersionMajor = 0
	VersionMinor = 9
	VersionPatch = 0
	VersionTag   = "" // example: "rc1"
)

func VersionString() string {
	version = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	if VersionTag != "" {
		version += "-" + VersionTag
	}
	parts := []string{}
	if build != "" {
		parts = append(parts, build)
	}
	if branch != "" {
		parts = append(parts, branch)
	}
	if len(parts) > 0 {
		version += "+" + strings.Join(parts, ".")
	}
	return version
}
