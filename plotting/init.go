package plotting

import (
	"github.com/sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

type logWrapper struct {
	*logrus.Entry
}

var (
	log = &logWrapper{
		Entry: logger.New().WithField("pkg", "evaluation/plotting"),
	}
)

func (l *logWrapper) Output(calldepth int, s string) error {
	// l.WithField("calldepth", calldepth).Debug(s)
	return nil
}

func init() {
	config.AfterInit(func() {
		log = &logWrapper{
			Entry: logger.New().WithField("pkg", "evaluation/plotting"),
		}
	})
}
