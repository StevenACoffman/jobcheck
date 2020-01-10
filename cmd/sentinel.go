package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/StevenACoffman/grace"
	"github.com/spf13/cobra"
)

// sentinelCmd represents the sentinel command
var sentinelCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sentinel called")
		var sleep int
		if len(args) > 0 {
			sleepVal := args[0]
			sleep, _ = strconv.Atoi(sleepVal)
		}
		filename := "/tmp/healthy"
		if sleep <= 0 {
			touchFile(filename)
		} else {
			touchFile(filename)
			wait, ctx := grace.NewWait()

			err := wait.WaitWithFunc(func() error {
				ticker := time.NewTicker(time.Duration(sleep) * time.Second)
				for {
					select {
					case <-ticker.C:
						touchFile(filename)
						// testcase what happens if an error occured
						// return fmt.Errorf("test error ticker 2s")
					case <-ctx.Done():
						log.Printf("closing ticker 2s goroutine\n")
						return nil
					}
				}
			})

			if err != nil {
				log.Printf("received error: %v", err)

			} else {
				log.Println("finished clean")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sentinelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sentinelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sentinelCmd.Flags().IntVarP(&sleep,"sleep", "s", 0, "Sleep duration in seconds")
}

func touchFile(filename string) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	// get last modified time
	file, err := os.Stat(filename)

	if os.IsNotExist(err) {
		// write timestamp to file so you can compare initial creation to last modification
		const stampF = "20060102150405"
		err = ioutil.WriteFile(filename, []byte(now.Format(stampF)), 0644)
		if err != nil {
			fmt.Println(err)
		}
		file, err = os.Stat(filename)
	}

	if err != nil {
		fmt.Println(err)
	}

	modTime := file.ModTime()
	fmt.Println("Last modified time : ", modTime)

	// change both atime and mtime to current
	err = os.Chtimes(filename, now, now)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Changed the file time : ", now.Format(time.RFC3339))
}
