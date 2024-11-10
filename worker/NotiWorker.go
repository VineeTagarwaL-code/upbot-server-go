package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"upbot-server-go/database"
	"upbot-server-go/libraries"
	"upbot-server-go/models"

	"github.com/resend/resend-go/v2"
)

func NotiWorker() {

	for {
		redisClient := libraries.GetInstance()

		tasks, err := redisClient.BRPop(context.Background(), 0, "noti_queue").Result()
		if err != nil {
			log.Printf("Error fetching from notification queue: %v", err)
			continue
		}

		if len(tasks) > 1 {
			task := tasks[1]
			handleTask(task)
		}
	}
}

func handleTask(task string) {

	taskID, err := strconv.Atoi(task)
	if err != nil {
		log.Printf("Invalid task ID format: %s", task)
		return
	}

	var dbTask models.Task
	if err := database.DB.First(&dbTask, taskID).Error; err != nil {
		log.Printf("Error fetching task with ID %d: %v", taskID, err)
		return
	}
	var user models.User
	if err := database.DB.First(&user, dbTask.UserID).Error; err != nil {
		log.Printf("Error fetching user with ID %d: %v", dbTask.UserID, err)
		return
	}

	if dbTask.NotifyDiscord {
		sendToDiscordWebhook(dbTask.WebHook, dbTask.URL)
	} else {
		sendFailureNotificationEmail(user.Email, dbTask.URL)
	}
}

type DiscordEmbed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Footer      EmbedFooter  `json:"footer,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedFooter struct {
	Text string `json:"text"`
	Icon string `json:"icon_url"`
}

type DiscordWebhookPayload struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

func sendToDiscordWebhook(webhook *string, url string) error {
	if webhook == nil || *webhook == "" {
		return fmt.Errorf("no Discord webhook URL provided")
	}

	// Create the embed message
	embed := DiscordEmbed{
		Title:       "🚨 Server Ping Failure Alert 🚨",
		Description: fmt.Sprintf("We have detected multiple ping failures for the server at %s.", url),
		Color:       16711680,
		Fields: []EmbedField{
			{
				Name:   "Server URL",
				Value:  fmt.Sprintf("[Visit Server](%s)", url),
				Inline: false,
			},
			{
				Name:   "Status",
				Value:  "❌ Failed to respond",
				Inline: true,
			},
			{
				Name:   "Failure Count",
				Value:  "2 consecutive failures",
				Inline: true,
			},
		},
		Footer: EmbedFooter{
			Text: "Please take immediate action to resolve this issue.",
			Icon: "https://yourdomain.com/path/to/icon.png",
		},
	}

	payload := DiscordWebhookPayload{
		Embeds: []DiscordEmbed{embed},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", *webhook+"?wait=true", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord webhook returned non-OK status: %s", resp.Status)
	}
	return nil
}

func sendFailureNotificationEmail(userEmail string, url string) error {
	subject := "⚠️ Server Ping Failure Alert for " + url
	htmlContent := fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; color: #333;">
		<table style="width: 100%%; max-width: 600px; margin: auto; background-color: #f9f9f9; padding: 20px; border-radius: 10px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);">
			<tr>
				<td style="text-align: center;">
					<h2 style="color: #d9534f;">🚨 Server Ping Failure Alert 🚨</h2>
					<p style="font-size: 18px; color: #555;">We've detected multiple failures for your monitored server.</p>
				</td>
			</tr>
			<tr>
				<td style="padding: 20px; background-color: #fff; border-radius: 8px;">
					<h3 style="color: #333; font-size: 20px;">Server Details</h3>
					<p style="font-size: 16px; margin: 5px 0;"><strong>Server URL:</strong> <a href="%s" style="color: #337ab7;"> %s </a></p>
					<p style="font-size: 16px; margin: 5px 0;"><strong>Failure Count:</strong> 2 consecutive failures</p>
				</td>
			</tr>
			<tr>
				<td style="padding: 20px;">
					<h3 style="color: #d9534f; font-size: 18px; text-align: center;">Immediate Action Recommended</h3>
					<p style="font-size: 16px; color: #666; text-align: center;">
						Our system has paused further monitoring to prevent repeated notifications. Please check your server's status and address any connectivity issues.
					</p>
					<div style="text-align: center; margin-top: 20px;">
						<a href="https://upbot.vineet.tech/dashboard" style="background-color: #5cb85c; color: white; padding: 12px 20px; border-radius: 5px; font-size: 16px; text-decoration: none;">
							Go to Dashboard
						</a>
					</div>
				</td>
			</tr>
			<tr>
				<td style="padding: 20px; text-align: center;">
					<p style="font-size: 14px; color: #999;">
						If you have any questions, please <a href="mailto:vineetagarwal.now@gmail.com" style="color: #337ab7;">contact support</a>.
					</p>
					<p style="font-size: 12px; color: #bbb;">&copy; 2024 upbot.vineet.tech  All rights reserved.</p>
				</td>
			</tr>
		</table>
	</div>
	`, url, url)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{userEmail},
		Subject: subject,
		Html:    htmlContent,
	}

	if err := libraries.SendEmail(params); err != nil {
		log.Printf("Error sending email to %s: %v", userEmail, err)
		return err
	}

	return nil
}
