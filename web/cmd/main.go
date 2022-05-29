package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"web"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rocboss/paopao-ce/global"
	"github.com/rocboss/paopao-ce/pkg/setting"
	"github.com/rocboss/paopao-ce/pkg/setup"
)

func init() {
	var err error
	if err = setupSetting(); err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	if err = setup.Logger(); err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
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

	if err = cfg.ReadSection("Web", &global.WebSetting); err != nil {
		return err
	}
	if err = cfg.ReadSection("Log", &global.LoggerSetting); err != nil {
		return err
	}

	global.WebSetting.ReadTimeout *= time.Second
	global.WebSetting.WriteTimeout *= time.Second
	global.Mutex = &sync.Mutex{}
	return nil
}

//func main() {
//	http.Handle("/",
//		http.FileServer(&assetfs.AssetFS{
//			Asset:     web.Asset,
//			AssetDir:  web.AssetDir,
//			AssetInfo: web.AssetInfo,
//			Prefix:    "web/dist",
//			Fallback:  "index.html",
//		}),
//	)
//	http.ListenAndServe(":"+global.WebSetting.HttpPort, nil)
//}

func main() {
	gin.SetMode(global.WebSetting.RunMode)

	s := &http.Server{
		Addr:           global.WebSetting.HttpIp + ":" + global.WebSetting.HttpPort,
		Handler:        NewRouter(),
		ReadTimeout:    global.WebSetting.ReadTimeout,
		WriteTimeout:   global.WebSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("\x1b[36m%s\x1b[0m\n", "PaoPao web listen on http://"+global.WebSetting.HttpIp+":"+global.
		WebSetting.HttpPort)
	s.ListenAndServe()
}

const staticPath = `web/dist/`

func NewRouter() *gin.Engine {
	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 跨域配置
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsConfig))

	var fs = assetfs.AssetFS{
		Asset:     web.Asset,
		AssetDir:  web.AssetDir,
		AssetInfo: nil,
		Prefix:    staticPath,
		Fallback:  "index.html",
	}
	r.StaticFS("/favicon.ico", &fs)

	var assets = assetfs.AssetFS{
		Asset:     web.Asset,
		AssetDir:  web.AssetDir,
		AssetInfo: nil,
		Prefix:    staticPath + "assets",
		Fallback:  "index.html",
	}
	r.StaticFS("/assets", &assets)

	r.GET("/", indexHandler)
	// 关键点【解决页面刷新404的问题】
	r.NoRoute(indexHandler)
	return r
}

func indexHandler(c *gin.Context) {
	//设置响应状态
	c.Writer.WriteHeader(http.StatusOK)
	//载入首页
	indexHTML, _ := web.Asset(staticPath + "index.html")
	c.Writer.Write(indexHTML)
	//响应HTML类型
	c.Writer.Header().Add("Accept", "text/html")
	//显示刷新
	c.Writer.Flush()
}
