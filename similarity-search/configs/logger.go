package configs

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger = StructuredLogs()

func StructuredLogs() *zap.SugaredLogger {

	logfileLocation := fmt.Sprintf("%s/logs", GetEnvWithKey(KEY_LOGFILE_PATH, "."))
	var cfg zap.Config
	cfg.OutputPaths = []string{logfileLocation, "stdout"}
	cfg.Encoding = GetEnvWithKey(KEY_LOGFILE_ENCODING, "json")
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cfg.EncoderConfig = zapcore.EncoderConfig{MessageKey: GetEnvWithKey("logfile_messageKey", "message"),
		TimeKey:      GetEnvWithKey("logfile_time", "time"),
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		LevelKey:     GetEnvWithKey("logfile_level", "level"),
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		CallerKey:    GetEnvWithKey("logfile_callerKey", "callerKey"),
		EncodeCaller: zapcore.ShortCallerEncoder}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := logger.Sugar()
	defer log.Sync()
	return log
}
