/* Authentication elements of the UTM server.
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
    "fmt"
    "github.com/goincremental/negroni-sessions"
    "github.com/u-blox/utm/service/models"
    "github.com/u-blox/utm/service/utilities"
    "net/http"
)

type Credentials struct {
    Email    string
    Password string
}

type Auth struct{}

/// Method for the authentication struct that takes
// an HTTP response and returns an HTTP request containing
// the user's details.
func (a *Auth) User(response http.ResponseWriter, request *http.Request) {
    db := utilities.GetDB(request)
    if db != nil {
        session := sessions.GetSession(request)
        user_id := session.Get("user_id")
        fmt.Println(user_id)
        if user_id == nil {
            response.WriteHeader(403)
        } else {
            user := new(models.User)
            user.Get(db, user_id.(string))
            fmt.Println(user)
            outData, _ := json.Marshal(user)
            response.Write(outData)
        }
    } else {
        response.WriteHeader(404)        
    }    
}

/// Method for the authentication struct that takes
// an HTTP response and returns an HTTP request containing
// all the details of all the users in a company.
func (a *Auth) Users(response http.ResponseWriter, request *http.Request) {
    db := utilities.GetDB(request)
    if db != nil {
        session := sessions.GetSession(request)
        user_company := session.Get("user_company")
        fmt.Println(user_company)
        if user_company == nil {
            response.WriteHeader(403)
        } else {
            user := new(models.User)
            users := user.GetCompanyUsers(db, user_company.(string))
            fmt.Println(users)
            outData, _ := json.Marshal(users)
            response.Write(outData)
        }
    } else {
        response.WriteHeader(404)        
    }    
}

/// Method for the authentication struct that takes
// an HTTP response and returns an HTTP request containing
// all the UUIDs for a user.
func (a *Auth) Uuids(response http.ResponseWriter, request *http.Request) {
    db := utilities.GetDB(request)
    if db != nil {
        session := sessions.GetSession(request)
        user_company := session.Get("user_company")
        fmt.Println(user_company)
        if user_company == nil {
            response.WriteHeader(403)
        } else {
            user := new(models.User)
            uuids := user.GetUserUuids(db, user_company.(string))
            fmt.Println(uuids)
            outData, _ := json.Marshal(uuids)
            response.Write(outData)
        }
    } else {
        response.WriteHeader(404)        
    }    
}

/// Method for the authentication struct that takes
// an HTTP response and returns an HTTP request containing
// the friendly names of the devices against a user's e-mail.
func (a *Auth) UserUEs(response http.ResponseWriter, request *http.Request) {
    db := utilities.GetDB(request)
    if db != nil {
        session := sessions.GetSession(request)
        user_email := session.Get("user_email")
        fmt.Println(user_email)
        if user_email == nil {
            response.WriteHeader(403)
        } else {
            user := new(models.User)
            ues := user.GetUEs(db, user_email.(string))
            fmt.Println(ues)
            outData, _ := json.Marshal(ues)
            response.Write(outData)
        }
    } else {
        response.WriteHeader(404)        
    }    
}

/* End Of File */
