/* Utilities to do with user stuff for the UTM server.
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
    "errors"
    "github.com/goincremental/negroni-sessions"
    "net/http"
)

/// Return the user ID for the current session
func GetUserId(r *http.Request) (string, error) {
    session := sessions.GetSession(r)
    user_id := session.Get("user_id")
    if user_id == nil {
        return "", errors.New("No user")
    } else {
        return user_id.(string), nil
    }
}

/// TODO
func AuthenticationHandler(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, err := GetUserId(r)
        if err != nil {
            w.WriteHeader(403)
        } else {
            next(w, r)
        }
    }
}

/* End Of File */
