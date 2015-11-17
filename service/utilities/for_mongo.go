/* Utilities to do with  the database for the UTM server.
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
    "github.com/gorilla/context"
    "gopkg.in/mgo.v2"
    "net/http"
)

/// Return the current database from an HTTP request
func GetDB(r *http.Request) *mgo.Database {
    db := context.Get(r, "db")
    return db.(*mgo.Database)
}

/// Set up the database from an HTTP request
func SetDB(r *http.Request, db *mgo.Database) {
    context.Set(r, "db", db)
}

/* End Of File */
