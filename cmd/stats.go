package cmd

import (
	"fmt"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Server statistics",
	Long:  "Get server statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		stats, err := conn.Stats()
		if err != nil {
			return fmt.Errorf("failed to get stats: %w", err)
		}

		color.Cyan("=== Server Statistics ===")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Metric", "Value"})
		table.SetBorder(true)

		for key, value := range stats {
			table.Append([]string{key, value})
		}

		table.Render()
		return nil
	},
}

var statsJobCmd = &cobra.Command{
	Use:   "stats-job [job-id]",
	Short: "Job statistics",
	Long:  "Get statistics for a specific job",
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

		stats, err := conn.StatsJob(jobID)
		if err != nil {
			return fmt.Errorf("failed to get job stats: %w", err)
		}

		color.Cyan("=== Job %d Statistics ===", jobID)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Metric", "Value"})
		table.SetBorder(true)

		for key, value := range stats {
			table.Append([]string{key, value})
		}

		table.Render()
		return nil
	},
}

var statsTubeCmd = &cobra.Command{
	Use:   "stats-tube [tube-name]",
	Short: "Tube statistics",
	Long:  "Get statistics for a specific tube",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tubeName := tube
		if len(args) > 0 {
			tubeName = args[0]
		}

		tubeSet := beanstalk.NewTube(conn, tubeName)
		stats, err := tubeSet.Stats()
		if err != nil {
			return fmt.Errorf("failed to get tube stats: %w", err)
		}

		color.Cyan("=== Tube '%s' Statistics ===", tubeName)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Metric", "Value"})
		table.SetBorder(true)

		for key, value := range stats {
			table.Append([]string{key, value})
		}

		table.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(statsJobCmd)
	rootCmd.AddCommand(statsTubeCmd)
}
