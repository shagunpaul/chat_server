package chatService

import (
	"context"
	"fmt"
	stream "github.com/GetStream/stream-chat-go/v5"
	"log"
	"time"
)

func CreateUser(ctx context.Context, client *stream.Client, users []*stream.User) time.Time {
	issuedAt := time.Now()

	resp, err := client.UpsertUsers(ctx, users...)
	if err != nil {
		fmt.Errorf("Error creating users: %v", err)
	}
	fmt.Println(resp)
	for _, user := range users {
		// Create a token with issuedAt time and no expiration
		token, err := client.CreateToken(user.ID, time.Time{}, issuedAt)
		if err != nil {
			log.Fatalf("Error generating JWT token for user %s: %v", user.ID, err)
		}
		fmt.Printf("Generated JWT token for user %s: %s\n", user.ID, token)

	}
	return issuedAt
}

func UpdateUser(ctx context.Context, client *stream.Client, userId string) {

	update := stream.PartialUserUpdate{
		ID: userId,
		Set: map[string]interface{}{
			"Teams": []string{"propertyId 2"},
		},
	}
	resp, err := client.PartialUpdateUsers(ctx, []stream.PartialUserUpdate{update})
	if err != nil {
		fmt.Errorf("Error updating user", userId, "error", err)

	}
	fmt.Println(resp)
}

func GetUserById(ctx context.Context, client *stream.Client, key, value string) *stream.QueryUsersResponse {
	filter := map[string]interface{}{
		key: map[string]string{"$eq": value},
	}

	// Query the user
	resp, err := client.QueryUsers(ctx, &stream.QueryOption{
		Filter: filter,
	})
	if err != nil {
		log.Fatalf("failed to query user: %v", err)
	}
	return resp
}

func GetUserByTeams(ctx context.Context, client *stream.Client, key, value string) *stream.QueryUsersResponse {
	filter := map[string]interface{}{
		key: map[string]string{"$contains": value},
	}
	// Query the user
	resp, err := client.QueryUsers(ctx, &stream.QueryOption{
		Filter: filter,
	})
	if err != nil {
		log.Fatalf("failed to query user: %v", err)
	}
	return resp
}

func CreateChannel(client *stream.Client, ctx context.Context, chanType, channelId, userID string, data *stream.ChannelRequest) (*stream.CreateChannelResponse, error) {
	channel, err := client.CreateChannel(ctx, chanType, channelId, userID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}
	return channel, nil
}

func AddUsersToChannel(client *stream.Client, ctx context.Context, channel *stream.CreateChannelResponse, userId string) {

	update := stream.PartialUserUpdate{
		ID: userId,
		Set: map[string]interface{}{
			"Teams": []string{"floor1"},
		},
	}

	updatedUser, err := client.PartialUpdateUsers(ctx, []stream.PartialUserUpdate{update})
	fmt.Println(updatedUser)
	if err != nil {
		fmt.Errorf("failed to update users: %w", err)
	}
	resp, err := channel.Channel.AddMembers(ctx, []string{userId})
	if err == nil {
		fmt.Printf("Users updated successfully: %v\n", resp)

	}

}

func GetChannel(ctx context.Context, channelId string, client *stream.Client) *stream.QueryChannelsResponse {
	// Query the updated channel details
	filter := map[string]interface{}{
		"id": map[string]interface{}{
			"$eq": channelId,
		},
	}
	options := &stream.QueryOption{
		Filter: filter,
	}

	queryResp, err := client.QueryChannels(ctx, options)
	if err != nil {
		log.Fatalf("Failed to query channel details: %v", err)
	}
	fmt.Println("Channel Details", queryResp)
	return queryResp
}

func SendMessage(ctx context.Context, channel *stream.CreateChannelResponse, userId, text string) {
	message := &stream.Message{
		Text: text,
	}
	_, err := channel.Channel.SendMessage(ctx, message, userId)
	if err != nil {
		fmt.Errorf("failed to send message: %w", err)
	}
}

func GetMessages(ctx context.Context, channel *stream.CreateChannelResponse, client *stream.Client) *stream.Message {
	// Query the updated channel details
	filter := map[string]interface{}{
		"id": map[string]interface{}{
			"$eq": channel.Channel.ID,
		},
	}
	options := &stream.QueryOption{
		Filter: filter,
	}

	queryResp, err := client.QueryChannels(ctx, options)
	if err != nil {
		log.Fatalf("Failed to query channel details: %v", err)
	}
	fmt.Println("Channel Details", queryResp)
	return queryResp.Channels[0].Messages[0]
}
