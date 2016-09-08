# Introduction to Message Router Server

This repo contains a broker server for the U-blox Test Mobile project. The server is written in Golang and is based on original utm-server projuect (`github.com/u-blox/utm-server`) provided by u-blox. 
The server receive AMQP message (which contains the "CoAP package" with all options and payload) from the client (`github.com/HamedSiasi/mbed-os-ublox-coap`) and route them to the CoAP server mentioned in CoAP package. it has also a "non-database memory" to deliver the CoAP server reply to the correct device by mapping "DeviceUUID" to the "AMQPmsgID".

# Prerequisites
(1) Download a golang binary release suitable for your system `https://golang.org/dl`

(2) Follow the installation instructions `https://golang.org/doc/install`


# Working directory structure

|_ bin

|_ pkg

|_ src _

        |_ github.com
        
        


# Server Installation

(1) Fetch and build the code from github:

`go get -u github.com/HamedSiasi/go-router-server`

(2) Copy `config.cfg` into the bin directoy

(3) Edit `config.cfg` for the correct port number (default: `3001`)

(4) Edit `config.cfg` for the correct `username` nad `password` (default: `hamed`,`neulneul`)

(5) Change directory to the bin directory

(6) Start the executable from the bin directory with:

`go-router-server`

