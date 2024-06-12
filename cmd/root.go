package cmd

import (
	"os"
	"fmt"
	"github.com/ManudL2000/tgcom-cobra/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
	"log"
)

/* In these variables we store the arguments passed to flags -f, -l, -d and -a */
var FileToRead string
var LineToRead string
var DryRun bool
var ActionToDo string

/* rootCmd is the command tgcom. "Use" is the name of the command, "Short" is a brief description of the command, "Long
is a longer description of the command, Run is the action that must be executed when command tgcom is called" */
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
		UnpackFile(cmd)
	},
}

/* the one command used to run the main function (set by default by cobra-cli) */
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

/* pass in this function the flags of the command tgcom. Flags can be Persistend (so that if tgcom has a sub-command, e.g.
subtgcom, the flag defined for tgcom can be used as flags of subtgcom) or local (so flags are usable only for tgcom command) */
func init() {
	/* In the next 2 lines we modify the action to perform when -h (or --help) flag is called and we set the usage func
	that in our case will be displayed in cases where we don't define arguments of flags or so on */
	rootCmd.SetHelpFunc(customHelpFunc)
	rootCmd.SetUsageFunc(customUsageFunc)

	rootCmd.PersistentFlags().StringVarP(&FileToRead, "file", "f", "", "pass argument to the flag and will print file content")
    rootCmd.PersistentFlags().StringVarP(&LineToRead, "line", "l", "", "pass argument to line flag and will print the line specified")
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "pass argument to dry-run flag and will print the result")
	rootCmd.PersistentFlags().StringVarP(&ActionToDo, "action", "a", "toggle", "pass argument to action to comment/uncomment/toggle some lines")
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
func UnpackFile(cmd *cobra.Command){
	if strings.Contains(FileToRead, ","){
		if cmd.Flags().Changed("line"){
			fmt.Println("Warning: when passed multiple file to flag -f don't use -l flag")
		}
		fileInfo := strings.Split(FileToRead, ",")
		for i:=0; i<len(fileInfo); i++ {
			if strings.Contains(fileInfo[i], ":"){
				parts := strings.Split(fileInfo[i], ":")
				if len(parts) != 2 {
					log.Fatalf("invalid syntax. Use 'FileToRead:lines'")
				}
				utils.ChangeFile(parts[0], parts[1], ActionToDo, DryRun)
			} else {
				log.Fatalf("invalid syntax. Use 'FileToRead:lines'")
			}
		}
	} else {
		utils.ChangeFile(FileToRead, LineToRead, ActionToDo, DryRun)
	}
}

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