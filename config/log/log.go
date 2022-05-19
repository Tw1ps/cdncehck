package log

import (
	"go.uber.org/zap"
)

var (
	logger  *zap.Logger
	Slogger *zap.SugaredLogger
)

func init() {
	logger, _ = zap.NewDevelopment()
	Slogger = logger.Sugar()

	// slogger.Debugf("debug message age is %d,agender is %s", 19, "man")
	// slogger.Info("Info() uses sprint")
	// slogger.Infof("Infof() uses %s", "sprintf")
	// slogger.Infow("Infow() allows tags", "name", "Legolas", "type", 1)
}
