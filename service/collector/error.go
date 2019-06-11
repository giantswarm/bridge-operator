package collector

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidFileError = &microerror.Error{
	Kind: "invalidFileError",
}

// IsInvalidFile asserts invalidFileError.
func IsInvalidFile(err error) bool {
	return microerror.Cause(err) == invalidFileError
}
