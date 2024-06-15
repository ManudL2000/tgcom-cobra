package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ManudL2000/tgcom-cobra/utils/modfile"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

/* In these variables we store the arguments passed to flags -f, -l, -d and -a */
var FileToRead string
var inputFlag modfile.Config
var sigFlag modfile.Track

/*
	rootCmd is the command tgcom. "Use" is the name of the command, "Short" is a brief description of the command, "Long

is a longer description of the command, Run is the action that must be executed when command tgcom is called"
*/
var rootCmd = &cobra.Command{
	Use:   "ciaoo",
	Short: "ciaoo is a verison of tgcom that uses cobra-cli toolkit",
	Long: `A longer description that spans multiple lines and likely contains
    examples and usage of using your application. For example:

	ciaoo is a CLI library written in Go that allows users to
	comment or uncomment pieces of code. It support many different
	languages including Go, C, Java, Python, Bash and many others....`,

	PreRun: func(cmd *cobra.Command, args []string) {
		check := false
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				check = true
			}
		})

		if !check {
			fmt.Println("No flag assigned: Welcome to tgcom")
			os.Exit(1)
		}
		fmt.Println("Some flag assined... will see")
	},

	Run: func(cmd *cobra.Command, args []string) {

		ReadFlags(cmd)
		// now sigFlag have been modified
		// During the read of the flag errors of using wrong flags have been already found

		/*
			// If user did not call any flag then print basic info of Usage function and exit
			if noFlagsGiven(cmd) {
				customUsageFunc(cmd)
				os.Exit(1)
			}

			// Otherwise user need to pass something to flag -f. If this does not happen print an error
			// message and exit
			if !cmd.Flags().Changed("file") {
				fmt.Println("Provide a valid file with the flag -f or pass it through the pipeline")
				os.Exit(1)
			}
			// If some arguments have been passed to -f flag then process the arguments of the flag with
			// the following function
			ReadFlags(cmd)
		*/
	},
}

/* the one command used to run the main function (set by default by cobra-cli) */
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpFunc(customHelpFunc)
	rootCmd.SetUsageFunc(customUsageFunc)

	rootCmd.PersistentFlags().StringVarP(&inputFlag.Filename, "file", "f", "", "pass argument to the flag and will modify the file content")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.LineNum, "line", "n", "", "pass argument to line flag and will modify the line in the specified range")
	rootCmd.PersistentFlags().BoolVarP(&inputFlag.DryRun, "dry-run", "d", false, "pass argument to dry-run flag and will print the result")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.Action, "action", "a", "toggle", "pass argument to action to comment/uncomment/toggle some lines")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.StartLabel, "start-label", "s", "", "pass argument to start-label to modify lines after start-label")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.EndLabel, "end-label", "e", "", "pass argument to end-label to modify lines up to end-label")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.Lang, "language", "l", "", "pass argument to language to specify the language of the input code")
	rootCmd.MarkFlagsRequiredTogether("start-label", "end-label")
	rootCmd.MarkFlagsMutuallyExclusive("line", "start-label")
	rootCmd.MarkFlagsMutuallyExclusive("line", "end-label")
	rootCmd.MarkFlagsOneRequired("file", "language")
	rootCmd.MarkFlagsMutuallyExclusive("file", "language")
}

/* check if no flag is given */
func noFlagsGiven(cmd *cobra.Command) bool {
	hasFlags := false
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			hasFlags = true
		}
	})
	return !hasFlags
}

/* add in Track true in correspondance of flag that have been  modified */
func RecordFlag(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		switch f.Name {
		case "file":
			if f.Changed {
				sigFlag.File = true
			} else {
				sigFlag.Stdin = true
			}
		case "line":
			if f.Changed {
				sigFlag.Line = true
			}
		case "start-label":
			if f.Changed {
				sigFlag.Start = true
			}
		case "end-label":
			if f.Changed {
				sigFlag.End = true
			}
		case "action":
			if f.Changed {
				sigFlag.Action = true
			}
		case "dryrun":
			if f.Changed {
				sigFlag.Dryrun = true
			}
		case "language":
			if f.Changed {
				sigFlag.Language = true
			}
		}
	})
}

/* analyze the argument of -f. If more files are given (e.g -f file1:line1,file2:line2,file3:line3) then split each
content and pass each pair of file and corresponding line to the ChangeFile() function. Otherwise pass directly content
of -f and -l flags. */
/* TODO: make this function more pretty: divide it into some more functions: if flag f has been assigned take and process it's
arguments (multiple or single files) (develop a function for first case). If flag f has not been given expect smth from stdin*/

func ReadFlags(cmd *cobra.Command) {
	fmt.Println(inputFlag.Filename)
	if strings.Contains(inputFlag.Filename, ",") {
		if cmd.Flags().Changed("line") {
			fmt.Println("Warning: when passed multiple file to flag -f don't use -l flag")
		}
		if cmd.Flags().Changed("start-label") && cmd.Flags().Changed("end-label") {
			fileInfo := strings.Split(inputFlag.Filename, ",")
			for i := 0; i < len(fileInfo); i++ {
				inputFlag.Filename = fileInfo[i]
				modfile.ChangeFile(inputFlag)
			}
		} else {
			fileInfo := strings.Split(inputFlag.Filename, ",")
			for i := 0; i < len(fileInfo); i++ {
				if strings.Contains(fileInfo[i], ":") {
					parts := strings.Split(fileInfo[i], ":")
					if len(parts) != 2 {
						log.Fatalf("invalid syntax. Use 'FileToRead:lines'")
					}
					inputFlag.Filename = parts[0]
					inputFlag.LineNum = parts[1]
					modfile.ChangeFile(inputFlag)
				} else {
					// HERE: insert code to say that is possible that start-end labels have been given
					log.Fatalf("invalid syntax. Use 'FileToRead:lines'")
				}
			}
		}
	} else {
		if cmd.Flags().Changed("line") || cmd.Flags().Changed("start-label") && cmd.Flags().Changed("end-label") {
			modfile.ChangeFile(inputFlag)
		} else {
			log.Fatalf("Not specified what you want to modify: add -l flag or -s and -e flags")
		}

	}
}

/* the following function decide in which mode we add/remove comments: currently (12/06/2024) only two modes exists: passing lines */

func customHelpFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Help Message for Tgcom application")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  tgcom [-f][single file or multiple files with lines] [-l][single line or range of lines] [-d][dry run]")
	fmt.Println()
	fmt.Println("Available Commands:")
	for _, c := range cmd.Commands() {
		fmt.Printf("  %s - %s\n", c.Name(), c.Short)
	}
	fmt.Println()
	fmt.Println("Flags:")
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		fmt.Printf("  --%s: %s\n", flag.Name, flag.Usage)
	})
	fmt.Println()
	fmt.Println("Use 'appname [command] --help' for more information about a command.")
}

func customUsageFunc(cmd *cobra.Command) error {
	fmt.Printf("Custom Usage Message for command: %s\n", cmd.Name())
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", cmd.UseLine())
	fmt.Println()
	fmt.Println("Flags:")
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		fmt.Printf("  --%s: %s\n", flag.Name, flag.Usage)
	})
	return nil
}
