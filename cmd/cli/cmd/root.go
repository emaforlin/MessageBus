package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	addr    string
	rootCmd = &cobra.Command{
		Use:   "mbus",
		Short: "CLI for MessageBus over gRPC",
		Long:  `mbus is a CLI tool to interact with the MessageBus server via gRPC.`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&addr, "addr", "localhost:50051", "gRPC server address")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
