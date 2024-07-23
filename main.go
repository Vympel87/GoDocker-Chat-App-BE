package main

import (
	"log"
	"os"
	"server/db"
	"server/internal/user"
	"server/internal/ws"
	"server/router"

	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    dbConn, err := db.NewDatabase()
    if err != nil {
        log.Fatalf("could not initialize database connection: %s", err)
    }

    userRep := user.NewRepository(dbConn.GetDB())
    userSvc := user.NewService(userRep)
    userHandler := user.NewHandler(userSvc)

    hub := ws.NewHub()
    wsHandler := ws.NewHandler(hub)
    go hub.Run()

    port := os.Getenv("PORT")
    router.InitRouter(userHandler, wsHandler)
    router.Start("0.0.0.0:" + port)
}
