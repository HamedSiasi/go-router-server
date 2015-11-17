/* Utility to get environment variables for the UTM server.
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

package utilities

import (
    "log"
    "os"
)

func EnvStr(name string) string {
    s := os.Getenv(name)
    if s == "" {
        log.Fatal("empty config env var:", name)
    }
    return s
}

/* End Of File */
