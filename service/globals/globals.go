/* Global types/variables file for UTM server.
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

package globals

import (
    "fmt"
    "log"
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

type DebugLevel int

const (
	DEBUG_OFF   DebugLevel = 0
    DEBUG_ERROR DebugLevel = 1
    DEBUG_TRACE DebugLevel = 2
    DEBUG_INFO  DebugLevel = 3
)

/// Definition of an error
type Error struct {
    InnerError error
    Buf        []byte
    Status     int
    Annotation string
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// Level of debug required
const Dbg DebugLevel = DEBUG_TRACE



// Log  prefix so that we can tell who we are
var LogTag string = "UTM-API"


//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

/// Fail func
func FailOnError (err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
	    panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

/// TODO
func (apperr *Error) Error() string {
    if apperr.Annotation == "" {
        if apperr.InnerError != nil {
            return apperr.InnerError.Error()
        } else {
            return "Unspecified error"
        }
    } else {
        if apperr.InnerError != nil {
            return apperr.Annotation + " :\n" + apperr.InnerError.Error()
        } else {
            return apperr.InnerError.Error()
        }
    }
}

// General debugging
func (d DebugLevel) PrintfInfo(s string, a ...interface{}) {
	if (d > 0) && (d >= DEBUG_INFO) {
		log.Printf(s, a...)
	}
}
func (d DebugLevel) PrintfTrace(s string, a ...interface{}) {
	if (d > 0) && (d >= DEBUG_TRACE) {
		log.Printf(s, a...)
	}
}
func (d DebugLevel) PrintfError(s string, a ...interface{}) {
	if (d > 0) && (d >= DEBUG_ERROR)  {
		log.Printf(s, a...)
	}
}

/* End Of File */