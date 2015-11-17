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
)

// Definition of a UTM
type Uuid struct {
    id      bson.ObjectId `bson:"_id,omitempty" json:"id"`
    uid     string        `bson:"uid" json:"uid"`
    rssi    string        `bson:"rssi" json:"rssi"`
    company string        `bson:"company" json:"company"`
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

/// Test data
var v = [...]Uuid{
    {uid: "861f9e8c-5b8d-11e5-885d-feff819cdc9a", rssi: "20", company: "u-blox"},
    {uid: "861fa274-5b8d-11e5-885d-feff819cdc9b", rssi: "27", company: "u-blox"},
    {uid: "861fa3fa-5b8d-11e5-885d-feff819cdc9c", rssi: "32", company: "u-blox"},
    {uid: "861fa53a-5b8d-11e5-885d-feff819cdc9d", rssi: "62", company: "vodafone"},
    {uid: "861fa670-5b8d-11e5-885d-feff819cdc9e", rssi: "72", company: "vodafone"},
    {uid: "861fa79c-5b8d-11e5-885d-feff819cdc9f", rssi: "52", company: "vodafone"},
}

// Add a UTM to the database based on UUID
func (u *Uuid) Insert(db *mgo.Database) {

    c := db.C("uuids")
    for _, e := range v {
        e.id = bson.NewObjectId()
        c.Insert(&e)
    }
}

/* End Of File */
