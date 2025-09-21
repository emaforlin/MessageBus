package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		// c.Unsubscribe()
		c.Close()
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
