package cmd

import (
	"fmt"
	"strings"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var peekCmd = &cobra.Command{
	Use:   "peek [job-id]",
	Short: "Peek at a job",
	Long:  "Return job information without reserving it",
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

		body, err := conn.Peek(jobID)
		if err != nil {
			return fmt.Errorf("failed to peek job: %w", err)
		}

		color.Green("✓ Job %d", jobID)
		fmt.Printf("Body: %s\n", string(body))
		return nil
	},
}

var peekReadyCmd = &cobra.Command{
	Use:   "peek-ready",
	Short: "Peek at the next ready job",
	Long:  "Return the next ready job in the currently used tube",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tubeSet := beanstalk.NewTube(conn, tube)
		id, body, err := tubeSet.PeekReady()
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return fmt.Errorf("no ready jobs found in tube '%s'", tube)
			}
			return fmt.Errorf("failed to peek ready job: %w", err)
		}

		color.Green("✓ Next ready job in tube: %s", tube)
		fmt.Printf("Job ID: %d\n", id)
		fmt.Printf("Body: %s\n", string(body))
		return nil
	},
}

var peekDelayedCmd = &cobra.Command{
	Use:   "peek-delayed",
	Short: "Peek at the next delayed job",
	Long:  "Return the delayed job with the shortest delay left",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tubeSet := beanstalk.NewTube(conn, tube)
		id, body, err := tubeSet.PeekDelayed()
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return fmt.Errorf("no delayed jobs found in tube '%s'", tube)
			}
			return fmt.Errorf("failed to peek delayed job: %w", err)
		}

		color.Green("✓ Next delayed job in tube: %s", tube)
		fmt.Printf("Job ID: %d\n", id)
		fmt.Printf("Body: %s\n", string(body))
		return nil
	},
}

var peekBuriedCmd = &cobra.Command{
	Use:   "peek-buried",
	Short: "Peek at the next buried job",
	Long:  "Return the next buried job",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tubeSet := beanstalk.NewTube(conn, tube)
		id, body, err := tubeSet.PeekBuried()
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return fmt.Errorf("no buried jobs found in tube '%s'", tube)
			}
			return fmt.Errorf("failed to peek buried job: %w", err)
		}

		color.Green("✓ Next buried job in tube: %s", tube)
		fmt.Printf("Job ID: %d\n", id)
		fmt.Printf("Body: %s\n", string(body))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(peekCmd)
	rootCmd.AddCommand(peekReadyCmd)
	rootCmd.AddCommand(peekDelayedCmd)
	rootCmd.AddCommand(peekBuriedCmd)
}
