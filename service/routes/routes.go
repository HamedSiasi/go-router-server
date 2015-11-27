/* How to get to things on the UTM server.
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

package routes

import (
	"github.com/gorilla/mux"
	"github.com/robmeades/utm/service/controllers"
	"github.com/robmeades/utm/service/utilities"
)

func LoadRoutes() *mux.Router {

	index := controllers.Index{}
	auth := controllers.Auth{}
	user := controllers.User{}
 
	router := mux.NewRouter()

	router.HandleFunc("/", index.Welcome).Methods("GET")
	router.HandleFunc("/sitemap", index.Sitemap)

	a := router.PathPrefix("/auth").Subrouter()
	//a.HandleFunc("/login", auth.Login).Methods("POST")
	//a.HandleFunc("/logout", auth.Logout).Methods("GET")
	a.HandleFunc("/user", auth.User).Methods("GET")
	a.HandleFunc("/users", auth.Users).Methods("GET")
	a.HandleFunc("/uuids", auth.UserUEs).Methods("GET")
	//a.HandleFunc("/register", auth.Register).Methods("POST")

	a = router.PathPrefix("/users").Subrouter()
	a.HandleFunc("/profile", utilities.AuthenticationHandler(user.Profile)).Methods("GET")
	a.HandleFunc("/{id}", utilities.AuthenticationHandler(user.Get)).Methods("GET")

	return router
}

/* End Of File */
