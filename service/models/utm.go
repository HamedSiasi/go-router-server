/* Description of UTMs for the UTM server.
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

package models

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
)

// TODO Rob to understand this later
type UtmMsg struct {
    id      bson.ObjectId `bson:"_id,omitempty" json:"id"`
  	date    time.Time     `bson:"date" json:"date"`
    uuid    string        `bson:"uuid" json:"uuid"`
    msg     string        `bson:"msg" json:"msg"`
}

// Representations of the AMQP comms
type AmqpMessage struct {
    deviceUuid   string `bson:"device_uuid" json:"device_uuid"`
    endpointUuid int    `bson:"endpoint_uuid" json:"endpoint_uuid"`
    payload      []byte `bson:"payload json:"payload"`
}

type AmqpReceiveMessage struct {
    AmqpMessage
    deviceName string `bson:"device_name" json:"device_name`
    id         string `bson:"id" json:"id"`
}

type AmqpResponseMessage struct {
    AmqpMessage
    deviceName string `bson:"device_name" json:"device_name`
    command    string `bson:"command" json:"command"`
    data       []byte `bson:"data" json:"data"`
}

type AmqpErrorMessage struct {
    AmqpMessage
    queue   string `bson:"queue" json:"queue"`
    message string `bson:"message" json:"message"`
    reason  string `bson:"reason" json:"reason"`
}

func (u *UtmMsg) Insert(db *mgo.Database, uuid string, msg string) {

	Msg := UtmMsg{}
	utmColl := db.C("UtmMsgs")
	Msg.id = bson.NewObjectId()
	Msg.date = time.Now()
	Msg.uuid = uuid
	Msg.msg = msg
	utmColl.Insert(&Msg)

}

/* End Of File */
