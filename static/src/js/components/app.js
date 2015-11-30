var React = require('react');
var Template = require('./app-template.js');
var Router = require('react-router-component');
var Display = require('./display/display.react');
var Mode = require('./mode/mode.react');
var Login = require('./user/login.react');
var Register = require('./user/register.react');
var Index = require('./index');

var Locations = Router.Locations;
var Location  = Router.Location;

var App = React.createClass({
    render:function(){
        return (
            <Template>
                <Locations>
                    <Location path="/" handler={Index} />
                    <Location path="#/display" handler={Display} />
                    <Location path="#/mode" handler={Mode} />
                    <Location path="#/login" handler={Login} />
                    <Location path="#/register" handler={Register} />
                </Locations>
            </Template>
        );
    }
});

module.exports = App;