/* Utilities to do with GET/POST etc. behaviours for the UTM server.
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
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/u-blox/utm/service/globals"
)

func ValidateGetRequest (request *http.Request) *globals.Error {
    // Ensure this is a GET request
    if (request.Method != "GET") || (request.Method == "") {
        globals.Dbg.PrintfError("%s [dl_msgs] --> received unsupported REST request %s %s.\n", globals.LogTag, request.Method, request.URL)
        return ClientError("unsupported method", http.StatusBadRequest)
    }
    
    return nil    
}

func ValidatePostRequest (request *http.Request) *globals.Error {
    // Ensure this is a POST request
    if (request.Method != "POST") || (request.Method == "") {
        globals.Dbg.PrintfError("%s [dl_msgs] --> received unsupported REST request %s %s.\n", globals.LogTag, request.Method, request.URL)
        return ClientError("unsupported method", http.StatusBadRequest)
    }
    
    return nil    
}

/// Definition of a REST client
type RestClient struct {
    address string
    client  http.Client
}

/// Create a new REST client with a given address
func NewRestClient(address string) *RestClient {
    var c RestClient
    c.address = address
    return &c
}

/// Perform a GET request on a given URL for a given REST client
func (c *RestClient) Get(url string) (*http.Response, error) {
    resp, err := c.client.Get(c.address + url)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

/// Perform a GET request on a given URL for a given REST
// client and the convert the result into JSON.
// The obtained JSON is deserialised into v
func (c *RestClient) GetJSON(url string, v interface{}) error {
    resp, err := c.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        DumpResponse(resp)
        return fmt.Errorf("Invalid response from %s: %s", url, resp.Status)
    }
    err = json.NewDecoder(resp.Body).Decode(v)
    if err != nil {
        return err
    }
    return nil
}

/// Perform a POST request on a given URL for a given REST
// client, sending v as JSON.
func (c *RestClient) Post(url string, v interface{}) (*http.Response, error) {
    var buf bytes.Buffer
    err := json.NewEncoder(&buf).Encode(v)
    if err != nil {
        return nil, err
    }
    resp, err := c.client.Post(c.address+url, "application/json", &buf)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

/// Perform a POST request on a given URL for a given REST
// client, sending v as JSON and returning nil if successful,
// otherwise an error string.
func (c *RestClient) PostOK(url string, v interface{}) error {
    resp, err := c.Post(url, v)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        DumpResponse(resp)
        return fmt.Errorf("Invalid response from %s: %s", url, resp.Status)
    }
    return nil
}

/// Perform a PUT request on a given URL for a given REST
// client, sending v as JSON and returning nil if successful,
// otherwise an error string.
func (c *RestClient) PutOK(url string, v interface{}) error {
    var buf bytes.Buffer
    err := json.NewEncoder(&buf).Encode(v)
    if err != nil {
        return err
    }
    request, err := http.NewRequest("PUT", url, &buf)
    if err != nil {
        return err
    }
    request.Header.Set("Content-Type", "application/json")
    resp, err := c.client.Do(request)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        DumpResponse(resp)
        return fmt.Errorf("Invalid response from %s: %s", url, resp.Status)
    }
    return nil
}

/// TODO
// The interface req (if not nil) is serialised as JSON content in the
// delete request and any response content is deserialised into the
// interface res (if not nil)
func (c *RestClient) DeleteOK(url string) error {
    // Construct a DELETE request
    request, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }

    // Issue the DELETE request and check the response
    response, err := c.client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        DumpResponse(response)
        return fmt.Errorf("Invalid response from %s: %s", url, response.Status)
    }
    return nil
}

/* End Of File */
