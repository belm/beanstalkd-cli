package cmd

import (
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/spf13/cobra"
)

var (
	host string
	port string
	tube string
)

var rootCmd = &cobra.Command{
	Use:   "beanstalkd-cli",
	Short: "A comprehensive command-line tool for Beanstalkd",
	Long: `Beanstalkd CLI is a feature-complete command-line interface for interacting with Beanstalkd,
a simple, fast work queue. It supports all Beanstalkd operations including job management,
tube operations, and server statistics.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Beanstalkd server host")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "11300", "Beanstalkd server port")
	rootCmd.PersistentFlags().StringVarP(&tube, "tube", "t", "default", "Tube name")
}

func connect() (*beanstalk.Conn, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", addr, err)
	}
	return conn, nil
}

func connectWithTimeout(timeout time.Duration) (*beanstalk.Conn, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := beanstalk.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", addr, err)
	}
	return conn, nil
}
