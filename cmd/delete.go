package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [job-id]",
	Short: "Delete a job",
	Long:  "Delete a job by its ID",
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

		if err := conn.Delete(jobID); err != nil {
			return fmt.Errorf("failed to delete job: %w", err)
		}

		color.Green("✓ Job %d deleted successfully", jobID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
