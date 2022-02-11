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


/* log levels */
const DBGRM string = "DBGRM"
const DEBUG string = "DEBUG"
const INFO string = "INFO"
const WARNING string = "WARNING"
const ERROR string = "ERROR"


var loglevel map[string]uint8

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
const default_LOG_LEVEL string = "debug"
const log_FILE_SIZE int64 = 20971520 // 20 MB
