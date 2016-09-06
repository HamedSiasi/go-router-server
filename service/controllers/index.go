/* Page-serving elements of the UTM server.
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
    "github.com/u-blox/utm/service/utilities"
    "net/http"
)

type Index struct{}
type SiteMap []map[string]string

/// Serve the main page
func (i *Index) Welcome(w http.ResponseWriter, req *http.Request) {
    http.ServeFile(w, req, "static/index.html")
}

/// TODO
func (i *Index) Sitemap(w http.ResponseWriter, req *http.Request) {
    sitemap := SiteMap{
        {"url": "#/login",
            "title": "Login"},
    }
    sitemap_for_user := SiteMap{
        {"url": "/auth/logout",
            "title": "Logout"},
        {"url": "#/users/profile",
            "title": "Profile"},
    }
    _, err := utilities.GetUserId(req)
    if err == nil {
        js, _ := json.Marshal(sitemap_for_user)
        w.Write(js)
    } else {
        js, _ := json.Marshal(sitemap)
        w.Write(js)
    }
}

/* End Of File */
