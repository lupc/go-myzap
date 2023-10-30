package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/lupc/go-myzap"
	"go.uber.org/zap"
)

func main() {

	var logger = myzap.NewDefaultLogger()
	var wg = sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for {
				logger.Info(fmt.Sprintf("info测试内容，time:%v", time.Now()), zap.Any("time", time.Now()))
				logger.Debug(fmt.Sprintf("Debug测试内容，time:%v", time.Now()))
				logger.Error(fmt.Sprintf("Error测试内容，time:%v", time.Now()))

				logger.Info("test")
			}
		}()
	}
	wg.Wait()

}
