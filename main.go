package main

import (
	"chat_system/chatService"
	"context"
	"fmt"
	stream "github.com/GetStream/stream-chat-go/v5"
	"log"
)

func InitializeStreamClient(ctx context.Context, apiKey, apiSecret string) (*stream.Client, error) {
	client, err := stream.NewClient(apiKey, apiSecret)
	client.UpdateAppSettings(ctx, stream.NewAppSettings().SetMultiTenant(true))
	if err != nil {
		return nil, fmt.Errorf("failed to create Stream client: %w", err)
	}
	return client, nil
}

func main() {
	//Add API Credentials from stream.io
	apiKey := ""
	apiSecret := ""
	propertyID := "99743c5e-03ec-4f13-9d34-9d13efd92016"
	teamName := "Prestige"
	ctx := context.Background()
	chanType := "messaging"

	client, err := InitializeStreamClient(ctx, apiKey, apiSecret)
	if err != nil {
		log.Fatalf("Error initializing Stream client: %v", err)
	}

	//Create Users
	users := []*stream.User{
		{ID: "user1", Name: "User One", Teams: []string{propertyID}},
		{ID: "user2", Name: "User Two", Teams: []string{propertyID}},
	}
	issuedAt := chatService.CreateUser(ctx, client, users)

	//Update Users
	chatService.UpdateUser(ctx, client, "user1")

	// Fetch Updated user details by filters
	resp := chatService.GetUserById(ctx, client, "id", "user1")
	fmt.Println(resp)

	resp = chatService.GetUserByTeams(ctx, client, "teams", propertyID)
	log.Println(resp)

	//Create Channel
	data := &stream.ChannelRequest{
		Members: []string{"user1", "user2"},
	}
	channelId := propertyID + teamName
	channel, err := chatService.CreateChannel(client, ctx, chanType, channelId, "user1", data)
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}
	log.Printf("Channel created successfully: %v\n", channel)

	//Add new users to the channel
	userId := "user3"
	chatService.AddUsersToChannel(client, ctx, channel, userId)
	log.Println("New users added to the channel successfully.")

	//Fetch Channel Details
	channelDetails := chatService.GetChannel(ctx, channelId, client)
	log.Println("Channel Details", channelDetails)

	//Send message
	text := "@Bob I told them I was pesca-pescatarian. Which is one who eats solely fish who eat other fish."
	chatService.SendMessage(ctx, channel, userId, text)

	//GetMessage
	messages := chatService.GetMessages(ctx, channel, client)
	log.Println("Message List", messages)

	//Revoke Token for user
	_, revokeErr := client.RevokeUserToken(ctx, "user1", &issuedAt)
	if revokeErr != nil {
		log.Println("Error revoking token for user", err)
	} else {
		log.Println("Revoked token for user user1")

	}
}
