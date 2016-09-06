/* Middleware stuff for the UTM server.
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

package system

import (
    "github.com/u-blox/utm/service/utilities"
    "gopkg.in/mgo.v2"
    "net/http"
)

var initialSession *mgo.Session = nil

// A handler function pulled in by Negroni to establish a database session
func MgoOpenDbSession(response http.ResponseWriter, request *http.Request, nextFunc http.HandlerFunc) {

    var err error = nil
    
    if initialSession == nil {
        initialSession, err = mgo.Dial("127.0.0.1:27017")
    }

    if err == nil {
        session := initialSession.Clone()
        db := session.DB("utm-db")
        utilities.SetDB(request, db)
    }
    
    nextFunc(response, request)
}

// A handler function pulled in by Negroni to release a database session
func MgoCloseDbSession(response http.ResponseWriter, request *http.Request, nextFunc http.HandlerFunc) {
    db := utilities.GetDB(request)
    if db != nil {
        db.Session.Close()
    }    

    nextFunc(response, request)
}

// A clean-up function that should be called at end of time to close the initial session
func MgoCleanup() {
    if initialSession != nil {
        initialSession.Close()
    }    
}

/* End Of File */
