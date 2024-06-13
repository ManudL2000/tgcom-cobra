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

	Run: func(cmd *cobra.Command, args []string) {
		/* If user did not call any flag then print basic info of Usage function and exit */
		if noFlagsGiven(cmd) {
			customUsageFunc(cmd)
			os.Exit(1)
		}

		/* Otherwise user need to pass something to flag -f. If this does not happen print an error
		message and exit  */
		if !cmd.Flags().Changed("file") {
			fmt.Println("Provide a valid file with the flag -f or pass it through the pipeline")
			os.Exit(1)
		}
		/* If some arguments have been passed to -f flag then process the arguments of the flag with
		the following function */
		ReadFlags(cmd)
	},
}

/* the one command used to run the main function (set by default by cobra-cli) */
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

/*
	pass in this function the flags of the command tgcom. Flags can be Persistend (so that if tgcom has a sub-command, e.g.

subtgcom, the flag defined for tgcom can be used as flags of subtgcom) or local (so flags are usable only for tgcom command)
*/
func init() {
	/* In the next 2 lines we modify the action to perform when -h (or --help) flag is called and we set the usage func
	that in our case will be displayed in cases where we don't define arguments of flags or so on */
	rootCmd.SetHelpFunc(customHelpFunc)
	rootCmd.SetUsageFunc(customUsageFunc)

	rootCmd.PersistentFlags().StringVarP(&FileToRead, "file", "f", "", "pass argument to the flag and will print file content")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.LineNum, "line", "l", "", "pass argument to line flag and will print the line specified")
	rootCmd.PersistentFlags().BoolVarP(&inputFlag.DryRun, "dry-run", "d", false, "pass argument to dry-run flag and will print the result")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.Action, "action", "a", "toggle", "pass argument to action to comment/uncomment/toggle some lines")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.StartLabel, "start-label", "s", "", "pass argument to start-label to modify lines after start-label")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.EndLabel, "end-label", "e", "", "pass argument to end-label to modify lines after end-label")
}

/* function to see if no flag is given */
func noFlagsGiven(cmd *cobra.Command) bool {
	hasFlags := false
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			hasFlags = true
		}
	})
	return !hasFlags
}

/* analyze the argument of -f. If more files are given (e.g -f file1:line1,file2:line2,file3:line3) then split each
content and pass each pair of file and corresponding line to the ChangeFile() function. Otherwise pass directly content
of -f and -l flags. */
/* TODO: make this function more pretty */

func ReadFlags(cmd *cobra.Command) {
	if strings.Contains(inputFlag.Filename, ",") {
		if cmd.Flags().Changed("line") {
			fmt.Println("Warning: when passed multiple file to flag -f don't use -l flag")
		}
		// TODO: add the possibility of adding labels (same start and same end label) for all the files
		// two possibility: if both start and end label are given then there should be no lines, otherwise
		// you should procede as before
		if cmd.Flags().Changed("start-label") && cmd.Flags().Changed("end-label") {
			fileInfo := strings.Split(FileToRead, ",")
			for i := 0; i < len(fileInfo); i++ {
				inputFlag.Filename = fileInfo[i]
				modfile.ChangeFile(inputFlag)
			}
		} else {
			fileInfo := strings.Split(FileToRead, ",")
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
			inputFlag.Filename = FileToRead
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
