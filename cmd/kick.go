package cmd

import (
	"fmt"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var kickCmd = &cobra.Command{
	Use:   "kick [bound]",
	Short: "Kick buried/delayed jobs",
	Long:  "Kick at most [bound] buried or delayed jobs back into the ready queue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		var bound int
		if _, err := fmt.Sscanf(args[0], "%d", &bound); err != nil {
			return fmt.Errorf("invalid bound: %w", err)
		}

		tubeSet := beanstalk.NewTube(conn, tube)
		kicked, err := tubeSet.Kick(bound)
		if err != nil {
			return fmt.Errorf("failed to kick jobs: %w", err)
		}

		color.Green("✓ Kicked %d jobs successfully", kicked)
		fmt.Printf("Tube: %s\n", tube)
		return nil
	},
}

var kickJobCmd = &cobra.Command{
	Use:   "kick-job [job-id]",
	Short: "Kick a specific job by ID",
	Long:  "Kick a specific buried or delayed job by its ID",
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

		if err := conn.KickJob(jobID); err != nil {
			return fmt.Errorf("failed to kick job: %w", err)
		}

		color.Green("✓ Job %d kicked successfully", jobID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(kickCmd)
	rootCmd.AddCommand(kickJobCmd)
}
