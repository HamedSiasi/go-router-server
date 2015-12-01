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
func GetDB(request *http.Request) *mgo.Database {
    db := context.Get(request, "db")
    if db != nil {
        return db.(*mgo.Database)
    }
    
    return nil
}

/// Set up the database from an HTTP request
func SetDB(request *http.Request, db *mgo.Database) {
    context.Set(request, "db", db)
}

// Insert a struct into the database.  Note that this only
// works if passed a pointer to a struct.
func InsertDB (collectionString string, documents ...interface{}) error {    
    session, err := mgo.Dial("127.0.0.1:27017")
    if err == nil {
        defer session.Close()
        db := session.DB("utm-db")
        collection := db.C(collectionString)
        for _, document := range documents {
            err = collection.Insert(document)
            if err != nil {
                break
            }
        }
    }    
    return err
}

/* End Of File */
