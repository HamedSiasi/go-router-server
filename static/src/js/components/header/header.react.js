var React = require('react');
var Link = require('react-router-component').Link;

var Header = React.createClass({
  render: function() {
    return (

      <nav className="navbar navbar-default navbar-static-top" role="navigation" style={{marginBottom: 0}}>
        <div className="navbar-header">               
          <a className="brand " href="#/utmlist">
            <img src="static/dist/assets/images/logo.png" alt="u-blox" style={{maxWidth: 130, padding: 5}} />
          </a>
        </div>
        {/* /.navbar-header */}
        <ul className="nav navbar-top-links navbar-right">
          <li className="dropdown">
            <a className="dropdown-toggle" data-toggle="dropdown">
              <i className="fa fa-user fa-fw" /> User <i className="fa fa-caret-down" />
            </a>
            <ul className="dropdown-menu dropdown-user">
              <li><a href="#"><i className="fa fa-user fa-fw" /> Add User</a>
              </li>
              <li><a href="#"><i className="fa fa-user fa-fw" /> User Profile</a>
              </li>
              <li className="divider" />
              <li><a href="login.html"><i className="fa fa-sign-out fa-fw" /> Logout</a>
              </li>
              <li><a href="login.html"><i className="fa fa-sign-in fa-fw" /> Login</a>
              </li>
            </ul>
            {/* /.dropdown-user */}
          </li>
          {/* /.dropdown */}
          <li className="dropdown">
            <a className="dropdown-toggle" data-toggle="dropdown">
              <i className="fa fa-cloud-download" /> Setting  <i className="fa fa-caret-down" />
            </a>
            <ul className="dropdown-menu dropdown-user">
              <li><Link href="/mode"><i className="fa fa-exchange" /> Mode</Link></li>
              <li><Link href="/"><i className="fa fa-tachometer" /> Dash Board</Link></li>
              <li><a href="#"><i className="fa fa-file-archive-o" /> Upload Files</a></li>
              <li><a href="#"><i className="fa fa-file-excel-o" /> Download Files</a></li>
            </ul>
            {/* /.dropdown-user */}
          </li>
          {/* /.dropdown */}
          <li className="dropdown">
            <a className="dropdown-toggle" data-toggle="dropdown">
              <i className="fa fa-bar-chart" /> Reports  <i className="fa fa-caret-down" />
            </a>
            <ul className="dropdown-menu dropdown-user">
              <li><a href="#"><i className="fa fa-area-chart" /> Frame Loss</a>
              </li><li><a href="#"><i className="fa fa-bar-chart" /> Energy</a>
              </li>
            </ul>
            {/* /.dropdown-user */}
          </li>
          <li className="dropdown">
            <a className="dropdown-toggle" data-toggle="dropdown">
              <i className="fa fa-graduation-cap" /> Help  <i className="fa fa-caret-down" />
            </a>
            <ul className="dropdown-menu dropdown-user">
              <li><a href="#"><i className="fa fa-book">
                  </i> User Manual</a>
              </li><li><a href="#">
                  <i className="fa fa-question-circle" /> FAQ</a>
              </li>
            </ul>
            {/* /.dropdown-user */}
          </li>
          {/* /.dropdown-user */}
          {/* /.dropdown */}
          {/* /.dropdown */}
        </ul>
        {/* /.navbar-top-links */}
      </nav>
    );
  }
});

module.exports = Header;