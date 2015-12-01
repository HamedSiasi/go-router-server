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
    "github.com/robmeades/utm/service/globals"
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
            globals.Dbg.PrintfTrace("%s [amqp] --> AMQP waiting for UL stimulus.\n", globals.LogTag)
            select {
                case <-q.Quit:
                    globals.Dbg.PrintfTrace("%s [amqp] --> AMQP UL quitting.\n", globals.LogTag)
                    return
                case msg, ok = <-responseChan:
                    if (!ok) {
                        globals.Dbg.PrintfTrace("%s [amqp] --> responseChan got a NOT OK thing, returning.\n", globals.LogTag)
                        // TODO let's see if it recovers (was return)
                    }
                    receivedMsg = true
                    m := AmqpResponseMessage{}
                    err = json.Unmarshal(msg.Body, &m)
                    if err == nil {
                        globals.Dbg.PrintfTrace("%s [amqp] --> UTM UUID is %+v.\n", globals.LogTag, m.DeviceUuid)
                        globals.Dbg.PrintfTrace("%s [amqp] --> UTM name is \"%+v\".\n", globals.LogTag, m.DeviceName)
                        q.Msgs <- &m
                    }
                case msg, ok = <-errorChan:
                    if (!ok) {
                        globals.Dbg.PrintfTrace("%s [amqp] --> errorChan got a NOT OK thing, returning.\n", globals.LogTag)
                        // TODO let's see if it recovers (was return)
                    }
                    receivedMsg = true
                    globals.Dbg.PrintfTrace("%s [amqp] --> AMQP error channel says %v.\n", globals.LogTag, msg)
                    m := AmqpErrorMessage{}
                    err = json.Unmarshal(msg.Body, &m)
                    if err == nil {
                        globals.Dbg.PrintfTrace("%s--> error is %v.\n", globals.LogTag, m)
                    }
                        
            }
            
            if receivedMsg {
                if err == nil {
                    globals.Dbg.PrintfTrace("%s [amqp] --> received UL:\n\n%+v\n\n", globals.LogTag, string(msg.Body))
                } else {
                    globals.Dbg.PrintfTrace("%s [amqp] --> received UL:\n\n%+v\n\n...which is not JSON decodable: \"%s\".\n", globals.LogTag, string(msg.Body), err.Error())
                }
            }
        }
    }()
    
    // Downlink loop
    go func() {        
        // Continually process DL AMQP Uplink messages until commanded to Quit
        for {
            globals.Dbg.PrintfTrace("%s [amqp] --> AMQP waiting for DL stimulus.\n", globals.LogTag)
            select {
                case <-q.Quit:
                    globals.Dbg.PrintfTrace("%s [amqp] --> AMQP DL loop quitting.\n", globals.LogTag)
                    return
                case dlMsg, ok := <-q.Downlink:
                    if (!ok) {
                        globals.Dbg.PrintfTrace("%s [amqp] --> q.Downlink got a NOT OK thing, returning.\n", globals.LogTag)
                        return
                    }
                    serialisedData, err := json.Marshal(dlMsg)
                    if err != nil {
                        globals.Dbg.PrintfError("%s [amqp] --> attempting to JSONify DL AMQP message:\n\n%+v\n\n...results in error: \"%s\".\n", globals.LogTag, dlMsg, err.Error())
                    } else {
                        publishedMsg := amqp.Publishing{
                            DeliveryMode: amqp.Persistent,
                            Timestamp:    time.Now(),
                            ContentType:  "application/json",
                            Body:         serialisedData,
                        }
                        err = channel.Publish(username, "send", false, false, publishedMsg)
                        if err != nil {
                            globals.Dbg.PrintfError("%s [amqp] --> unable to publish DL message:\n\n%+v\n\n...due to error: \"%s\".\n", globals.LogTag, publishedMsg, err.Error())
                        } else {
                            globals.Dbg.PrintfTrace("%s [amqp] --> published DL message:\n\n%+v\n\n%s\n", globals.LogTag, publishedMsg, string(serialisedData))
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
