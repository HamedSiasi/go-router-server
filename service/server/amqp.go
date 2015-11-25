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

// Generic struct describing an AMQP message
type AmqpMessage struct {
    DeviceUuid   string `bson:"device_uuid" json:"device_uuid"`
    EndpointUuid int    `bson:"endpoint_uuid" json:"endpoint_uuid"`
    Payload      []int  `bson:"payload" json:"payload"`
}

// Struct describing an AMQP response (i.e. uplink) message
type AmqpResponseMessage struct {
    AmqpMessage
    DeviceName string `bson:"device_name" json:"device_name"`
    Command    string `bson:"command" json:"command"`
    Data       []byte `bson:"data" json:"data"`
}

// Struct describing an AMQP error message
type AmqpErrorMessage struct {
    AmqpMessage
    Queue   string `bson:"queue" json:"queue"`
    Message string `bson:"message" json:"message"`
    Reason  string `bson:"reason" json:"reason"`
}

// Struct describing an AMQP queue
type Queue struct {
    Conn     *amqp.Connection
    Quit     chan interface{}
    Msgs     chan interface{}
    Downlink chan AmqpMessage
}


// Set up a channel to consume a queue
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

/// Open an AMQP and loop around UL and DL message processing.
// Note that the UL and DL message processing is in separate loops,
// otherwise they can block each other.
func OpenQueue(username, amqpAddress string) (*Queue, error) {
    q := Queue{}
    q.Quit = make(chan interface{})
    q.Msgs = make(chan interface{})
    q.Downlink = make(chan AmqpMessage)

    Conn, err := amqp.Dial(amqpAddress)
    if err != nil {
        return nil, fmt.Errorf("--> connecting to RabbitMQ: \"%s\"", err.Error())
    }
    q.Conn = Conn

    channel, err := Conn.Channel()
    if err != nil {
        q.Close()
        return nil, err
    }

    responseChan, err := consumeQueue(channel, username+".response")
    if err != nil {
        q.Close()
        return nil, err
    }
    errorChan, err := consumeQueue(channel, username+".error")
    if err != nil {
        q.Close()
        return nil, err
    }

    //Open DB
    //amqpSession, err := mgo.Dial("127.0.0.1:27017")

    // rec_Collection := amqpSession.DB("utm-db").C("receivemessages")
    // res_Collection := amqpSession.DB("utm-db").C("responsemessages")
    // err_Collection := amqpSession.DB("utm-db").C("errormessages")

    // Uplink loop
    go func() {
        defer func() {
            // Close the downlinkMessages channel when
            //the AMQP handler is closed down and set it to nil
            downlinkMessages = nil
            close(q.Downlink)
            //amqpSession.Close()
        }()
        
        // Continually process AMQP UL messages until commanded to Quit
        for {
            var msg amqp.Delivery
            var ok bool
            receivedMsg := false
	        Dbg.PrintfTrace("%s --> AMQP waiting for UL stimulus.\n", logTag)
            select {
	            case <-q.Quit:
        	        Dbg.PrintfTrace("%s --> AMQP UL quitting.\n", logTag)
	                return
	            case msg, ok = <-responseChan:
	                if (!ok) {
	                    return
	                }
	                receivedMsg = true
	                m := AmqpResponseMessage{}
	                err = json.Unmarshal(msg.Body, &m)
	                if err == nil {
	                    q.Msgs <- &m
	                    Dbg.PrintfTrace("%s --> UTM UUID is %+v.\n", logTag, m.DeviceUuid)
	                    Dbg.PrintfTrace("%s --> UTM name is \"%+v\".\n", logTag, m.DeviceName)
	                }
	            case msg, ok = <-errorChan:
	                if (!ok) {
	                    return
	                }
	                receivedMsg = true
        	        Dbg.PrintfTrace("%s --> AMQP error channel says %v.\n", logTag, msg)
	                m := AmqpErrorMessage{}
	                err = json.Unmarshal(msg.Body, &m)
	                if err == nil {
    					Dbg.PrintfTrace("%s--> error is %v.\n", logTag, m)
	                }
    					
            }
            
            if receivedMsg {
                if err == nil {
                    Dbg.PrintfTrace("%s --> received UL:\n\n%+v\n\n", logTag, string(msg.Body))
                } else {
                    Dbg.PrintfTrace("%s --> received UL:\n\n%+v\n\n...which is undecodable: \"%s\".\n", logTag, string(msg.Body), err.Error())
                }
            }
        }
    }()
    
    // Downlink loop
    go func() {        
        // Continually process DL AMQP Uplink messages until commanded to Quit
        for {
	        Dbg.PrintfTrace("%s --> AMQP waiting for DL stimulus.\n", logTag)
            select {
	            case <-q.Quit:
        	        Dbg.PrintfTrace("%s --> AMQP DL loop quitting.\n", logTag)
	                return
	            case dlMsg, ok := <-q.Downlink:
	                if (!ok) {
	                    return
	                }
	                serialisedData, err := json.Marshal(dlMsg)
	                if err != nil {
	                    Dbg.PrintfError("%s --> attempting to JSONify DL AMQP message:\n\n%+v\n\n...results in error: \"%s\".\n", logTag, dlMsg, err.Error())
	                } else {
	                    publishedMsg := amqp.Publishing{
	                        DeliveryMode: amqp.Persistent,
	                        Timestamp:    time.Now(),
	                        ContentType:  "application/json",
	                        Body:         serialisedData,
	                    }
	                    err = channel.Publish(username, "send", false, false, publishedMsg)
	                    if err != nil {
	                        Dbg.PrintfError("%s --> unable to publish DL message:\n\n%+v\n\n...due to error: \"%s\".\n", logTag, publishedMsg, err.Error())
	                    } else {
	                        Dbg.PrintfTrace("%s --> published DL message:\n\n%+v\n\n%s\n", logTag, publishedMsg, string(serialisedData))
	                    }
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
