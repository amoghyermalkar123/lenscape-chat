package api

import (
	"fmt"
	apitypes "lenscape-chat/api-types"
	"lenscape-chat/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Api struct {
	DB             *db.DB
	ChannelManager *ChannelManager
}

var upgrader = websocket.Upgrader{
	//check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool {
		//Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		return true
	},
}

func (a *Api) LoadChannelMessages(c *gin.Context) {}
func (a *Api) LoadUserChannels(c *gin.Context)    {}

func (a *Api) ConnectToChannel(c *gin.Context) {
	channel := c.Param("channel")
	listener := a.ChannelManager.NewSession(channel)

	//upgrade get request to websocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("upgrade failed:", err)
		return
	}
	defer ws.Close()

	clientGone := c.Request.Context().Done()
	for {
		select {
		case <-clientGone:
			return
		case message := <-listener:
			//Response message to client
			err = ws.WriteJSON(message)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func (a *Api) SendMessage(c *gin.Context) {
	channel := c.Param("channel")

	chanMessage := &apitypes.ChatMessage{}
	if err := c.BindJSON(chanMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": fmt.Errorf("failed request operation: %v", err).Error()})
		return
	}
	err := a.DB.SendMessage(channel, chanMessage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": fmt.Errorf("failed request operation: %v", err).Error()})
		return
	}
	a.ChannelManager.CreateChannelMessage(channel, chanMessage)
	c.JSON(http.StatusOK, gin.H{"status": "message sent", "channel_id": channel})
}

func (a *Api) CreateChannel(c *gin.Context) {
	channelDetails := &apitypes.ChannelCreation{}
	if err := c.BindJSON(channelDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": fmt.Errorf("failed request operation: %v", err).Error()})
		return
	}
	fmt.Println(channelDetails)
	channelID, err := a.DB.CreateChannel(channelDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": fmt.Errorf("failed request operation: %v", err).Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "channel created", "channel_id": channelID})
}
