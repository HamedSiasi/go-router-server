/* AMQP interface definitions for the UTM server.
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
    "fmt"
    "github.com/streadway/amqp"
    "time"

    //"gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
)

var uuidMap = make(map[string]*DisplayRow)
var uuidSlice = make([]*DisplayRow, len(uuidMap))

type AmqpMessage struct {
    DeviceUuid   string `bson:"device_uuid" json:"device_uuid"`
    EndpointUuid int    `bson:"endpoint_uuid" json:"endpoint_uuid"`
    Payload      []int  `bson:"payload json:"payload"`
}

type AmqpReceiveMessage struct {
    AmqpMessage
    DeviceName string `bson:"device_name" json:"device_name"`
    Id         string `bson:"id" json:"id"`
}

type AmqpResponseMessage struct {
    AmqpMessage
    DeviceName string `bson:"device_name" json:"device_name"`
    Command    string `bson:"command" json:"command"`
    Data       []byte `bson:"data" json:"data"`
}

type AmqpErrorMessage struct {
    AmqpMessage
    Queue   string `bson:"queue" json:"queue"`
    Message string `bson:"message" json:"message"`
    Reason  string `bson:"reason" json:"reason"`
}

type Queue struct {
    Conn     *amqp.Connection
    Quit     chan interface{}
    Msgs     chan interface{}
    Downlink chan AmqpMessage
}

var totalMsgs uint64
var totalBytes uint64

var row = &DisplayRow{}
var rowsList = make([]*DisplayRow, 1)

func consumeQueue(channel *amqp.Channel, chanName string) (<-chan amqp.Delivery, error) {
    msgChan, err := channel.Consume(
        chanName, // Queue
        "",       // consumer
        true,     // auto-ack
        false,    // exclusive
        false,    // no-local
        false,    // no-wait
        nil,      // args
    )
    if err != nil {
        return nil, fmt.Errorf("--> consuming '%s' channel: %s", chanName, err.Error())
    }
    return msgChan, nil
}

func OpenQueue(username, amqpAddress string) (*Queue, error) {
    q := Queue{}
    q.Quit = make(chan interface{})
    q.Msgs = make(chan interface{})
    q.Downlink = make(chan AmqpMessage)

    Conn, err := amqp.Dial(amqpAddress)
    if err != nil {
        return nil, fmt.Errorf("--> connecting to RabbitMQ: %s", err.Error)
    }
    q.Conn = Conn

    channel, err := Conn.Channel()
    if err != nil {
        return nil, err
    }

    receiveChan, err := consumeQueue(channel, username+".receive")
    if err != nil {
        return nil, err
    }
    responseChan, err := consumeQueue(channel, username+".response")
    if err != nil {
        return nil, err
    }
    errorChan, err := consumeQueue(channel, username+".error")
    if err != nil {
        return nil, err
    }

    //Open DB
    //amqpSession, err := mgo.Dial("127.0.0.1:27017")

    // rec_Collection := amqpSession.DB("utm-db").C("receivemessages")
    // res_Collection := amqpSession.DB("utm-db").C("responsemessages")
    // err_Collection := amqpSession.DB("utm-db").C("errormessages")

    if err != nil {
        panic(err)
    }

    //var m AmqpMessage
    go func() {
        defer func() {
            // Close the downlinkMessages channel when
            //the AMQP handler is closed down and set it to nil
            downlinkMessages = nil
            close(q.Downlink)
            //amqpSession.Close()
        }()
        // Continually process AMQP messages until commanded to Quit
        for {
            var msg amqp.Delivery
            receivedMsg := false
            select {
	            case <-q.Quit:
	                return
	            case msg = <-receiveChan:
	                receivedMsg = true
	                m := AmqpReceiveMessage{}
	                err = json.Unmarshal(msg.Body, &m)
	                if err == nil {
	                    //rec_Collection.Insert(&m)
	                    Dbg.PrintfTrace("%s --> sending the following downlink message:\n\n%+v\n\n", logTag, &m)
	                    q.Msgs <- &m
	                } else {
		                Dbg.PrintfError("%s -->maybe there is an error here...\n\n%+v\n\n", logTag, &m)	                	
	                }
	            case msg = <-responseChan:
	                receivedMsg = true
	                m := AmqpResponseMessage{}
	                err = json.Unmarshal(msg.Body, &m)
	                if err == nil {
	                    q.Msgs <- &m
	                    Dbg.PrintfTrace("%s --> UTM UUID is %+v.\n", logTag, m.DeviceUuid)
	                    Dbg.PrintfTrace("%s --> UTM name is '%+v'.\n", logTag, m.DeviceName)
	                    row.Uuid = m.DeviceUuid
	                    row.UnitName = m.DeviceName
	                    //res_Collection.Insert(&m)
	                }
	            case msg = <-errorChan:
	                receivedMsg = true
	                m := AmqpErrorMessage{}
	                err = json.Unmarshal(msg.Body, &m)
	                if err == nil {
	                    //err_Collection.Insert(&m)
	                    q.Msgs <- &m
	                }
	            case dlMsg := <-q.Downlink:
	                serialisedData, err := json.Marshal(dlMsg)
	                if err != nil {
	                    Dbg.PrintfError("%s --> attempting to JSONify AMQP message:\n\n%+v\n\n...results in error: %s.\n", logTag, dlMsg, err.Error())
	                } else {
	                    publishedMsg := amqp.Publishing{
	                        DeliveryMode: amqp.Persistent,
	                        Timestamp:    time.Now(),
	                        ContentType:  "application/json",
	                        Body:         serialisedData,
	                    }
	                    err = channel.Publish(username, "send", false, false, publishedMsg)
	                    if err != nil {
	                        Dbg.PrintfError("%s --> unable to publish downlink message:\n\n%+v\n\n...due to error: %s.\n", logTag, publishedMsg, err.Error())
	                    } else {
	                        Dbg.PrintfTrace("%s --> published downlink message:\n\n%+v\n\n%s\n", logTag, publishedMsg, string(serialisedData))
	                    }
	                }
            }
            if receivedMsg {
                if err == nil {
                    Dbg.PrintfTrace("%s --> received:\n\n%+v\n\n", logTag, string(msg.Body))
                } else {
                    Dbg.PrintfTrace("%s --> received:\n\n%+v\n\n...which is undecodable: %s.\n", logTag, string(msg.Body), err.Error())
                }
            }
        }
    }()
    return &q, nil
}

func (q *Queue) Close() {
    // FIXME: maybe dont need a Quit chan... can get EOF info from AMQP chans somehow.
    q.Conn.Close()
    close(q.Quit)
}

/* End Of File */
