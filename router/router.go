package router

import (
	"lenscape-chat/api"
	"lenscape-chat/db"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	host := "127.0.0.1"

	db, err := db.GetDB(host)
	if err != nil {
		panic(err)
	}

	api := api.Api{
		DB:             db,
		ChannelManager: api.NewChannelManager(),
	}

	commGroup := r.Group("/communications")
	{
		// load all channels
		commGroup.GET("/channels/:user-id", api.LoadUserChannels)
		// connect to a channel
		commGroup.GET("/connect/:channel", api.ConnectToChannel)
		// load messages from channel
		commGroup.GET("/:channel/messages", api.LoadChannelMessages)
		// send message to channel
		commGroup.POST("/:channel/message", api.SendMessage)
		// create channel
		commGroup.POST("/channel", api.CreateChannel)
	}
	return r
}
