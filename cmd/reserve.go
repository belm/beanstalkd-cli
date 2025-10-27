package cmd

import (
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var reserveTimeout time.Duration

var reserveCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve a job from a tube",
	Long:  "Reserve a job from the watched tubes",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tubeSet := beanstalk.NewTubeSet(conn, tube)

		var id uint64
		var body []byte

		if reserveTimeout > 0 {
			id, body, err = tubeSet.Reserve(reserveTimeout)
		} else {
			id, body, err = tubeSet.Reserve(0)
		}

		if err != nil {
			return fmt.Errorf("failed to reserve job: %w", err)
		}

		color.Green("✓ Job reserved successfully")
		fmt.Printf("Job ID: %d\n", id)
		fmt.Printf("Body: %s\n", string(body))

		return nil
	},
}

func init() {
	reserveCmd.Flags().DurationVarP(&reserveTimeout, "timeout", "T", 0, "Reserve timeout (0 for blocking)")
	rootCmd.AddCommand(reserveCmd)
}
