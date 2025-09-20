package cmd

import (
	"log"

	"github.com/emaforlin/messagebus/internal/client"
	"github.com/spf13/cobra"
)

var subscribeTopic string

var subscribeCmd = &cobra.Command{
	Use:     "subscribe <topic>",
	Short:   "Subscribe to a topic and receive messages",
	Example: "subscribe #notifications",
	Run: func(cmd *cobra.Command, args []string) {

		subscribeTopic = args[0]

		if subscribeTopic == "" {
			log.Fatal("topic is required")
		}

		c, err := client.NewGRPCClient(addr)
		if err != nil {
			log.Fatal("Failed to connect:", err)
		}
		defer c.Close()

		log.Printf("Subscribed to topic [%s]\n", subscribeTopic)
		c.Subscribe(subscribeTopic)
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
