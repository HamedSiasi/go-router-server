/* Structs to send DL messages via the client interface from the UTM server.
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

package models

import (
    "time"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

// TODO
// type ClientMeasurementControlGeneric struct
// type ClientMeasurementControlPowerState struct
// type ClientMeasurementControlUnion union

// ClientPingReqDlMsg
type ClientPingReqDlMsg struct {
    // Empty
}

// ClientRebootReqDlMsg
type ClientRebootReqDlMsg struct {
    SdCardNotRequired bool      `bson:"sd_card_not_required" json:"sd_card_not_required"`
    DisableModemDebug bool      `bson:"disable_modem_debug" json:"disable_modem_debug"`
    DisableButton     bool      `bson:"disable_button" json:"disable_button"`
    DisableServerPing bool      `bson:"disable_server_ping" json:"disable_server_ping"`
}

// ClientDateTimeSetReqDlMsg
type ClientDateTimeSetReqDlMsg struct {
    UtmTime           time.Time `bson:"time" json:"time"`
    SetDateOnly       bool      `bson:"set_date_only" json:"set_date_only"`
}

// ClientDateTimeGetReqDlMsg
type ClientDateTimeGetReqDlMsg struct {
    // Empty    
}

type ClientModeEnum uint32

const (
    CLIENT_MODE_NULL          ClientModeEnum = iota
    CLIENT_MODE_SELF_TEST
    CLIENT_MODE_COMMISSIONING
    CLIENT_MODE_STANDARD_TRX
    CLIENT_MODE_TRAFFIC_TEST
)

// ClientModeSetReqDlMsg
type ClientModeSetReqDlMsg struct {
    Mode        ClientModeEnum  `bson:"mode" json:"mode"`
}

// ClientModeGetReqDlMsg
type ClientModeGetReqDlMsg struct {
    // Empty    
}

// ClientIntervalsGetReqDlMsg
type ClientIntervalsGetReqDlMsg struct {
    // Empty    
}

// ClientReportingIntervalSetReqDlMsg
type ClientReportingIntervalSetReqDlMsg struct {
    ReportingInterval uint32  `bson:"reporting_interval" json:"reporting_interval"`
}

// ClientHeartbeatSetReqDlMsg
type ClientHeartbeatSetReqDlMsg struct {
    HeartbeatSeconds   uint32  `bson:"heartbeat_seconds" json:"heartbeat_seconds"`
    HeartbeatSnapToRtc bool    `bson:"heartbeat_snap_to_rtc" json:"heartbeat_snap_to_rtc"`
}

// ClientMeasurementsGetReqDlMsg
type ClientMeasurementsGetReqDlMsg struct {
    // Empty
}

// TODO
// type ClientMeasurementControlSetReqDlMsg struct
// type ClientMeasurementControlSetCnfUlMsg struct
// type ClientMeasurementControlGetCnfUlMsg struct
// type ClientMeasurementControlIndUlMsg struct
// type ClientMeasurementsControlDefaultsSetReqDlMsg struct
// type ClientMeasurementsControlDefaultsSetCnfUlMsg struct

// ClientTrafficReportGetReqDlMsg
type ClientTrafficReportGetReqDlMsg struct {
    // Empty
}

// ClientTrafficTestModeParametersSetReqDlMsg
type ClientTrafficTestModeParametersSetReqDlMsg struct {
    NumUlDatagrams      uint32  `bson:"num_ul_datagrams" json:"num_ul_datagrams"`
    LenUlDatagram       uint32  `bson:"len_ul_datagram" json:"len_ul_datagram"`
    NumDlDatagrams      uint32  `bson:"num_dl_datagrams" json:"num_dl_datagrams"`
    LenDlDatagram       uint32  `bson:"len_dl_datagram" json:"len_dl_datagram"`
    TimeoutSeconds      uint32  `bson:"timeout_seconds" json:"timeout_seconds"`
    NoReportsDuringTest bool
}

// ClientTrafficTestModeParametersGetReqDlMsg
type ClientTrafficTestModeParametersGetReqDlMsg struct {
    // Empty
}

// ClientTrafficTestModeReportGetReqDlMsg
type ClientTrafficTestModeReportGetReqDlMsg struct {
    // Empty
}

// ClientActivityReportGetReqDlMsg
type ClientActivityReportGetReqDlMsg struct {
    // Empty
}

/* End Of File */
