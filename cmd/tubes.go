package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var listTubesCmd = &cobra.Command{
	Use:   "list-tubes",
	Short: "List all tubes",
	Long:  "List all existing tubes",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tubes, err := conn.ListTubes()
		if err != nil {
			return fmt.Errorf("failed to list tubes: %w", err)
		}

		color.Cyan("=== All Tubes ===")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"#", "Tube Name"})
		table.SetBorder(true)

		for i, tube := range tubes {
			table.Append([]string{fmt.Sprintf("%d", i+1), tube})
		}

		table.Render()
		fmt.Printf("\nTotal: %d tubes\n", len(tubes))
		return nil
	},
}

var listTubeUsedCmd = &cobra.Command{
	Use:   "list-tube-used",
	Short: "Show currently used tube",
	Long:  "Show the tube currently being used",
	RunE: func(cmd *cobra.Command, args []string) error {
		color.Green("✓ Default tube configuration: %s", tube)
		fmt.Println("Note: This shows the tube specified with -t flag")
		return nil
	},
}

var listTubesWatchedCmd = &cobra.Command{
	Use:   "list-tubes-watched",
	Short: "List watched tubes",
	Long:  "List all tubes currently being watched",
	RunE: func(cmd *cobra.Command, args []string) error {
		color.Cyan("=== Currently Configured Tube ===")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"#", "Tube Name"})
		table.SetBorder(true)
		table.Append([]string{"1", tube})
		table.Render()
		fmt.Println("\nNote: Watching is managed per connection. Use -t flag to specify tube.")
		return nil
	},
}

var useTubeCmd = &cobra.Command{
	Use:   "use [tube-name]",
	Short: "Use a tube",
	Long:  "Use the specified tube for putting jobs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tubeName := args[0]
		color.Green("✓ To use tube '%s', add: -t %s to your commands", tubeName, tubeName)
		fmt.Printf("Example: ./beanstalkd-cli -t %s put \"job data\"\n", tubeName)
		return nil
	},
}

var watchTubeCmd = &cobra.Command{
	Use:   "watch [tube-name]",
	Short: "Watch a tube",
	Long:  "Add a tube to the watch list for reserving jobs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tubeName := args[0]
		color.Green("✓ To watch tube '%s', add: -t %s to your reserve commands", tubeName, tubeName)
		fmt.Printf("Example: ./beanstalkd-cli -t %s reserve\n", tubeName)
		return nil
	},
}

var ignoreTubeCmd = &cobra.Command{
	Use:   "ignore [tube-name]",
	Short: "Ignore a tube",
	Long:  "Remove a tube from the watch list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		color.Green("✓ To work with a different tube, use -t flag with the desired tube name")
		fmt.Printf("Example: ./beanstalkd-cli -t default reserve\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listTubesCmd)
	rootCmd.AddCommand(listTubeUsedCmd)
	rootCmd.AddCommand(listTubesWatchedCmd)
	rootCmd.AddCommand(useTubeCmd)
	rootCmd.AddCommand(watchTubeCmd)
	rootCmd.AddCommand(ignoreTubeCmd)
}
