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
	"github.com/u-blox/utm-server/service/globals"
	"log"
	"net/http"
	"runtime"
)

/// Create a new error
func NewAppError(optionalMessage string, optionalInnerError error, captureStackTrace bool, contextDependantStatusInt int) *globals.Error {
	err := globals.Error{
		InnerError: optionalInnerError,
		Buf:        nil,
		Status:     contextDependantStatusInt,
		Annotation: optionalMessage,
	}

	if captureStackTrace {
		err.Buf = make([]byte, 1<<16)
		runtime.Stack(err.Buf, false)
	}

	return &err
}

/// Create a server error
func ServerError(err error) *globals.Error {
	if err == nil {
		return nil
	}

	apperr := NewAppError("", err, true, 500)
	return apperr
}

/// Create a client error
func ClientError(message string, status int) *globals.Error {
	apperr := NewAppError("", errors.New(message), false, status)
	return apperr
}

/// A handler which takes rest request handler function as an argument
// and returns an error
type Handler func(w http.ResponseWriter, req *http.Request) *globals.Error

/// TODO
func (handler Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	apperr := handler(response, request)
	if apperr != nil {
		var message string
		if apperr.Status >= 500 {
			message = "Internal Server Error"
		} else {
			message = apperr.Error()
		}
		log.Printf("ERROR (%d): %v\n%s", apperr.Status, apperr.InnerError, apperr.Buf)
		http.Error(response, message, apperr.Status)
	}
}

/* End Of File */
