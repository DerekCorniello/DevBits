// The logger package helps set up the api system's logging system
// this package is used in conjunction with all of the other packages
// to provide both the backend and frontend good details on processes
// and any errors.
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true}) // Human-readable
	Log.SetLevel(logrus.InfoLevel)
}
