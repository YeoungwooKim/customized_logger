package colorLog

import (
	"fmt"
	"io"
	"log"
	"log_test/colorPreset"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var levels = [...]string{
	"TRACE",
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

var logger = newLogger(TRACE, true, true)

type Logger struct {
	LogLevel int
	SaveLog  bool
	PostLog  bool
}

type Caller struct {
	file     string
	funcName string
	line     int
}

func SetLogLevel(level int) {
	logger.LogLevel = level
}

// Options to save conditions according to level, save log to local, send logfile to remote
func newLogger(level int, saveLog bool, postLog bool) *Logger {
	catchShutdown( /*postingLogData*/ func() {
		fmt.Printf("before close method call\n")
	})
	return &Logger{
		LogLevel: level,
		SaveLog:  saveLog,
		PostLog:  postLog,
	}
}

func GetLogger() *Logger {
	return logger
}

func SetLogger() {
	logger = &Logger{
		LogLevel: INFO,
		SaveLog:  false,
		PostLog:  false,
	}
}

// Catch shutdown signal
func catchShutdown(gracefulShutdownFunc ...func()) {
	// create channel and asign signal (1 receive)
	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGTERM, // 15
		syscall.SIGHUP,  // 1
		syscall.SIGINT,  // 2
		syscall.SIGQUIT, // 3
		os.Interrupt,    // == SIGINT
	)

	// close
	go func() {
		sig := <-sigs
		file, multiWriter := initFileIo("output.log")
		defer file.Close()
		log.SetOutput(multiWriter)
		log.SetPrefix("")
		log.Println("::: Terminating... :::\ncaught signal : ", sig)

		for i := 0; i < len(gracefulShutdownFunc); i++ {
			gracefulShutdownFunc[i]()
		}

		log.Println("			Wait for 5 second to finish processing")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

func postingLogData() {
	log.Println("exit process - 2")
	// log.Println("hello i am post binary data(ex logfiles..)\n")
	// time.Sleep(time.Second * 2)
	// go cmd email
	// cmd echo "This is body msg" > body.txt | mpack -s "This is subject" -d $(pwd)/body.txt $(pwd)/output.log recipient "MAIL_ADDRESS"
	var bodyMsg, subject, recipient string

	bodyMsg = "this is dummy body\n\t" + fmt.Sprintf("%v", time.Now().UnixMilli()) + "\n\n\n"
	subject = "[" + time.Now().Format("2006-01-02 15:04:05.000") + "] this is dummy subject"
	recipient = "MAIL_ADDRESS"

	cmdline := fmt.Sprintf(`echo "%v" > body.txt | mpack -s "%v" -d $(pwd)/body.txt $(pwd)/output.log recipient %v`, bodyMsg, subject, recipient)
	cmd := exec.Command("sh", "-c", cmdline)
	log.Printf("	 command %v", cmdline)
	if err := cmd.Run(); err != nil {
		log.Panicf("cmd panic %v", err)
		return
	}
}

// returns specific location info that log declared.
func getCallerInfo() *Caller {
	pc, file, line, _ := runtime.Caller(3)
	funcName := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return &Caller{
		file:     file,
		funcName: funcName[len(funcName)-1],
		line:     line,
	}
}

// "must" close file properly
func initFileIo(filename string) (*os.File, io.Writer) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("error occured while opening file %v", err)
		return nil, nil
	}
	return file, io.MultiWriter(file, os.Stdout)
}

func (logger *Logger) print(logLevel int, message string, a ...interface{}) {
	var logType string
	logType, message = logger.highlightMode(logLevel, fmt.Sprintf(message, a...))
	callerInfo := getCallerInfo()
	caller := fmt.Sprintf("%v::%v(%v)", callerInfo.file, callerInfo.funcName, callerInfo.line)
	if len(caller) > 25 {
		caller = caller[len(caller)-25:]
	}

	if logLevel >= logger.LogLevel {
		file, multiWriter := initFileIo("output.log")
		defer file.Close()
		log.SetOutput(multiWriter)

		log.SetFlags(0)
		log.SetPrefix(fmt.Sprintf("[%v][%v][%v]| ", getDate(), logType, caller))
		log.Println(message)
	}
}

func (logger *Logger) highlightMode(logLevel int, message string) (string, string) {
	// TRACE, DEBUG, INFO, WARN, ERROR, FATAL.
	var textPreset, messagePreset string
	switch logLevel {
	case TRACE:
		// lowest log level
	case DEBUG:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgGreen)
	case INFO:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgCyan)
	case WARN:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgMagenta)
	case ERROR:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgYellow)
		messagePreset = fmt.Sprintf("%v%v%v%v",
			colorPreset.BoldOn, "", colorPreset.Yellow, colorPreset.BgLightGray)
	case FATAL:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgRed)
		messagePreset = fmt.Sprintf("%v%v%v%v",
			colorPreset.BoldOn, colorPreset.UnderLineOn, colorPreset.Red, colorPreset.BgWhite)
	}
	return fmt.Sprintf("%v%v%v", textPreset, levels[logLevel], colorPreset.Reset),
		fmt.Sprintf("%v%v%v", messagePreset, message, colorPreset.Reset)
}

func getDate() interface{} {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// TRACE, DEBUG, INFO, WARN, ERROR, FATAL.
func Trace(message string, a ...interface{}) {
	logger.print(TRACE, message, a...)
}

func Debug(message string, a ...interface{}) {
	logger.print(DEBUG, message, a...)
}

func Info(message string, a ...interface{}) {
	logger.print(INFO, message, a...)
}

func Warn(message string, a ...interface{}) {
	logger.print(WARN, message, a...)
}

func Error(message string, a ...interface{}) {
	logger.print(ERROR, message, a...)
}

func Fatal(message string, a ...interface{}) {
	logger.print(FATAL, message, a...)
}
