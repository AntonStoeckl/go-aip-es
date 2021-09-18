package shared

import "github.com/cockroachdb/errors"

var (
	ErrInputIsInvalid             = errors.New("input is invalid")
	ErrDomainConstraintsViolation = errors.New("domain constraints violation")
)

func MarkAndWrapError(original, markAs error, wrapWith string) error {
	return errors.Mark(errors.Wrap(original, wrapWith), markAs)
}
