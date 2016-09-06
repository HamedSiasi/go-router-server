/* XML utilites for the UTM server.
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
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
    "bytes"
    "strings"
)

type UtmXml struct {
    Id      bson.ObjectId `bson:"_id,omitempty" json:"id, omitempty"`
    Date    time.Time     `bson:"date" json:"date, omitempty"`
    Uuid    string        `bson:"uuid" json:"uuid, omitempty"`
    XmlData string        `bson:"xmldata" json:"xmldata, omitempty"`
}

type UtmXmlUuidData struct {
    Date    time.Time     `bson:"date" json:"date, omitempty"`
    XmlData string        `bson:"xmldata" json:"xmldata, omitempty"`
}

// Store a byte array of XML to the database
func XmlDataStore(data []byte, uuid string) error {

    utmXml := UtmXml{}
    utmXml.Id = bson.NewObjectId()
    utmXml.Date = time.Now()
    utmXml.Uuid = strings.ToLower(uuid)
    utmXml.XmlData = string(data[:bytes.IndexByte(data, 0)])
    
    session, err := mgo.Dial("127.0.0.1:27017")
    if err == nil {
        defer session.Close()
        collection := session.DB("utm-db").C("UtmXmlData")
        err = collection.Insert(&utmXml)
    }

    return err
}

// Query the database for XmlData
func XmlDataQuery (uuid string, startDateTime time.Time, endDateTime time.Time) (*[]UtmXmlUuidData, error) {
    
    session, err := mgo.Dial("127.0.0.1:27017")
    if err == nil {
        defer session.Close()
        collection := session.DB("utm-db").C("UtmXmlData")
        var query []UtmXmlUuidData
        if endDateTime.After(time.Time{}) {
            err = collection.Find(bson.M{"uuid": uuid, "date": bson.M{"$gte": startDateTime, "$lte": endDateTime }}).All (&query)
        } else {
            err = collection.Find(bson.M{"uuid": uuid, "date": bson.M{"$gte": startDateTime}}).All (&query)
        }
        return &query, err
    }    
    return nil, err
}

/* End Of File */
