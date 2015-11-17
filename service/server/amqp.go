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
    "log"
    "time"

    //"gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
)

var uuidMap = make(map[string]*DisplayRow)
var uuidSlice = make([]*DisplayRow, len(uuidMap))

type AmqpMessage struct {
    DeviceUuid   string `bson:"DeviceUuid" json:"DeviceUuid"`
    EndpointUuid int    `bson:"EndpointUuid" json:"EndpointUuid"`
    Payload      []int  `bson:"Payload json:"Payload"`
}

type AmqpReceiveMessage struct {
    AmqpMessage
    DeviceName string `bson:"DeviceName" json:"DeviceName`
    Id         string `bson:"Id" json:"Id"`
}

type AmqpResponseMessage struct {
    AmqpMessage
    DeviceName string `bson:"DeviceName" json:"DeviceName`
    Command    string `bson:"Command" json:"Command"`
    Data       []byte `bson:"Data" json:"Data"`
}

type AmqpErrorMessage struct {
    AmqpMessage
    Queue   string `bson:"Queue" json:"Queue"`
    Message string `bson:"Message" json:"Message"`
    Reason  string `bson:"Reason" json:"Reason"`
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
                log.Printf("%s -->maybe there is an error here...\n\n%+v\n\n", logTag, &m)
                if err == nil {
                    //rec_Collection.Insert(&m)
                    log.Printf("%s --> sending the following downlink message:\n\n%+v\n\n", logTag, &m)
                    q.Msgs <- &m
                }
            case msg = <-responseChan:
                receivedMsg = true
                m := AmqpResponseMessage{}
                err = json.Unmarshal(msg.Body, &m)
                if err == nil {
                    q.Msgs <- &m
                    log.Printf("%s --> UTM UUID is %+v.\n", logTag, m.DeviceUuid)
                    log.Printf("%s --> UTM name is '%+v'.\n", logTag, m.DeviceName)
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
                    log.Printf("%s --> attempting to JSONify AMQP message:\n\n%+v\n\n...results in error: %s.\n", logTag, dlMsg, err.Error())
                } else {
                    publishedMsg := amqp.Publishing{
                        DeliveryMode: amqp.Persistent,
                        Timestamp:    time.Now(),
                        ContentType:  "application/json",
                        Body:         serialisedData,
                    }
                    err = channel.Publish(username, "send", false, false, publishedMsg)
                    if err != nil {
                        log.Printf("%s --> unable to pulbish downlink message %+v due to error: %s.\n", logTag, publishedMsg, err.Error())
                    } else {
                        log.Printf("%s --> published downlink message %+v %s.\n", logTag, publishedMsg, string(serialisedData))
                    }
                }
            }
            if receivedMsg {
                if err == nil {
                    log.Printf("%s --> received %+v.\n", logTag, string(msg.Body))
                } else {
                    log.Printf("%s --> received %+v which is undecodable: %s.\n", logTag, string(msg.Body), err.Error())
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
