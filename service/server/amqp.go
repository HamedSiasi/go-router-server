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
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/u-blox/utm/service/globals"
	"time"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

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
	Conn       *amqp.Connection
	Quit       chan interface{}
	UlAmqpMsgs chan interface{}
	DlAmqpMsgs chan AmqpMessage
}

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

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

/// Open the necessary AMQP channels and create structured
// channels to pass data up and down throught the channels,
// then set off two Go routines that loop around UL and DL
// message processing.
// Note that the UL and DL message processing is in separate loops,
// otherwise they can block each other.
func OpenQueue(username, amqpAddress string) (*Queue, error) {
	q := Queue{}
	q.Quit = make(chan interface{})
	q.UlAmqpMsgs = make(chan interface{})
	q.DlAmqpMsgs = make(chan AmqpMessage)

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

	// This is an actual Go chanel, the uplink AMQP channel
	responseChan, err := consumeQueue(channel, username+".response")
	if err != nil {
		q.Close()
		return nil, err
	}

	// No-one seems to know what this is for but you apparently
	// have to consume it
	errorChan, err := consumeQueue(channel, username+".error")
	if err != nil {
		q.Close()
		return nil, err
	}

	// Uplink loop
	go func() {
		defer func() {
			// Close the downlink message channel when
			// the AMQP handler is closed down and set it to nil
			DlMsgs = nil
			close(q.DlAmqpMsgs)
		}()

		// Continually process AMQP UL messages until commanded to Quit
		for {
			var msg amqp.Delivery
			var isOpen bool
			receivedMsg := false
			globals.Dbg.PrintfTrace("%s [amqp] --> AMQP waiting for UL stimulus.\n", globals.LogTag)
			select {
			case <-q.Quit:
				{
					globals.Dbg.PrintfTrace("%s [amqp] --> AMQP UL quitting...\n", globals.LogTag)
					return
				}
			case msg, isOpen = <-responseChan:
				{
					if isOpen {
						receivedMsg = true
						m := AmqpResponseMessage{}
						err = json.Unmarshal(msg.Body, &m)
						if err == nil {
							q.UlAmqpMsgs <- &m
						}
					} else {
						globals.Dbg.PrintfTrace("%s [amqp] --> responseChan has closed on us, returning...\n", globals.LogTag)
						// Send error on the channel so that things above us can clear up
						err = errors.New("AMQP response channel has closed.\n")
						q.UlAmqpMsgs <- &err
						return
					}
				}
			case msg, isOpen = <-errorChan:
				{
					if isOpen {
						receivedMsg = true
						globals.Dbg.PrintfTrace("%s [amqp] --> AMQP error channel says %v.\n", globals.LogTag, msg)
						m := AmqpErrorMessage{}
						err = json.Unmarshal(msg.Body, &m)
						if err == nil {
							globals.Dbg.PrintfTrace("%s--> error is %v.\n", globals.LogTag, m)
						}
					} else {
						globals.Dbg.PrintfTrace("%s [amqp] --> errorChan has closed on us, returning...\n", globals.LogTag)
						// Send error on the channel so that things above us can clear up
						q.UlAmqpMsgs <- errors.New("AMQP error channel has closed.\n")
						return
					}
				}
			}

			if receivedMsg {
				if err == nil {
					//globals.Dbg.PrintfTrace("%s --> received: %+v\n\n", globals.LogTag, string(msg.Body))
				} else {
					globals.Dbg.PrintfTrace("%s --> received: %+v\n\n...which is not JSON decodable: \"%s\".\n", globals.LogTag, string(msg.Body), err.Error())
				}
			}
		}
	}()

	// Downlink loop
	go func() {
		// Continually process AMQP DL messages until commanded to Quit
		for {
			globals.Dbg.PrintfTrace("%s [amqp] --> AMQP waiting for DL stimulus.\n", globals.LogTag)
			select {
			case <-q.Quit:
				{
					globals.Dbg.PrintfTrace("%s [amqp] --> AMQP DL loop quitting.\n", globals.LogTag)
					return
				}
			case dlMsg, isOpen := <-q.DlAmqpMsgs:
				{
					if isOpen {
						fmt.Println(dlMsg)
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
							err = channel.Publish(
								username,
								"send",
								false,
								false,
								publishedMsg)
							if err != nil {
								globals.Dbg.PrintfError("%s [amqp] --> unable to publish DL message:\n\n%+v\n\n...due to error: \"%s\".\n", globals.LogTag, publishedMsg, err.Error())
							} else {
								globals.Dbg.PrintfTrace("%s [amqp] --> published DL message:\n\n%+v\n\n%s\n", globals.LogTag, publishedMsg, string(serialisedData))
							}
						}
					} else {
						globals.Dbg.PrintfTrace("%s [amqp] --> the downlink channel has closed on us, returning...\n", globals.LogTag)
						return
					}
				}
			} //select
		}
	}()

	return &q, nil
}

/// Close the AMQP connection
func (q *Queue) Close() {
	q.Conn.Close()
	close(q.Quit)
}

/* End Of File */
