package router

import (
	"server/internal/user"
	"server/internal/ws"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {
	r = gin.Default()

	r.POST("/api/signup", userHandler.CreateUser)
	r.POST("/api/login", userHandler.Login)
	r.GET("/api/logout", userHandler.Logout)

	r.POST("/websocket/createRoom", wsHandler.CreateRoom)
	r.GET("/websocket/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/websocket/getRooms", wsHandler.GetRooms)
	r.GET("/websocket/getClients/:roomId", wsHandler.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}