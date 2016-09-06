/* User data elements of the UTM server.
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

package controllers

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/u-blox/utm/service/models"
    "github.com/u-blox/utm/service/utilities"
)

type User struct{}

/// Method for the user struct that takes an HTTP request and
// returns an HTTP response containing the user's details for
// their ID.
func (u *User) Get(response http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    id := vars["id"]

    db := utilities.GetDB(request)
    if db != nil {
        user := new(models.User)
        err := user.Get(db, id)
        if err != nil {
            response.WriteHeader(404)
        } else {
            user.Password = ""
            out, _ := json.Marshal(user)
            response.Write(out)
        }
    } else {
        response.WriteHeader(404)        
    }    
}

// TODO
func (u *User) Profile(response http.ResponseWriter, request *http.Request) {
    user_id, _ := utilities.GetUserId(request)
    db := utilities.GetDB(request)
    if db != nil {
        user := new(models.User)
        user.Get(db, user_id)
        user.Password = ""
        out, _ := json.Marshal(user)
        response.Write(out)
    } else {
        response.WriteHeader(404)        
    }    
}

/* End Of File */
