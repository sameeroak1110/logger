/* *****************************************************************************
Copyright (c) 2022, sameeroak1110 (sameeroak1110@gmail.com)
All rights reserved.
BSD 3-Clause License.

Package     : github.com/sameeroak1110/logger
Filename    : github.com/sameeroak1110/logger/types.go
File-type   : GoLang source code file

Compiler/Runtime: go version go1.17 linux/amd64

Version History
Version     : 1.0
Author      : sameer oak (sameeroak1110@gmail.com)
Description :
- User defined data-types used by logger package.
***************************************************************************** */
package logger

type logmessage struct {
	componentFlag int8
	component     string
	logmsg        string
}

type LogConfig struct {
	SrcBaseDir      string `json:"srcBaseDir"`       // $PWD
	FileSize        int    `json:"fileSize"`         // 20971520 (20MB)
	MaxFilesCnt     int    `json:"maxFilesCnt"`      // 10
	DefaultLogLevel string `json:"defaultLogLevel"`  // DEBUG
}

type loglevel struct {
	str   string
	color string
	wt    uint8
}
