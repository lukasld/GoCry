package cliSimpleCall

import (
    "errors"
)

var (
    ErrStdErr   = errors.New("1PW-StdErr:")
)

type Error struct {
    appErr error
    stdErr error
}

func (e Error) Error() string {
    return errors.Join(e.appErr, e.stdErr).Error()
}

func NewError(appErr, stdErr error) error {
    return Error{
        appErr: appErr,
        stdErr: stdErr,
    }
}
