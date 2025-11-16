package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// struct that represents a discord reply
type Interaction struct {
	User1    string
	User2    string
	Message1 string
	Message2 string
}

func processMessage(ch chan *Interaction, ctx context.Context, driver neo4j.DriverWithContext) {
	interaction := <-ch
	user1 := interaction.User1
	user2 := interaction.User2

	//add this to the database
	result, err := neo4j.ExecuteQuery(ctx, driver, `
    CREATE (a:User {name: $name})
    CREATE (b:User {name: $friendName})
    CREATE (a)-[:KNOWS]->(b)
    `,
		map[string]any{
			"name":       user1,
			"friendName": user2,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		panic(err)
	}

	summary := result.Summary
	fmt.Printf("Created %v nodes in %+v.\n",
		summary.Counters().NodesCreated(),
		summary.ResultAvailableAfter())

}

// connect to neo4j aura db
func connectNeo4j() (context.Context, neo4j.DriverWithContext) {
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

// connect to discord chat websocket API
func connectDiscord(interactions chan<- Interaction) *discordgo.Session {

	discord, _ := discordgo.New("Bot " + os.Getenv("DISCORD_AUTH"))

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		interactions <- Interaction{
			User1:    m.Author.ID,
			User2:    m.ChannelID,
			Message1: m.Content,
			Message2: m.Content,
		}
	})

	discord.Identify.Intents = discordgo.IntentsGuildMessages
	discord.Open()

	fmt.Println("Discord bot connected")

	return discord
}

func main() {
	interactions := make(chan *Interaction, 100)

	ctx, driver := connectNeo4j()

	go processMessage(interactions, ctx, driver)

	select {}
}
