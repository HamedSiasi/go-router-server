/* Debug dump functions for the UTM server.
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
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
)

type DebugLevel int

const (
	DEBUG_OFF   DebugLevel = 0
    DEBUG_ERROR DebugLevel = 1
    DEBUG_TRACE DebugLevel = 2
    DEBUG_INFO  DebugLevel = 3
)

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

// Debugging functions for logging requests and responses
func DumpRequest(req *http.Request) {
    content, err := httputil.DumpRequest(req, true)
    if err != nil {
        log.Println("Error dumping request: %v\n", err)
    } else {
        fmt.Println(string(content))
    }
}

func DumpResponse(req *http.Response) {
    content, err := httputil.DumpResponse(req, true)
    if err != nil {
        log.Println("Error dumping response: %v\n", err)
    } else {
        fmt.Println(string(content))
    }
}

/* End Of File */
