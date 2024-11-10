package libraries

import (
	"log"
	"os"
	"sync"

	"github.com/resend/resend-go/v2"
)

var (
	client     *resend.Client
	clientOnce sync.Once
)

func getClient() *resend.Client {
	clientOnce.Do(func() {
		RESEND_API_KEY := os.Getenv("RESEND_API_KEY")
		log.Println("Initializing Resend client")
		client = resend.NewClient(RESEND_API_KEY)
	})
	return client
}

func SendEmail(params *resend.SendEmailRequest) error {
	client := getClient()

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}
	return nil
}
