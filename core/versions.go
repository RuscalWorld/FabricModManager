package core

import (
	"github.com/hashicorp/go-version"
	"strings"
)

func CheckVersions(ver string, constraint interface{}) bool {
	if versions, ok := constraint.([]interface{}); ok {
		for _, required := range versions {
			if required == ver {
				return true
			}
		}

		return false
	}

	required := constraint.(string)
	if required == "*" {
		return true
	}

	ver = strings.TrimPrefix(ver, "v")
	required = strings.TrimPrefix(required, "v")

	vConstraint, err := version.NewConstraint(required)
	if err != nil {
		return false
	}

	vVer, err := version.NewVersion(ver)
	if err != nil {
		return false
	}

	return vConstraint.Check(vVer)
}
