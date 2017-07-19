package httptool

import (
	"github.com/go-irain/logger"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	logger.New()
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
	logger.SetRollingFile("log", "test.log", int32(10), int64(50), logger.MB)
	for i := 0; i < 1; i++ {
		go Post("1111111", "http://180.97.81.222:6680", "asdfa=asdfasdf&a123=asdf", false, true)
	}
	time.Sleep(time.Duration(200) * time.Second)
}
