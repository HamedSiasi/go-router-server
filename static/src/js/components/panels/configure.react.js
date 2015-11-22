var React = require('react');
var Link = require('react-router-component').Link;

var Configure = React.createClass({
  render:function(){
    return (
   
    	<div className="col-lg-4">
            <div  style={{height: 60, width: 200, marginTop: 100}}>
              <div className="panel-body">
                <p>               
                  <Link href="/mode">
                    <b className="fa fa-cogs  fa-2x" />
                  </Link>
                  <b style={{float: 'right'}}>Configure Ticked</b>               
                </p>
              </div>
            </div>
          </div>
     
    );
  }
});

module.exports = Configure;
