package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var touchCmd = &cobra.Command{
	Use:   "touch [job-id]",
	Short: "Touch a reserved job",
	Long:  "Request more time to work on a reserved job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		var jobID uint64
		if _, err := fmt.Sscanf(args[0], "%d", &jobID); err != nil {
			return fmt.Errorf("invalid job ID: %w", err)
		}

		if err := conn.Touch(jobID); err != nil {
			return fmt.Errorf("failed to touch job: %w", err)
		}

		color.Green("✓ Job %d touched successfully", jobID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(touchCmd)
}
