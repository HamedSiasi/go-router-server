# Message Router Server

This repository contains the web server and web client for the U-blox Test Mobile project.  The server  (residing in the service directory) is written in Golang and is based on original code provided by Neul as part of the tedI project.  The client (residing in the static directory with all the javascript/HTML code actually piled into the `static/dist` directory (using Gulp)) is written in Javascript using React/Flux and HTML5, original code created by a contractor (mmanjoura) who worked on the project.


# Server Installation

(1) Fetch and build the code from github:

`go get -u github.com/HamedSiasi/go-router-server`

(2) Copy `config.cfg` into the bin directoy:

`cp ~/code/gocode/src/github.com/u-blox/utm-server/config.cfg ~/code/gocode/bin`

(3) Edit `config.cfg` for the correct port number (default: `3001`)

(4) Edit `config.cfg` for the correct `username` nad `password` (default: `hamed`,`neulneul`)

(5) Change directory to the bin directory:

`cd ~/code/gocode/bin`

(6) Start the executable from the bin directory with:

`go-router-server`

