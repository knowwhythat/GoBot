package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type LogType zerolog.Logger

var Logger LogType

func Init() {
	logLevel := viper.GetString("log.level")
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		zerolog.SetGlobalLevel(level)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	timeFormat := "2006-01-02 15:04:05"
	zerolog.TimeFieldFormat = timeFormat

	logDir := viper.GetString("log.path")
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("Mkdir failed, err:", err)
		return
	}

	fileName := logDir + time.Now().Format("2006-01-02") + ".log"
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}
	consoleWriter.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	consoleWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	consoleWriter.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s;", i)
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	Logger = LogType(zerolog.New(multi).With().Timestamp().Logger().With().Caller().Logger())
}

func (LogType) Print(message string) {
	logger := zerolog.Logger(Logger)
	logger.Info().Msg(message)
}
func (LogType) Trace(message string) {
	logger := zerolog.Logger(Logger)
	logger.Trace().Msg(message)
}
func (LogType) Debug(message string) {
	logger := zerolog.Logger(Logger)
	logger.Debug().Msg(message)
}
func (LogType) Info(message string) {
	logger := zerolog.Logger(Logger)
	logger.Info().Msg(message)
}
func (LogType) Warning(message string) {
	logger := zerolog.Logger(Logger)
	logger.Warn().Msg(message)
}
func (LogType) Error(message string) {
	logger := zerolog.Logger(Logger)
	logger.Error().Msg(message)
}
func (LogType) Fatal(message string) {
	logger := zerolog.Logger(Logger)
	logger.Fatal().Msg(message)
}
