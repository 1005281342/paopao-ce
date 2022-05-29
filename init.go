package main

import (
	"log"
	"sync"
	"time"

	"github.com/rocboss/paopao-ce/global"
	"github.com/rocboss/paopao-ce/pkg/setting"
	"github.com/rocboss/paopao-ce/pkg/setup"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
}

func setupSetting() error {
	var (
		cfg *setting.Setting
		err error
	)
	if cfg, err = setting.NewSetting(); err != nil {
		return err
	}

	if err = cfg.ReadSection("Server", &global.ServerSetting); err != nil {
		return err
	}
	if err = cfg.ReadSection("App", &global.AppSetting); err != nil {
		return err
	}
	if err = cfg.ReadSection("Log", &global.LoggerSetting); err != nil {
		return err
	}
	if err = cfg.ReadSection("Database", &global.DatabaseSetting); err != nil {
		return err
	}
	if err = cfg.ReadSection("Search", &global.SearchSetting); err != nil {
		return err
	}
	if err = cfg.ReadSection("Redis", &global.RedisSetting); err != nil {
		return err
	}
	if err = cfg.ReadSection("JWT", &global.JWTSetting); err != nil {
		return err
	}
	if err = cfg.ReadSection("Storage", &global.AliossSetting); err != nil {
		return err
	}

	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.Mutex = &sync.Mutex{}
	return nil
}

func setupLogger() error {
	return setup.Logger()
}

func setupDBEngine() error {
	return setup.DBEngine()
}
