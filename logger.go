/* ****************************************************************************
Copyright (c) 2022, sameeroak1110 (sameeroak1110@gmail.com)
All rights reserved.
BSD 3-Clause License.

Package     : github.com/sameeroak1110/logger
Filename    : github.com/sameeroak1110/logger/logger.go
File-type   : GoLang source code file

Compiler/Runtime: go version go1.17 linux/amd64

Version History
Version     : 1.0
Author      : sameer oak (sameeroak1110@gmail.com)
Description :
- logger library.
- Consumer of the logger library are any other code place.
- Consumer sends message to the logger function named logger.Log().
- Structure of a log message:
[<module_name>] [<dd-mm-yyyy:hhmmss-nnnnnnnnn-zzz>] [<loglevel>] [<filename>: <linenumber>] [<package_name>.<function_name>]:
<log_message>

- A dispatcher go routine fetches from the channel, extracts the log message, and dumps the same in the log file.
- Logfile(s) are located at: $PWD/logs
logfile name: server.log.<no>, where "no" stands for logfile number.
- Current logfile has extension .1.
- Max allowed size of a logfile is 20 MB (20971520 Bytes) and logfiles are rolled over after 10 log files.
**************************************************************************** */
package logger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"syscall"
	"strings"
)


/* ****************************************************************************
Description :
- Extracts sourceFilePath - defaultPath from sourceFilePath.

Arguments   :
1> sourceFilePath string: Absolute path of source file where logger.Log() has been called from.
2> defaultPath string: Default path component.

Return value:
1> bool: true is successful, false otherwise.
2> string: Absolute-path less default path.

Additional note: na
**************************************************************************** */
func getFilePath(sourceFilePath string, defaultPath string) (bool, string) {
	filePath := ""
	if len(defaultPath) > len(sourceFilePath) {
		return false, filePath
	}

	fmt.Printf("dbgrm::  sourceFilePath: %s,  defaultPath: %s\n", sourceFilePath, defaultPath)

	length := len(sourceFilePath) - len(defaultPath)
	var i int
	for i = 0; i < length; i++ {
		if sourceFilePath[i] == defaultPath[0] {
			if sourceFilePath[i:i+len(defaultPath)] == defaultPath {
				break
			}
		}
	}

	filePath = sourceFilePath[i+len(defaultPath) : len(sourceFilePath)]
	fmt.Printf("dbgrm::  filePath: %s\n", filePath)
	return true, filePath
}


/* ****************************************************************************
Description :
- Constructs a type logmessage variable.
- Dumps the same in the logmsg_buffered_channel

Arguments   :
1> strcomponent string: Modulename.
2> loglevelStr string:
- There exist 4 loglevels: ERROR, WARNING, INFO, and DEBUG.
The loglevels are incremental where DEBUG being the highest one and
includes all log levels.

Return value: na

Additional note: na
**************************************************************************** */
func Log(strcomponent string, loglevelStr string, msg string, args ...interface{}) {
	defer func() {  // chanbuffLog has been closed.
		if recoverVal := recover(); recoverVal != nil {
			fmt.Println("[WARNING]::  Log(): recover value:", recoverVal)
		}
	}()

	currentLoglevel := loglevelMap[current_LOG_LEVEL]  // 0: DBGRM, 1: DEBUG, 2: INFO, 3: WARNING, 4: ERROR
	msgLoglevel, isOK := loglevelMap[loglevelStr]
	if !isOK {
		return
	}
	if (msgLoglevel.wt < currentLoglevel.wt) && (currentLoglevel.wt != 1) { // silently slips through a DBGRM message when currentLoglevel.wt is 1(DEBUG)
		return
	}

	t := time.Now()
	zonename, _ := t.In(time.Local).Zone()
	msgTimeStamp := fmt.Sprintf("%02d-%02d-%d:%02d%02d%02d-%06d-%s", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zonename)
	pc, fn, line, _ := runtime.Caller(1)

	//gwd, _ := os.Getwd()
	//fmt.Printf("dbgrm::  gwd: %s\n", gwd)
	////_, filePath := getFilePath(fn, srcBaseDir)
	////srcFile1 := strings.Split(str1, str2)
	//filePath := strings.Split(fn, srcBaseDir)
	//srcFile := srcBaseDir + filePath[len(filePath) - 1]

	tmp1 := strings.Split((runtime.FuncForPC(pc).Name()), ".")
	pkgname := tmp1[0]
	srcFile := pkgname + "/" + path.Base(fn)
	funcName := tmp1[1]

	msgPrefix := ""
	if loglevelStr == "DBGRM" {
		msgPrefix = "#### "
	}

	////logMsg := fmt.Sprintf("[%s] [%s] [%s] [%s: %d] [%s]:\n", strcomponent, msgTimeStamp, loglevelStr, filepath.Base(fn), line, runtime.FuncForPC(pc).Name())
	////logMsg := fmt.Sprintf("[%s] [%s] [%s] [%s: %d] [%s]:\n",
		////strcomponent, msgTimeStamp, loglevelStr, filePath[len(filePath) - 1], line, runtime.FuncForPC(pc).Name())
	//logMsg := fmt.Sprintf("[%s] [%s] [%s] [%s: %d] [%s]:\n", strcomponent, msgTimeStamp, loglevelStr, srcFile, line, runtime.FuncForPC(pc).Name())
	logMsg := fmt.Sprintf("[%s] [%s] [%s] [%s +%d]@[%s]:\n", strcomponent, msgTimeStamp, loglevelStr, srcFile, line, funcName)
	logMsg = fmt.Sprintf(logMsg+msg, args...)
	logMsg = msgPrefix + logMsg + "\n"

	if !isLoggerInstanceInit {
		//logMsg = msgLoglevel.color + msgPrefix + logMsg + colorNornal + "\n"
		logMsg = msgLoglevel.color + logMsg + colorNornal
		fmt.Printf(logMsg)
		return
	}

	logMessage := logmessage {
		component: strcomponent,
		logmsg: logMsg,
	}

	if doneChanFlag == false {
		chanbuffLog <- logMessage
	}
}


/* ****************************************************************************
Description :
- A go routine, invoked through Logger()
- Waits onto buffered channel name chanbuffLog infinitely.
- Extracts data from the channel, it's of type logmessage.
- Dumps log into the file pointed by pServerLogFile.

Arguments   : na for now.
1> wg *sync.WaitGroup: waitgroup handler for conveying done status to the caller.
2> doneChan chan bool: done channel to terminate logger thread.

Return Value: na

Additional note: na
**************************************************************************** */
//func LogDispatcher(ploggerWG *sync.WaitGroup) {
func LogDispatcher(ploggerWG *sync.WaitGroup, doneChan chan bool) {
	defer func() {
		fmt.Println("logger exiting.")
		ploggerWG.Done()
	}()

	/* for {
		select {
			case logMsg := <-chanbuffLog: // pushes dummy logmessage onto the channel
				dumpServerLog(logMsg.logmsg)
		}
	} */


	runFlag := true
	for runFlag {
		select {
			case logMsg, isOK := <-chanbuffLog: // pushes dummy logmessage onto the channel
				if !isOK {
					runFlag = false
					break
				}
				dumpServerLog(logMsg.logmsg)
				break

			case <-doneChan:  // chanbuffLog has been closed. pull all the logs from the channel and dump them to file-system.
				doneChanFlag = true
				runFlag = false
				dumpServerLog("[WARNING]:: logger exiting. breaking out on closed log message-queue.\nstarting to flush all the blocked logs.\n")
				time.Sleep(10 * time.Second)
				close(chanbuffLog)
				for logMsg := range chanbuffLog {
					dumpServerLog(logMsg.logmsg)
				}
				break
		}
	}

	/* for runFlag {
		select {
			case <-doneChan:  // chanbuffLog needs to be closed. pull all the logs from the channel and dump them to file-system.
				runFlag = false
				dumpServerLog("[WARNING]:: logger exiting. breaking out on closed log message-queue.\nstarting to flush all the blocked logs.\n")
				close(chanbuffLog)
				for logMsg := range chanbuffLog {
					dumpServerLog(logMsg.logmsg)
				}
				break
			default:
				break
		}
		select {
			case logMsg, isOK := <-chanbuffLog: // pushes dummy logmessage onto the channel
				if !isOK {
					runFlag = false
					break
				}

				dumpServerLog(logMsg.logmsg)
				break
			default:
				break
		}
	} */
}


/* ****************************************************************************
Description :
- Dumps logMsg into target logfile pointed to by plogfile file handler.
- Dumps logMsg into the database table.

Arguments   :
1> logMsg string: log message to be dumped in the logfile.

Return Value: na

Additional note:
TODO: Dump log message into nosql db.
**************************************************************************** */
func dumpServerLog(logMsg string) {
	if pServerLogFile == nil {
		fmt.Printf("error-5\n")  // nil file handler
		os.Exit(1)
	}

	if logMsg == "" {
		return
	}

	pServerLogFile.WriteString(logMsg)
	//fmt.Printf(logMsg) // TODO-REM: remove this fmp.Printf() call later

	fi, err := pServerLogFile.Stat()
	if err != nil {
		fmt.Printf("error-6: %s\n", err.Error())  // Couldn't obtain stat
		return
	}

	fileSize := fi.Size()
	if fileSize >= log_FILE_SIZE {
		pServerLogFile.Close()
		pServerLogFile = nil
		err = os.Rename(logfileNameList[0], dummyLogfile)
		if err != nil {
			fmt.Printf("error-7: %s\n", err.Error())  // mv %s to %s, error: %s\n", logfileNameList[0], dummyLogfile, err.Error())
			pServerLogFile, err = os.OpenFile(logfileNameList[0], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
			return
		}

		pServerLogFile, err = os.OpenFile(logfileNameList[0], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("error-8: %s\n", err.Error())  // recreating logfile: %s,  error: %s\n", logfileNameList[0], err.Error())
			return
		}

		if currentLogfileCnt < 10 {
			currentLogfileCnt = currentLogfileCnt + 1
		}

		go handleLogRotate()
	}
}


/* ****************************************************************************
Description :
- Rotates logs to subsequent log file (n % 10). Each log file is 20MB (20971520 Bytes) size.
- Rolls over and starts from 1st log file if 10th log file is rotated.

Arguments   : na

Return Value: na

Additional note: na
**************************************************************************** */
func handleLogRotate() {
	for i := currentLogfileCnt; i > 2; i-- {
		err := os.Rename(logfileNameList[i-2], logfileNameList[i-1])
		if err != nil {
			// mv %s to %s. error: %s\n", logfileNameList[i-2], logfileNameList[i-1], err.Error())
			fmt.Printf("error-10: %s\n", err.Error())
			return
		}
	}

	err := os.Rename(dummyLogfile, logfileNameList[1])
	if err != nil {
		// while mv %s to %s. error: %s\n", dummyLogfile, logfileNameList[1], err.Error())
		fmt.Printf("error-11: %s\n", err.Error())
		return
	}
}


/* *****************************************************************************
Description :
- Initializes logger package data.
- Creates a directory $PWD/logs if doesn't exist and creates first logfile
underneath.

Arguments   :
1> isLoggerInit bool: true if logger data to be initialized. false in case logs are sent to stdout and not to any log file.
2> tmpSrcBaseDir: absolute path of base directory of source code tree.
3> logBaseDir: absolute path of base directory where logs will be stored.
4> logLevel: either of DEBUG, INFO, WARNING, ERROR.

Return value:
1> bool: True if successful, false otherwise.

Additional note: na
***************************************************************************** */
func Init(isLoggerInit bool, tmpSrcBaseDir string, logBaseDir string, logLevel string) bool {
	if isInit {
		return true
	}

	var err error

	if tmpSrcBaseDir = strings.TrimSpace(tmpSrcBaseDir); tmpSrcBaseDir == "" {
		fmt.Printf("Error-1: %s\nSource code BaseDir: %s\n", err.Error(), tmpSrcBaseDir)  // Error: abs path: %s\n", err.Error())
		return false
	}
	tmpSrcBaseDir = strings.TrimLeft(tmpSrcBaseDir, "/")
	tmpSrcBaseDir = strings.TrimRight(tmpSrcBaseDir, "/")
	srcBaseDir = "/" + tmpSrcBaseDir

	if logBaseDir = strings.TrimSpace(logBaseDir); logBaseDir == "" {
		if logBaseDir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			fmt.Printf("Error-1: %s\nlogBaseDir: %s\n", err.Error(), logBaseDir)  // Error: abs path: %s\n", err.Error())
			return false
		}
	}

	logLevel = strings.ToUpper(strings.TrimSpace(logLevel))
	if (logLevel != DBGRM) && (logLevel != DEBUG) && (logLevel != INFO ) && (logLevel != WARNING) && (logLevel != ERROR) {  // covers logLevel == ""
		fmt.Printf("Error-2: Incorrect log-level. Possible values are: DEBUG, INFO, WARNING, ERROR\n")
		return false
	}

	//logdir := filepath.Join(logBaseDir, filepath.Join("logs", filepath.Join("server")))
	logdir := filepath.Join(logBaseDir, "logs")
	if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
		fmt.Printf("error-3: %s\n", err.Error())  // Error: while creating logenv: %s\n", err.Error())
		return false
	}

	logfileNameList = make([]string, log_MAX_FILES)

	chanbuffLog = make(chan logmessage, 10)

	logFile := filepath.Join(logdir, log_FILE_NAME_PREFIX) + ".1"
	tmplogFile := filepath.Join(logdir, log_FILE_NAME_PREFIX)
	dummyLogfile = logFile + ".dummy"

	pServerLogFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("error-4: %s\n", err.Error())  //Error: while creating logfile: %s, error: %s\n", logFile, err.Error())
		return false
	}

	for i := int8(0); i < log_MAX_FILES; i++ {
		logfileNameList[i] = fmt.Sprintf("%s.%d", tmplogFile, i+1)
	}

	// if isLoggerInit == true: closes stderr so that error and panic logs can be captured in the logfile itself.
	if isLoggerInit {
		isLoggerInstanceInit = true
		if errDup2 := syscall.Dup2(int(pServerLogFile.Fd()), syscall.Stderr); errDup2 != nil {
			fmt.Printf("Error: Failed to reuse STDERR.\n")
		} else {
			fmt.Printf("Debug: Reused STDERR.\n")
		}

		if errDup2 := syscall.Dup2(int(pServerLogFile.Fd()), syscall.Stdout); errDup2 != nil {
			fmt.Printf("Error: Failed to reuse STDOUT.\n")
		} else {
			fmt.Printf("Debug: Reused STDOUT.\n")
		}
	}

	isInit = true
	return true
}


/*
func DeInit() bool {
	// this's a very rudimentary approach for letting logger-dispatcher to kick start.
	// but this should work since there's only main go-routing running in the same process context.
	// as well, we can safely ignore the current load on the the server hardware as 2 seconds time is sufficiently enough.
	time.Sleep(3 * time.Second)
	loggerWG.Done()
	return true
}*/
