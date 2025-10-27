package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	releasePriority uint32
	releaseDelay    time.Duration
)

var releaseCmd = &cobra.Command{
	Use:   "release [job-id]",
	Short: "Release a reserved job",
	Long:  "Put a reserved job back into the ready queue",
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

		if err := conn.Release(jobID, releasePriority, releaseDelay); err != nil {
			return fmt.Errorf("failed to release job: %w", err)
		}

		color.Green("✓ Job %d released successfully", jobID)
		fmt.Printf("Priority: %d\n", releasePriority)
		fmt.Printf("Delay: %s\n", releaseDelay)
		return nil
	},
}

func init() {
	releaseCmd.Flags().Uint32VarP(&releasePriority, "priority", "r", 1024, "New priority")
	releaseCmd.Flags().DurationVarP(&releaseDelay, "delay", "d", 0, "Delay before job becomes ready")
	rootCmd.AddCommand(releaseCmd)
}
