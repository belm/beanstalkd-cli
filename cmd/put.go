package cmd

import (
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	putPriority uint32
	putDelay    time.Duration
	putTTR      time.Duration
)

var putCmd = &cobra.Command{
	Use:   "put [data]",
	Short: "Put a job into a tube",
	Long:  "Insert a job with the specified data into the tube",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tubeSet := beanstalk.NewTube(conn, tube)
		id, err := tubeSet.Put([]byte(args[0]), putPriority, putDelay, putTTR)
		if err != nil {
			return fmt.Errorf("failed to put job: %w", err)
		}

		color.Green("✓ Job inserted successfully")
		fmt.Printf("Job ID: %d\n", id)
		fmt.Printf("Tube: %s\n", tube)
		fmt.Printf("Priority: %d\n", putPriority)
		fmt.Printf("Delay: %s\n", putDelay)
		fmt.Printf("TTR: %s\n", putTTR)

		return nil
	},
}

func init() {
	putCmd.Flags().Uint32VarP(&putPriority, "priority", "r", 1024, "Job priority (0 is highest)")
	putCmd.Flags().DurationVarP(&putDelay, "delay", "d", 0, "Job delay duration")
	putCmd.Flags().DurationVarP(&putTTR, "ttr", "T", 60*time.Second, "Time to run (TTR)")
	rootCmd.AddCommand(putCmd)
}
