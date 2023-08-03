package main

import (
	"embed"
	"one-api/common"
	"one-api/controller"
	"one-api/middleware"
	"one-api/model"
	"one-api/router"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//go:embed web/build
var buildFS embed.FS

//go:embed web/build/index.html
var indexPage []byte

func main() {
	common.SetupGinLog()
	common.SysLog("破壳AI " + common.Version + " started")
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	// Initialize SQL Database
	err := model.InitDB()
	if err != nil {
		common.FatalLog("failed to initialize database: " + err.Error())
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.FatalLog("failed to close database: " + err.Error())
		}
	}()

	// Initialize Redis
	err = common.InitRedisClient()
	if err != nil {
		common.FatalLog("failed to initialize Redis: " + err.Error())
	}

	// Initialize options
	model.InitOptionMap()
	if common.RedisEnabled {
		model.InitChannelCache()
	}
	if os.Getenv("SYNC_FREQUENCY") != "" {
		frequency, err := strconv.Atoi(os.Getenv("SYNC_FREQUENCY"))
		if err != nil {
			common.FatalLog("failed to parse SYNC_FREQUENCY: " + err.Error())
		}
		common.SyncFrequency = frequency
		go model.SyncOptions(frequency)
		if common.RedisEnabled {
			go model.SyncChannelCache(frequency)
		}
	}
	if os.Getenv("CHANNEL_UPDATE_FREQUENCY") != "" {
		frequency, err := strconv.Atoi(os.Getenv("CHANNEL_UPDATE_FREQUENCY"))
		if err != nil {
			common.FatalLog("failed to parse CHANNEL_UPDATE_FREQUENCY: " + err.Error())
		}
		go controller.AutomaticallyUpdateChannels(frequency)
	}
	if os.Getenv("CHANNEL_TEST_FREQUENCY") != "" {
		frequency, err := strconv.Atoi(os.Getenv("CHANNEL_TEST_FREQUENCY"))
		if err != nil {
			common.FatalLog("failed to parse CHANNEL_TEST_FREQUENCY: " + err.Error())
		}
		go controller.AutomaticallyTestChannels(frequency)
	}

	// Initialize HTTP server
	server := gin.Default()
	// This will cause SSE not to work!!!
	//server.Use(gzip.Gzip(gzip.DefaultCompression))
	server.Use(middleware.CORS())

	// Initialize session store

	var store sessions.Store
	store = cookie.NewStore([]byte(common.SessionSecret))
	// if os.Getenv("REDIS_CONN_STRING") == "" {
	// 	store = cookie.NewStore([]byte(common.SessionSecret))
	// } else {
	// 	store, _ = redis.NewStore(10, "tcp", common.RDB.Options().Addr, common.RDB.Options().Password, []byte(common.SessionSecret))
	// 	store.Options(sessions.Options{
	// 		//7天过期
	// 		MaxAge: int(3600 * 24 * 7),
	// 	})
	// }
	server.Use(sessions.Sessions("session", store))

	router.SetRouter(server, buildFS, indexPage)
	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}
	err = server.Run(":" + port)
	if err != nil {
		common.FatalLog("failed to start HTTP server: " + err.Error())
	}
}
