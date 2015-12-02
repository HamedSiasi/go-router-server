/**
 * Copyright (C) u-blox Melbourn Ltd
 * u-blox Melbourn Ltd, Melbourn, UK
 * 
 * All rights reserved.
 *
 * This source file is the sole property of u-blox Melbourn Ltd.
 * Reproduction or utilization of this source in whole or part is
 * forbidden without the written consent of u-blox Melbourn Ltd.
 */

// Present a date/time as a local one:
LocalTime =  function(time) {
    offset = new Date().getTimezoneOffset() / 60;
	
    return time - offset; 
}

module.exports = LocalTime;