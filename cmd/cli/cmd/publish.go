package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/emaforlin/messagebus/internal/client"
	"github.com/spf13/cobra"
)

var (
	publishTopic   string
	publishMessage string
)

var publishCmd = &cobra.Command{
	Use: "publish <topic> <message>",
	Example: `  publish #notifications 'User created'
  echo "message" | publish #notifications
  cat file.json | publish #notifications`,
	Short: "Publish a message to a topic",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Topic is required")
		}

		publishTopic = args[0]

		// Check if we're receiving data from pipe
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Reading from pipe
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal("Failed to read from stdin:", err)
			}
			publishMessage = strings.TrimSpace(string(bytes))
		} else if len(args) > 1 {
			// Direct message from command line
			publishMessage = args[1]
		} else {
			log.Fatal("Message is required (either as argument or through pipe)")
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
}
