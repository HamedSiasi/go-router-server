/* Utils to deal with error cases on the UTM server.
 *
 * Copyright (C) u-blox Melbourn Ltd
 * u-blox Melbourn Ltd, Melbourn, UK
 *
 * All rights reserved.
 *
 * This source file is the sole property of u-blox Melbourn Ltd.
 * Reproduction or utilization of this source in whole or part is
 * forbidden without the written consent of u-blox Melbourn Ltd.
 */

package utilities

import (
    "errors"
    "log"
    "net/http"
    "runtime"
)

/// Definition of an error
type Error struct {
    innerError error
    buf        []byte
    status     int
    annotation string
}

/// Create a new error
func NewAppError(optionalMessage string, optionalInnerError error, captureStackTrace bool, contextDependantStatusInt int) *Error {
    err := Error{
        innerError: optionalInnerError,
        buf:        nil,
        status:     contextDependantStatusInt,
        annotation: optionalMessage,
    }

    if captureStackTrace {
        err.buf = make([]byte, 1<<16)
        runtime.Stack(err.buf, false)
    }

    return &err
}

/// Create a server error
func ServerError(err error) *Error {
    if err == nil {
        return nil
    }

    apperr := NewAppError("", err, true, 500)
    return apperr
}

/// Create a client error
func ClientError(message string, status int) *Error {
    apperr := NewAppError("", errors.New(message), false, status)
    return apperr
}

/// A handler which takes rest request handler function as an argument
// and returns an error
type Handler func(w http.ResponseWriter, req *http.Request) *Error

/// TODO
func (handler Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    apperr := handler(w, req)
    if apperr != nil {
        var message string
        if apperr.status >= 500 {
            message = "Internal Server Error"
        } else {
            message = apperr.Error()
        }
        log.Printf("ERROR (%d): %v\n%s", apperr.status, apperr.innerError, apperr.buf)
        http.Error(w, message, apperr.status)
    }
}

/// TODO
func (apperr *Error) Error() string {
    if apperr.annotation == "" {
        if apperr.innerError != nil {
            return apperr.innerError.Error()
        } else {
            return "Unspecified error"
        }
    } else {
        if apperr.innerError != nil {
            return apperr.annotation + " :\n" + apperr.innerError.Error()
        } else {
            return apperr.innerError.Error()
        }
    }
}

/* End Of File */
