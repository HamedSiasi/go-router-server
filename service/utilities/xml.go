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
)

type UtmXml struct {
	Id      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Date    time.Time     `bson:"date" json:"date"`
	Uuid    string        `bson:"uuid" json:"uuid"`
	XmlData string        `bson:"xmldata" json:"xmldata"`
}

// TODO Rob to understand this later
func (u *UtmXml) Insert(db *mgo.Database, xmlData string, uuid string) {

	var utmXml = UtmXml{}

	xml := db.C("UtmsXml")
	u.Id = bson.NewObjectId()
	u.Date = time.Now()
	u.Uuid = uuid
	u.XmlData = xmlData
	xml.Insert(&utmXml)

}

// TODO Rob to understand this later
func (u *UtmXml) Get(db *mgo.Database, uuid string) error {

	return db.C("UtmsXml").FindId(bson.ObjectIdHex(uuid)).One(&u)

}

/* End Of File */
