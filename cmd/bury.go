package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var buryPriority uint32

var buryCmd = &cobra.Command{
	Use:   "bury [job-id]",
	Short: "Bury a reserved job",
	Long:  "Put a reserved job into the buried state",
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

		if err := conn.Bury(jobID, buryPriority); err != nil {
			return fmt.Errorf("failed to bury job: %w", err)
		}

		color.Green("✓ Job %d buried successfully", jobID)
		fmt.Printf("Priority: %d\n", buryPriority)
		return nil
	},
}

func init() {
	buryCmd.Flags().Uint32VarP(&buryPriority, "priority", "r", 1024, "New priority")
	rootCmd.AddCommand(buryCmd)
}
