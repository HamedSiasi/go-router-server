/* Helper functions for the UTM server.
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

var WakeUpCodeLookUp = map[WakeUpEnum]string {
    0:   "OK",
    1:   "Watchdog",
    2:   "Network problem",
    3:   "SD Card problem",
    4:   "Supply problem",
    5:   "Protocol problem",
    6:   "Module not responding",
    7:   "Hardware problem",
    8:   "Memory allocation problem",
    9:   "Generic failure",
    10:  "Commanded reboot",
}

var ModeLookUp = map[ModeEnum]string {
    0: "Null",
    1: "Self Test",
    2: "Commissioning",
    3: "Standard TRX",
    4: "Traffic Test",
}

var DiskSpaceLeftLookUp = map[DiskSpaceLeftEnum]string {
    0:  " < 1GB",
    1:  " > 1GB",
    2:  " > 2GB",
    3:  " > 4GB",
    4:  " max",
}

var EnergyLeftLookUp = map[EnergyLeftEnum]string {
    0:  " > 5%",
    1:  " < 10%",
    2:  " > 20%",
    3:  " > 30%",
    4:  " > 50%",
    5:  " > 70%",
    6:  " > 90%",
    7:  " max left",
}

var TimeSetByLookUp = map[TimeSetByEnum]string {
    0: "Unknown",
    1: "GNSS",
    2: "PC",
    3: "Web API",
}

var ChargerStateEnumLookUp = map[ChargerStateEnum]string {
    0: "Unknown",
    1: "Off",
    2: "On",
    3: "Fault",
}

/* End Of File */
