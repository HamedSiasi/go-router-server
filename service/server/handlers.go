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
	"github.com/davecgh/go-spew/spew"
	"github.com/goincremental/negroni-sessions"
	"github.com/robmeades/utm/service/globals"
	"github.com/robmeades/utm/service/models"
	"github.com/robmeades/utm/service/utilities"
	"net/http"
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

func LogoutHandler(request http.ResponseWriter, response *http.Request) {
	session := sessions.GetSession(response)
	user_id := session.Get("user_id")

	globals.Dbg.PrintfTrace("Calling Logout User_id session \n%s\n", spew.Sdump(user_id))

	if user_id == nil {
		request.WriteHeader(403)
		http.Redirect(request, response, "/", 403)
	} else {
		session.Delete("user_id")
		http.Redirect(request, response, "/", 202)
	}
}

func ShowDisplayHandler(request http.ResponseWriter, response *http.Request) {
	session := sessions.GetSession(response)
	user_id := session.Get("user_id")

	globals.Dbg.PrintfTrace("Calling ShowDisplay User_id session \n%s\n", spew.Sdump(user_id))
	if user_id == nil {
		request.WriteHeader(403)
		http.Redirect(request, response, "/", 403)
	} else {
		http.Redirect(request, response, "/display", 202)
	}
}

func ShowModeHandler(request http.ResponseWriter, response *http.Request) {
	session := sessions.GetSession(response)
	user_id := session.Get("user_id")
	globals.Dbg.PrintfTrace("Calling ShowMode User ID session \n%s\n", spew.Sdump(user_id))
	if user_id == nil {
		request.WriteHeader(403)
		http.Redirect(request, response, "/", 403)
	} else {
		http.Redirect(request, response, "/mode", 202)
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {

	email := request.FormValue("email")
	password := request.FormValue("password")
	session := sessions.GetSession(request)

	db := utilities.GetDB(request)
    if db != nil {
    	user := new(models.User)
    	err := user.Authenticate(db, email, password)
    	globals.Dbg.PrintfTrace("Calling login User session \n%s\n", spew.Sdump(user.ID))
    	if err == nil {
    		session.Set("user_id", user.ID.Hex())
    		session.Set("user_company", user.Company)
    		session.Set("user_email", user.Email)
    		http.Redirect(response, request, "/display", 302)
    	} else {
    		http.Redirect(response, request, "/", 302)
    
    	}
    } else {
        response.WriteHeader(404)        
    }        	
}

func registerHandler(response http.ResponseWriter, request *http.Request) {

	company := request.FormValue("company_name")
	firstName := request.FormValue("user_firstName")
	lastName := request.FormValue("user_lastName")
	email := request.FormValue("email")
	password := request.FormValue("password")

	db := utilities.GetDB(request)
    if db != nil {
    	user := new(models.User)
    
    	user.NewUser(db, company, firstName, lastName, email, password)
    
    	http.Redirect(response, request, "/display", 302)
    } else {
        response.WriteHeader(404)        
    }        	
}

/// Get the summary data for the front page
func getFrontPageData(response http.ResponseWriter, request *http.Request) *globals.Error {
	err := utilities.ValidateGetRequest(request)
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
