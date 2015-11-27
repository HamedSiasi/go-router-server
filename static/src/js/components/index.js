var React = require('react');


var Index = React.createClass({

	render:function(){
		return (
			  <div className="row" ><br />
                          <div className="panel panel-default" >
                            <div className="_panel-heading" style={{width:'80%'}}>
                              <div className="panel-body">
                                <div className="dataTable_wrapper">
		<div>
		        <div>
        <h2>
          <a name="_Toc425263541">1.1 Introduction</a>
        </h2>
        <p>
          The software is shipped with a default configuration, which may be re-enabled via the remote interface or USB port. It has a number of basic test modes.
          Remote adjustment of the features and configurations can only be done in Standard TRX mode, and may have immediate impact on current behaviour. Most
          changes will persist across a power cycle. The USB port, or for remote access, Standard TRX mode, is also used to configure and manage entry into either
          Receive-Only or Transmit-Only modes, though these modes do not persist across a power cycle.
        </p>
        <p>
          Changes to the configuration are recorded in the host FLASH storage and therefore survive power cycling.
        </p>
        <h3>
          <a name="_Toc425263542">1.1.1 Basic Modes of Operation:</a>
        </h3>
        <p>
          1. Self-Test
        </p>
        <p>
          2. Standard TRX
        </p>
        <p>
          3. Receive only
        </p>
        <p>
          4. Transmit only
        </p>
        <h4>
          1.1.1.1 Self-Test
        </h4>
        <p>
          At boot up the host application will perform a number of tests to verify UTM operation. These include:
        </p>
        <p>
          1. Battery Voltage above required threshold
        </p>
        <p>
          2. Module ‘OK’ through AT command dialogue
        </p>
        <p>
          3. SD Card present
        </p>
        <p>
          Any failures at this stage will result in an error light and halt further operation, except for access via USB which may be used to extend tests via direct
          interaction.
        </p>
        <p>
          On successful completion of this stage the UTM will always transition to Standard TRX mode.
        </p>
        <h4>
          1.1.1.2 Standard TRX
        </h4>
        <p>
          This is the default mode after boot up and successful self-test.
        </p>
        <p>
          In this mode, the software operates with 2 basic internal timers:
        </p>
        <p>
          1. Heartbeat timer to wake up the UTM and perform one or more measurements
        </p>
        <p>
          2. Reporting interval timer to trigger a transmission of measurement(s) as user data
        </p>
        <p>
          For each measurement (see section 2.4.2) it is possible to configure the following:
        </p>
        <p>
          1. Each measurement is defined as occurring every K heartbeats
        </p>
        <p>
          2. Enable the measurement for transmission as user data
        </p>
        <p>
          Optionally, the criteria for transmission of the measurement may be set as:
        </p>
        <p>
          1. Only report the measurement if the data hysteresis is greater than D from the last
        </p>
        <p>
          2. Only report data if above a threshold A, or below a threshold B
        </p>
        <p>
          3. As (2) but on a “one-shot” basis
        </p>
        <p>
          4. Set a maximum reporting interval of MxN heartbeats irrespective of the criteria above
        </p>
        <p>
          The reporting interval is defined as occurring every N heartbeats (default is 1).
        </p>
        <p>
          In addition to the 2 internal timer mechanisms for UTM triggered operations, the eNodeB is responsible for configuring the the UE DRX Idle Timeout (the
          time the UE stays awake after an uplink transmission in order to receive downlink messages) and the UE DRX Sleep Duration (the periodicity of UE wake-up
          for downlink reception if there is nothing to transmit on the uplink). The UE DRX Sleep Duration is defined as 2^P x Q, where Q is ~5 seconds at MCS4, and
          ~320ms at MCS2, where P is in the range 0 to 15 up to a maximum period of ~2days at MCS4.
        </p>
        <p>
          It is possible to configure the UTM to request continuous uplink grants that will be filled with a numerically incrementing data stream. A timer may be set
          to terminate this function, or by remote access.
        </p>
        <h4>
          1.1.1.3 Receive Only
        </h4>
        <p>
          This mode can be pre-configured and entered via the Standard TRX mode. In this mode the module will enable its receiver and either
        </p>
        <p>
          1. perform a full band scan, after which it will return to Standard TRX mode, or
        </p>
        <p>
          2. perform a partial scan based on provided channel numbers
        </p>
        <p>
          3. perform a continuous RSSI measurement at the pre-configured frequency. In the latter case the duration of the RSSI measurement can be pre-configured to
          force a return to Standard TRX mode, or remain permanently in this mode until power down.
        </p>
        <p>
          Measurements will be logged to the SD Card if enabled, but none will be transmitted as user data whilst in this mode. The UTM will not be in communication
          with the eNodeB during this time. In the case of full band scan, the UTM can optionally transmit the full RSSI results as user data when returning to
          Standard TRX mode.
        </p>
        <p>
          Note that a full band scan can take a significant length of time (&gt;1hr) and in RX Mode power consumption will be high and battery life will be much
          shorter than Standard TRX operation. The UTM will terminate this mode when the battery voltage is measured at the minimum threshold, and after reporting
          the battery low condition in Standard TRX mode it will halt by asserting a battery low condition on the status LED.
        </p>
        <h4>
          1.1.1.4 Transmit Only
        </h4>
        <p>
          This mode can be pre-configured and entered via the Standard TRX mode. In this mode the module will enable its transmitter and either
        </p>
        <p>
          1. Transmit for a pre-configured time period, after which it will return to Standard TRX mode, or
        </p>
        <p>
          2. Transmit permanently in this mode until power down
        </p>
        <p>
          The transmission is continuous and may be used to measure radiation profiles. This mode should only be used in the lab to avoid spurious transmissions on
          licensed spectrum in the field. Note that in Standard TRX mode it is possible for the UTM to be configured to ask for continuous uplink grants.
        </p>
        <p>
          Measurements will be logged to the SD Card if enabled, but only pre-configured RF Burst content will be transmitted and the UTM will not be in
          communication with the eNodeB during this time.
        </p>
        <p>
          Note that power consumption will be significant during this time and battery life will be much shorter than Standard TRX operation. The UTM will terminate
          this mode when the battery voltage is measured at the minimum threshold, and after reporting the battery low condition in Standard TRX mode it will halt by
          asserting a battery low condition on the status LED.
        </p>
      </div>
		 </div>
		             </div>
                              </div>
                            </div>
                          </div>
                        </div>
		);
		
	}
});

module.exports = Index;
