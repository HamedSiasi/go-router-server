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
    "github.com/mmanjoura/utm-admin-v0.8/service/utilities"
    "gopkg.in/mgo.v2"
    "net/http"
)

// A handler function pulled in by Negroni to establish a database session
func MgoMiddleware(response http.ResponseWriter, request *http.Request, nextFunc http.HandlerFunc) {
    session, err := mgo.Dial("127.0.0.1:27017")

    if err != nil {
        panic(err)
    }

    //reqSession := session.Clone()
    //defer reqSession.Close()
    //defer session.Close()
    dbs := session.DB("utm-db")
    utilities.SetDB(request, dbs)
    nextFunc(response, request)
}

/* End Of File */
