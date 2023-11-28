package logger

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Logger struct { // TODO
}

var errorLevels = map[logrus.Level]string{
	logrus.PanicLevel: "panic",
	logrus.FatalLevel: "fatal",
	logrus.ErrorLevel: "error",
	logrus.WarnLevel:  "warning",
	logrus.InfoLevel:  "info",
	logrus.DebugLevel: "debug",
	logrus.TraceLevel: "trace",
}

func New(level string) *Logger {
	return &Logger{}
}

func (l Logger) Info(msg string) {
	pushLogger(msg, logrus.InfoLevel)
}

func (l Logger) Debug(msg string) {
	pushLogger(msg, logrus.DebugLevel)
}

func (l Logger) Trace(msg string) {
	pushLogger(msg, logrus.TraceLevel)
}

func (l Logger) Warning(msg string) {
	pushLogger(msg, logrus.WarnLevel)
}

func (l Logger) Error(msg string) {
	pushLogger(msg, logrus.ErrorLevel)
}

func (l Logger) Fatal(msg string) {
	pushLogger(msg, logrus.FatalLevel)
}

func (l Logger) Panic(msg string) {
	pushLogger(msg, logrus.PanicLevel)
}

func pushLogger(message string, currentLevel logrus.Level) {
	configLogLevel := os.Getenv("LOG_LEVEL")

	if len(configLogLevel) == 0 {
		configLogLevel = "2"
	}

	levelValue, errLevel := strconv.Atoi(configLogLevel)
	var logLevel logrus.Level

	if errLevel != nil {
		log.Println(errLevel)
	} else {
		logLevel = logrus.Level(levelValue)
	}

	if currentLevel > logLevel {
		return
	}

	flag.Parse()
	logsFilePath := getLogFilePath()
	logFile, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := &logrus.Logger{
		Out:   logFile,
		Level: logrus.TraceLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			LogFormat:       "[%time%] %msg%",
		},
	}

	levelMessage := errorLevels[currentLevel]
	logger.Printf("[%s] [%s] [%s] %s \n",
		getHostName(), "finery-pos-monitoring", levelMessage, message)
}

func getLogFilePath() string {
	containerName := os.Getenv("CONTAINER_NAME")

	if len(containerName) == 0 {
		containerName = "finery-pos-monitoring"
	}

	return fmt.Sprintf("./log/%s.log", containerName)
}

func getHostName() string {
	var hostName string
	hostNameFile, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		serverName, _ := os.Hostname()
		hostName = serverName
	} else {
		hostName = strings.ReplaceAll(string(hostNameFile), "\n", "")
	}

	return hostName
}
