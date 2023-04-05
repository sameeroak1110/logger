/* *****************************************************************************
Copyright (c) 2022, sameeroak1110 (sameeroak1110@gmail.com)
All rights reserved.
BSD 3-Clause License.

Package     : github.com/sameeroak1110/logger
Filename    : github.com/sameeroak1110/logger/data.go
File-type   : GoLang source code file

Compiler/Runtime: go version go1.17 linux/amd64

Version History
Version     : 1.0
Author      : sameer oak (sameeroak1110@gmail.com)
Description :
- Data used by logger package.
***************************************************************************** */
package logger

import (
	"os"
)


// log levels
const DBGRM string = "DBGRM"      // green
const DEBUG string = "DEBUG"      // normal
const INFO string = "INFO"        // normal
const WARNING string = "WARNING"  // yellow
const ERROR string = "ERROR"      // red

// buffered channel with size 10.
var chanbuffLog chan logmessage

// log-file file handler.
var pServerLogFile *os.File

var currentLogfileCnt uint8 = 1
var logfileNameList []string
var dummyLogfile string
//var loggerWG sync.WaitGroup

const log_MAX_FILES int8 = 10
const log_FILE_NAME_PREFIX string = "server.log"
const log_FILE_SIZE int64 = 20971520 // 20 MB
var current_LOG_LEVEL string = "DBGRM"

var isInit bool
var isLoggerInstanceInit bool
var srcBaseDir string

const colorNornal string = "\033[0m"
const colorErrorRed string = "\033[31m"
const colorDbgrmGreen string = "\033[32m"
const colorWarnYellow string = "\033[33m"

var loglevelMap = map[string]loglevel {
	"DBGRM": loglevel {
		str:   DBGRM,
		color: colorDbgrmGreen,
		wt:    0,
	},
	"DEBUG": loglevel {
		str:   DEBUG,
		color: colorNornal,
		wt:    1,
	},
	"INFO": loglevel {
		str:   INFO,
		color: colorNornal,
		wt:    2,
	},
	"WARNING": loglevel {
		str:   WARNING,
		color: colorWarnYellow,
		wt:    3,
	},
	"ERROR": loglevel {
		str:   ERROR,
		color: colorErrorRed,
		wt:    4,
	},
}

var doneChan chan bool
