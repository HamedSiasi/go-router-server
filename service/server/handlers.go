/* Web request handlers for UTM server.
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

package server

import (
    "encoding/json"
	"github.com/robmeades/utm/service/utilities"
	"github.com/robmeades/utm/service/models"
	"github.com/robmeades/utm/service/globals"
	"net/http"
    "github.com/davecgh/go-spew/spew"
    "github.com/goincremental/negroni-sessions"
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

func loginHandler(response http.ResponseWriter, request *http.Request) {

	email := request.FormValue("email")
	password := request.FormValue("password")
	session := sessions.GetSession(request)

	db := utilities.GetDB(request)
	user := new(models.User)
	err := user.Authenticate(db, email, password)
	if err == nil {
		session.Set("user_id", user.ID.Hex())
		session.Set("user_company", user.Company)
		session.Set("user_email", user.Email)
		//Fmt.Fprintf(response, "User %s success!\n", session.Get("user_email"))
		http.Redirect(response, request, "/", 302)
	}
}

func registerHandler(response http.ResponseWriter, request *http.Request) {

	company := request.FormValue("company_name")
	firstName := request.FormValue("user_firstName")
	lastName := request.FormValue("user_lastName")
	email := request.FormValue("email")
	password := request.FormValue("password")

	db := utilities.GetDB(request)
	user := new(models.User)

	user.NewUser(db, company, firstName, lastName, email, password)
	//Fmt.Fprintf(response, "User %s created successfully!\n", firstName)
	http.Redirect(response, request, "/login", 302)
}

/// Get the summary data for the front page
func getFrontPageData (response http.ResponseWriter, request *http.Request) *globals.Error {
	err := utilities.ValidateGetRequest (request)
	if err == nil {
        displayData := displayFrontPageData()
    	// Send the requested data
    	response.Header().Set("Content-Type", "application/json")
    	response.WriteHeader(http.StatusOK)
    	err := json.NewEncoder(response).Encode(displayData)
    	if err != nil {
    		globals.Dbg.PrintfError("%s [handler] --> received REST request %s but attempting to serialise the result:\n%s\n...yielded error %s.\n", globals.LogTag, request.URL, spew.Sdump(displayData), err.Error())
    		return utilities.ServerError(err)
    	}
    }	

	return err
}

/* End Of File */