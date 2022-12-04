package shutdown

import (
	"io"
	"os"
	"os/signal"

	"github.com/Meystergod/media-editor-project/api-service/pkg/logging"
)

func Graceful(logger *logging.Logger, signals []os.Signal, closeItems ...io.Closer) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, signals...)
	sig := <-sigChannel
	logger.Infof("caught signal %s. shutting down...", sig)

	for _, item := range closeItems {
		if err := item.Close(); err != nil {
			logger.Errorf("failed to close %v: %v", item, err)
		}
	}
}
