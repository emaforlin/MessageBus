package cmd

import (
	"fmt"
	"log"

	"github.com/emaforlin/inmembus/internal/client"
	"github.com/spf13/cobra"
)

var (
	publishTopic   string
	publishMessage string
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a message to a topic",
	Run: func(cmd *cobra.Command, args []string) {
		if publishTopic == "" || publishMessage == "" {
			log.Fatal("Both --topic and --message are required")
		}

		c, err := client.NewGRPCClient(addr)
		if err != nil {
			log.Fatal("Failed to connect:", err)
		}
		defer c.Close()

		if err := c.Publish(publishTopic, publishMessage); err != nil {
			log.Fatal("Publish error:", err)
		}

		fmt.Println("Message published successfully")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVar(&publishMessage, "message", "", "Message content")
	publishCmd.Flags().StringVar(&publishTopic, "topic", "", "Topic name")
}
