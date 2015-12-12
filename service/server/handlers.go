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
    "io"
    "github.com/davecgh/go-spew/spew"
    "github.com/goincremental/negroni-sessions"
    "github.com/robmeades/utm/service/globals"
    "github.com/robmeades/utm/service/models"
    "github.com/robmeades/utm/service/utilities"
    "net/http"
    "time"
    "strconv"
    "strings"
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

func LoginHandler(response http.ResponseWriter, request *http.Request) {

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

func RegisterHandler(response http.ResponseWriter, request *http.Request) {

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

/// Query the database
func QueryHandler(response http.ResponseWriter, request *http.Request) {

    var duration int64
    var startDateTime time.Time
    var endDateTime time.Time
    var err error
    
    uuid := strings.ToLower(request.FormValue("uuid"))        
    durationString := request.FormValue("duration")
    startDateTimeString := request.FormValue("startDateTime")
    
    globals.Dbg.PrintfTrace ("%s [handler] --> uuid %s, startDateTimeString %s, durationString %s.\n", globals.LogTag, uuid, startDateTimeString, durationString)
    
    if (len(uuid) > 0) {
        if len(startDateTimeString) > 0 {
            startDateTime, err = time.Parse("2006-01-02T15:04", startDateTimeString)
            if err != nil {
                globals.Dbg.PrintfTrace ("%s [handler] --> couldn't parse date/time string \"%s\".\n", globals.LogTag, startDateTimeString)
            }
        }
        
        if len(durationString) > 0 {
            duration, err = strconv.ParseInt (durationString, 10, 0)
            if err != nil {
                globals.Dbg.PrintfTrace ("%s [handler] --> couldn't parse duration string \"%s\".\n", globals.LogTag, durationString)
            } else {
                endDateTime = startDateTime.Add(time.Duration(duration) * time.Minute)
            }
        }    
        
        db := utilities.GetDB(request)
        if db != nil {
            globals.Dbg.PrintfTrace ("%s [handler] --> Querying UtmXmlData collection for uuid: %s, start date/time: %v, end date/time: %v.\n",
                globals.LogTag, uuid, startDateTime, endDateTime)
            query, err := utilities.XmlDataQuery(uuid, startDateTime, endDateTime)
            if err == nil {
                // Send the requested data
                response.Header().Set("content-type", "application/text")
                response.Header().Set("content-disposition", "attachment; filename=\"" + uuid +".txt\"");
                response.WriteHeader(http.StatusOK)
                for _, item := range *query {
                    var output string
                    output = item.Date.Format(time.RFC3339) + "\t" + item.XmlData + "\n"
                    io.WriteString(response, output)
                }
            } else {
                response.WriteHeader(http.StatusNoContent)
            }
        } else {
            response.WriteHeader(404)        
        }
    } else {
        response.WriteHeader(http.StatusNoContent)
    }
}

/// Get the summary data for the front page
func GetFrontPageData(response http.ResponseWriter, request *http.Request) *globals.Error {
    err := utilities.ValidateGetRequest(request)
    if err == nil {
        displayData := displayFrontPageData()
        if displayData != nil {
            // Send the requested data
            response.Header().Set("Content-Type", "application/json")
            response.WriteHeader(http.StatusOK)
            err := json.NewEncoder(response).Encode(displayData)
            if err != nil {
                globals.Dbg.PrintfError("%s [handler] --> received REST request %s but attempting to serialise the result:\n%s\n...yielded error %s.\n", globals.LogTag, request.URL, spew.Sdump(displayData), err.Error())
                return utilities.ServerError(err)
            }
        } else {
            return utilities.ServerError(err)            
        }
    }

    return err
}

/* End Of File */
