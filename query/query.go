package main

import (
    "fmt"
    "context"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/joho/godotenv"
	"os"
	"github.com/bwmarrin/discordgo"
)


//struct that represents a discord reply
type Interaction struct {
	User1 string
	User2 string
	Message1 string
	Message2 string
}

//goroutine that listens for new messages, if message is a reply, send the data as an Interaction  
func listener(discord discordgo, result chan Interaction) {
	
}

func addDatabase() {

}

//connect to neo4j aura db
func connectNeo4j() {
	ctx := context.Background()
	godotenv.Load("../.env")
	
	dbUri := os.Getenv("NEO4J_URI")
	dbUser := os.Getenv("NEO4J_USERNAME")
	dbPassword := os.Getenv("NEO4J_PASSWORD")
    driver, err := neo4j.NewDriverWithContext(
        dbUri,
        neo4j.BasicAuth(dbUser, dbPassword, ""))

    defer driver.Close(ctx)
    err = driver.VerifyConnectivity(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println("Connection established.")
	return ctx, driver
}


//conenect to discord chat websocket API
func connectDiscord() {
	discord, err := discordgo.New("Bot " + "authentication token")
	return discord
}


func main() {
	ctx, driver = connectNeo4j()
	discord = connectDiscord()

	go listener(discord)
}