package collector

import (
	"github.com/giantswarm/microerror"
)

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return microerror.Cause(err) == executionFailedError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var nameMatchError = &microerror.Error{
	Kind: "nameMatchError",
}

// IsNameMatch asserts nameMatchError.
func IsNameMatch(err error) bool {
	return microerror.Cause(err) == nameMatchError
}
