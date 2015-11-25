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
    "github.com/robmeades/utm/service/models"
    "github.com/robmeades/utm/service/utilities"
    "net/http"
)

type Credentials struct {
    Email    string
    Password string
}

type Auth struct{}

/// Login method for the authentication struct that takes
// an HTTP request and returns an HTTP response.
// A session is created if successful.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the login page!")
	// decoder := json.NewDecoder(r.Body)
	// credentials := new(Credentials)
	// err := decoder.Decode(&credentials)
	// if err != nil {
	// 	panic(err)
	// }

	// db := utilities.GetDB(r)
	// user := new(models.User)
	// err = user.Authenticate(db, credentials.Email, credentials.Password)
	// if err == nil {
	// 	session := sessions.GetSession(r)
	// 	session.Set("user_id", user.ID.Hex())
	// 	session.Set("user_company", user.Company)
	// 	session.Set("user_email", user.Email)
	// 	w.WriteHeader(202)

	// } else {
	// 	w.WriteHeader(404)
	// }

	email := r.FormValue("email")
	password := r.FormValue("password")

	db := utilities.GetDB(r)
	user := new(models.User)
	err := user.Authenticate(db, email, password)
	if err == nil {
		//session := sessions.GetSession(request)
		// session.Set("user_id", user.ID.Hex())
		// session.Set("user_company", user.Company)
		// session.Set("user_email", user.Email)
		//response.WriteHeader(202)
		http.Redirect(w, r, "http://localhost:3000", 202)
		http.Redirect(w, r, "/", 202)
		fmt.Fprintf(w, "error is nil!")
		//http.Redirect(response, request, "/", 202)

	} else {
		//fmt.Fprintf(w, "There is Error!")
		//http.Redirect(w, r, "http://localhost:3000", 202)
		// http.Redirect(w, r, "http://localhost:3000/lastestState", 202)
		// w.WriteHeader(404)
	}
}

/// Logout method for the authentication struct that takes
// an HTTP request and returns an HTTP response.
// The session is deleted, if it is found.
func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
    session := sessions.GetSession(r)
    user_id := session.Get("user_id")
    fmt.Println(user_id)
    if user_id == nil {
        w.WriteHeader(403)
        http.Redirect(w, r, "/", 403)
    } else {
        session.Delete("user_id")
        http.Redirect(w, r, "/", 202)
    }
}

/// Method for the authentication struct that takes
// an HTTP request and returns an HTTP response containing
// the user's details.
func (a *Auth) User(w http.ResponseWriter, r *http.Request) {
    db := utilities.GetDB(r)
    session := sessions.GetSession(r)
    user_id := session.Get("user_id")
    fmt.Println(user_id)
    if user_id == nil {
        w.WriteHeader(403)
    } else {
        user := new(models.User)
        user.Get(db, user_id.(string))
        fmt.Println(user)
        outData, _ := json.Marshal(user)
        w.Write(outData)
    }
}

/// Method for the authentication struct that takes
// an HTTP request and returns an HTTP response containing
// all the details of all the users in a company.
func (a *Auth) Users(w http.ResponseWriter, r *http.Request) {
    db := utilities.GetDB(r)
    session := sessions.GetSession(r)
    user_company := session.Get("user_company")
    fmt.Println(user_company)
    if user_company == nil {
        w.WriteHeader(403)
    } else {
        user := new(models.User)
        users := user.GetCompanyUsers(db, user_company.(string))
        fmt.Println(users)
        outData, _ := json.Marshal(users)
        w.Write(outData)
    }
}

/// Method for the authentication struct that takes
// an HTTP request and returns an HTTP response containing
// all the UUIDs for a user.
func (a *Auth) Uuids(w http.ResponseWriter, r *http.Request) {
    db := utilities.GetDB(r)
    session := sessions.GetSession(r)
    user_company := session.Get("user_company")
    fmt.Println(user_company)
    if user_company == nil {
        w.WriteHeader(403)
    } else {
        user := new(models.User)
        uuids := user.GetUserUuids(db, user_company.(string))
        fmt.Println(uuids)
        outData, _ := json.Marshal(uuids)
        w.Write(outData)
    }
}

/// Method for the authentication struct that takes
// an HTTP request and returns an HTTP response containing
// the friendly names of the devices against a user's e-mail.
func (a *Auth) UserUEs(w http.ResponseWriter, r *http.Request) {
    db := utilities.GetDB(r)
    session := sessions.GetSession(r)
    user_email := session.Get("user_email")
    fmt.Println(user_email)
    if user_email == nil {
        w.WriteHeader(403)
    } else {
        user := new(models.User)
        ues := user.GetUEs(db, user_email.(string))
        fmt.Println(ues)
        outData, _ := json.Marshal(ues)
        w.Write(outData)
    }
}

/// TODO: Rob to understand this later
func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the register page!")
	// decoder := json.NewDecoder(r.Body)
	// data := map[string]string{"company": "", "firstName": "", "lastName": "", "email": "", "password": ""}
	// err := decoder.Decode(&data)
	// if err != nil {
	// 	panic(err)
	// }

	// db := utilities.GetDB(r)
	// user := new(models.User)

	// user.NewUser(db, data["company"], data["firstName"], data["lastName"], data["email"], data["password"])
	// fmt.Println(user)
}

/* End Of File */
