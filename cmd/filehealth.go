package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// filehealthCmd represents the filehealth command
var filehealthCmd = &cobra.Command{
	Use:   "filehealth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("filehealth called")
		filename := "/tmp/healthy"
		var threshold float64 = 30
		healthy := livenessProbe(filename, threshold)

		if !healthy {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(filehealthCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filehealthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filehealthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func livenessProbe(filename string, threshold float64) bool {
	// if you aren't using UTC you are wrong.
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	// get last modified time
	file, err := os.Stat(filename)

	if err != nil {
		fmt.Println(err)
		//if file does not exist, or is unreadable, then it should fail liveness probe
		return false
	}

	modTime := file.ModTime()
	elapsed := now.Sub(modTime).Seconds()

	fmt.Println("Last modified time : ", modTime)
	fmt.Println("Seconds since last file modification : ", elapsed)
	return elapsed < threshold
}
