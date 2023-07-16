/*
Copyright © 2023 Malte Schink
*/
package cmd

import (
	"fmt"
	"io"
	"notebookeR/notebook"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "notebookeR <source> [<target>]",
	Short: "Convert a jupyter notebook with R code to an R source file",
	Long: `Convert a jupyter notebook that has markdown and rlang code blocks
to a R code. You can specify a target filename to write directly in a file or omit
to write the R code to stdout`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			x := recover()
			if x != nil {
				err = fmt.Errorf("%v", x)
			}
		}()
		if len(args) < 1 {
			return fmt.Errorf("missing notebook filename")
		}
		notebookFileName := args[0]
		var rFileName string
		var rwriter io.Writer
		if len(args) > 1 {
			rFileName = args[1]
			rf, err := os.Create(rFileName)
			if err != nil {
				return err
			}
			defer rf.Close()
			rwriter = rf
		} else {
			rwriter = os.Stdout
		}

		f, err := os.Open(notebookFileName)
		if err != nil {
			return err
		}
		defer f.Close()

		n, err := notebook.Parse(f)
		if err != nil {
			return err
		}

		_, err = rwriter.Write([]byte(n.ToR()))
		if err != nil {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.notebookeR.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
