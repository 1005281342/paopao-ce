package setup

import (
	"github.com/rocboss/paopao-ce/global"
	"github.com/rocboss/paopao-ce/internal/model"
	"github.com/rocboss/paopao-ce/pkg/logger"

	"github.com/go-redis/redis/v8"
)

func Logger() error {
	tLogger, err := logger.New(global.LoggerSetting)
	if err != nil {
		return err
	}
	global.Logger = tLogger

	return nil
}

func DBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.Redis = redis.NewClient(&redis.Options{
		Addr:     global.RedisSetting.Host,
		Password: global.RedisSetting.Password,
		DB:       global.RedisSetting.DB,
	})

	return nil
}
