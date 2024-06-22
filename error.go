package ll1

import "errors"

// ErrInvalidArgument is a general error returned when the input was not valid.
var ErrInvalidArgument = errors.New("invalid argument")

// ErrEmptyExpr is returned when a method was supplied an expression which was empty
// or the input simplified to an empty expression. It can be used as a signal from
// simplifiers upstream. For instance if an expression in a sequence simplifies to the
// empty expression, it can be omitted. It is a kind of ErrInvalidArgument.
var ErrEmptyExpr = errors.New("empty expression: invalid argument")
