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

var logTag string = "UTM"

/// TODO
func MgoMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    session, err := mgo.Dial("127.0.0.1:27017")

    if err != nil {
        panic(err)
    }

    reqSession := session.Clone()
    defer reqSession.Close()
    dbs := reqSession.DB("utm-db")
    utilities.SetDB(r, dbs)
    next(rw, r)
}

/* End Of File */
