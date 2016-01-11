This repository contains the web server and web client for the U-blox Test Mobile project.  The server  (residing in the service directory) is written in Golang and is based on original code provided by Neul as part of the tedI project.  The client (residing in the static directory with all the javascript/HTML code actually piled into the static/dist directory (using Gulp)) is written in Javascript using React/Flux and HTML5, original code created by a contractor (mmanjoura) who worked on the project.

# Setting Up For Windows Development

There is an MS Word document in this directory which explains how to set up for development on Windows, using Eclipse as IDE.

# Server Installation

Concerning server installation, the web server is run on a VM in Sgonico (151.9.34.90, http://ciot.it-sgn.u-blox.com/), user account "itadmin"). All files are kept under the directory ~/code/gocode.

The bash script is setup to execute these commands before loading the shell:

export GOROOT=$HOME/code/golang/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=$HOME/code/gocode
export PATH=/usr/local/git/bin:$PATH

Mongo database must be installed.

Start mongo with:

sudo mongod --config /etc/mongodb.conf

[Probably shouldn't need sudo but I found that mongo couldn't write to the journal file without this].

Fetch and build the UTM code from github with:

go get -u github.com/u-blox/utm

When it has built, copy the static directory into the bin directory:

cp -r ~/code/gocode/src/github.com/u-blox/utm/static ~/code/gocode/bin

Copy config.cfg into the bin directoy:

cp ~/code/gocode/src/github.com/u-blox/utm/config.cfg ~/code/gocode/bin

Edit config.cfg for the correct port number (default 8080).

Change directory to the bin directory:

cd ~/code/gocode/bin

Start the executable from the bin directory with:

nohup ./utm &

You should now be able to browse to the server at http://ciot.it-sgn.u-blox.com:8080/

# Mongo DB And Security

In Mongo on the Sgonico server a single set of web log-in details have been created.  They are username (actually e-mail, but we never e-mail to it): "one@astellia", password: "crazy8".  If you ever need to create new ones, this can be done through the Add User menu on the wev interface.  If you ever need to disable security for any reason, find the file app_store.js in static/src/js/stores and change the line:

var _isLoggedIn = false;

to:

var _isLoggedIn = true;

No, it's not very secure.  Oh, but do make sure that Gulp is running when you make the changes or otherwise the change won't be piled into your dist folder.

The Mongo shell can be entered by typing:

mongo

Useful mongo commands are:

* Show the databases: show dbs.
* Use a database (e.g. utm-db): use utm-db.
* Show the collections in a database: show collections.
* Display the contents of a collection (e.g. users) after "use"ing the relevant database: db.users.find().
* Remove the an entire collection (e.g. users) after "use"ing the relevant database: db.users.remove({}).