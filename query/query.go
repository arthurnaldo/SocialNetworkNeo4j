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


func processMessage(ch chan Interaction, ctx *context.Background, driver *new4j.driver) {
	interaction := <-ch
	user1 := interaction.user1
	user2 := interaction.user2
	message1 := interaction.message1
	message2 := interaction.message2

	//add this to the database
	

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


//connect to discord chat websocket API
func connectDiscord(interactions chan<- Interaction) *discordgo.Session {
	
	discord, _ := discordgo.New("Bot " + os.Getenv("DISCORD_AUTH"))

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		interactions <- Interaction{
			UserID:    m.Author.ID,
			ChannelID: m.ChannelID,
			Content:   m.Content,
		}
	})

	discord.Identify.Intents = discordgo.IntentsGuildMessages
	discord.Open()

	fmt.Println("Discord bot connected")

	return discord
}


func main() {
	interactions := make(chan *Interaction, 100)

	discord := connectDiscord(interactions)

	go processMessage(interactions, ctx, driver)

	select {}
}