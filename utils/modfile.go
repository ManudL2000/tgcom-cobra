package utils

import (
    "bufio"
    "fmt"
    "log"
    "os"
	"strings"
	"strconv"
	"io"
	"errors"
	"path/filepath"
)

/* Take in input the name of a file in the  current folder, a string that contains info about lines to be commented/uncommented, the action to do (comment,
uncomment or toggle, if no argument is passed to the flag -a the defualt will be toggle) and dryrun. If true the modifications will be displayed on the
terminal but will not be saved on the file. Otherwise the files will be modified */
func ChangeFile(filename string, line string, action string, dryrun bool) {

    // Open the file
    file, err := os.Open(filename)
    if err != nil {
        log.Fatalf("failed to open file: %s", err)
    }
    // Ensure file is closed at the end
    defer file.Close()

	char := selectCommentChars(filename)

    // find lines
	start, end:= FindLines(line)

	switch dryrun {
	case true:
		// Create a new scanner for the file
    	scanner := bufio.NewScanner(file)
		switch action{
		case "comment":
			currentLine := 1
    		for scanner.Scan() {
				lineContent := scanner.Text()
				if start <= currentLine && currentLine <= end {
					fmt.Println(lineContent + " " + "->" + " " + Comment(lineContent, char))
				} else {
    		    fmt.Println(lineContent)
				}
				currentLine++
    		}
			fmt.Println("\n")
			// Check for scanning errors
    		if err := scanner.Err(); err != nil {
        		log.Fatalf("error reading file: %s", err)
    		}
		case "uncomment":
			currentLine := 1
    		for scanner.Scan() {
				lineContent := scanner.Text()
				if start <= currentLine && currentLine <= end {
					fmt.Println(lineContent + " " + "->" + " " + Uncomment(lineContent, char))
				} else {
    		    fmt.Println(lineContent)
				}
				currentLine++
    		}
			fmt.Println("\n")
			// Check for scanning errors
    		if err := scanner.Err(); err != nil {
        		log.Fatalf("error reading file: %s", err)
    		}
		case "toggle":
			currentLine := 1
    		for scanner.Scan() {
				lineContent := scanner.Text()
				if start <= currentLine && currentLine <= end {
					fmt.Println(lineContent + " " + "->" + " " + ToggleComments(lineContent, char))
				} else {
    		    fmt.Println(lineContent)
				}
				currentLine++
    		}
			fmt.Println("\n")
			// Check for scanning errors
    		if err := scanner.Err(); err != nil {
        		log.Fatalf("error reading file: %s", err)
    		}
		}
	case false:
		// Create a backup of the original file
		backupFilename := filename + ".bak"
		createBackup(filename, backupFilename)

		// Create a temporary file
		tmpFilename := filename + ".tmp"
		tmpFile, err := os.Create(tmpFilename)
		if err != nil {
			restoreBackup(filename, backupFilename)
			log.Fatalf("Errore: %v", err)
		}
		defer tmpFile.Close()

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			restoreBackup(filename, backupFilename)
			tmpFile.Close()
			os.Remove(tmpFilename)
			log.Fatalf("Errore: %v", err)
		}

		err = writeChanges(file, tmpFile, start, end, action, char)

		if err != nil {
			restoreBackup(filename, backupFilename)
			tmpFile.Close()
			os.Remove(tmpFilename)
			log.Fatalf("Errore: %v", err)
		}

		if err := file.Close(); err != nil {
			restoreBackup(filename, backupFilename)
			tmpFile.Close()
			os.Remove(tmpFilename)
			log.Fatalf("Errore: %v", err)
		}

		// Close the temporary file before renaming
		if err := tmpFile.Close(); err != nil {
			os.Remove(tmpFilename)
			log.Fatalf("Errore: %v", err)
		}

		// Rename temporary file to original file
		if err := os.Rename(tmpFilename, filename); err != nil {
			restoreBackup(filename, backupFilename)
			log.Fatalf("Errore: %v", err)
		}

		// Remove backup file after successful processing
		os.Remove(backupFilename)
	}
}

func createBackup(filename, backupFilename string) {
	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Errore: %v", err)
	}
	defer inputFile.Close()

	backupFile, err := os.Create(backupFilename)
	if err != nil {
		log.Fatalf("Errore: %v", err)
	}
	defer backupFile.Close()

	_, err = io.Copy(backupFile, inputFile)
	if err != nil {
		log.Fatalf("Errore: %v", err)
	}
}

func restoreBackup(filename, backupFilename string) {
	// Remove the potentially corrupted file
	os.Remove(filename)
	// Restore the backup file
	os.Rename(backupFilename, filename)
}

func FindLines(lineStr string) (startLine int, endLine int) {
	if strings.Contains(lineStr, "-") {
		parts := strings.Split(lineStr, "-")
		if len(parts) != 2 {
			log.Fatalf("invalid range format. Use 'start-end'")
		}
		startLine, err := strconv.Atoi(parts[0])
		if err != nil || startLine <= 0 {
			log.Fatalf("invalid start line number")
		}
		endLine, err = strconv.Atoi(parts[1])
		if err != nil || endLine < startLine {
			log.Fatalf("invalid end line number")
		}
		return startLine, endLine
	} else {
		startLine, err := strconv.Atoi(lineStr)
		if err != nil || startLine <= 0 {
			log.Fatalf("please provide a valid positive integer for the line number or a range")
		}
		endLine = startLine
		return startLine, endLine
	}
}

func writeChanges(inputFile *os.File, outputFile *os.File, start int, end int, action string, char string) error {
	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	
	switch action{
	case "comment":
		currentLine := 1
		for scanner.Scan() {
			lineContent := scanner.Text()
			if start <= currentLine && currentLine <= end {
				lineContent = Comment(lineContent, char)
			}
			if _, err := writer.WriteString(lineContent + "\n"); err != nil {
				return err
			}
			currentLine++
		}
		if end > currentLine {
			return errors.New("line number is out of range")
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return writer.Flush()
	case "uncomment":
		currentLine := 1
		for scanner.Scan() {
			lineContent := scanner.Text()
			if start <= currentLine && currentLine <= end {
				lineContent = Uncomment(lineContent, char)
			}
			if _, err := writer.WriteString(lineContent + "\n"); err != nil {
				return err
			}
			currentLine++
		}
		if end > currentLine {
			return errors.New("line number is out of range")
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return writer.Flush()
	case "toggle":
		currentLine := 1
		for scanner.Scan() {
			lineContent := scanner.Text()
			if start <= currentLine && currentLine <= end {
				lineContent = ToggleComments(lineContent, char)
			}
			if _, err := writer.WriteString(lineContent + "\n"); err != nil {
				return err
			}
			currentLine++
		}
		if end > currentLine {
			return errors.New("line number is out of range")
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return writer.Flush()
	}
	return errors.New("Action provided is not valid")
}

func Comment(line string, char string) string {
	return char + " " + line
}

func Uncomment(line string, char string) string {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, char) {
		// Check for both `//` and `// ` prefixes.
		if strings.HasPrefix(trimmedLine, char + " ") {
			return strings.Replace(line, char + " ", "", 1)
		}
		return strings.Replace(line, char, "", 1)
	}
	return line
}

func ToggleComments(line string, char string) string {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, char) {
		return Uncomment(line, char)
	} else {
		return Comment(line, char)
	}
}

func selectCommentChars(filename string) string {
	extension := filepath.Ext(filename)
	var commentChars string
	switch extension {
	case ".go":
		commentChars = CommentChars["GoLang"]
	case ".js":
		commentChars = CommentChars["JS"]
	case ".sh", ".bash":
		commentChars = CommentChars["Bash"]
	case ".cpp", ".cc", ".h", ".c":
		commentChars = CommentChars["C++/C"]
	case ".java":
		commentChars = CommentChars["Java"]
	case ".py":
		commentChars = CommentChars["Pyhton"]
	case ".rb":
		commentChars = CommentChars["Ruby"]
	case ".pl":
		commentChars = CommentChars["Perl"]
	case ".php":
		commentChars = CommentChars["PHP"]
	case ".swift":
		commentChars = CommentChars["swift"]
	case ".kt", ".kts":
		commentChars = CommentChars["Kotlin"]
	case ".R":
		commentChars = CommentChars["R"]
	case ".hs":
		commentChars = CommentChars["Haskell"]
	case ".sql":
		commentChars = CommentChars["SQL"]
	case ".rs":
		commentChars = CommentChars["Rust"]
	case ".scala":
		commentChars = CommentChars["Scala"]
	case ".dart":
		commentChars = CommentChars["Dart"]
	case ".mm":
		commentChars = CommentChars["Objective-C"]
	case ".m":
		commentChars = CommentChars["MATLAB"]
	case ".lua":
		commentChars = CommentChars["Lua"]
	case ".erl":
		commentChars = CommentChars["Erlang"]
	case ".ex", ".exs":
		commentChars = CommentChars["Elixir"]
	case ".ts":
		commentChars = CommentChars["TS"]
	case ".vhdl", ".vhd":
		commentChars = CommentChars["VHDL"]
	case ".v", ".sv":
		commentChars = CommentChars["Verilog"]
	default:
		fmt.Printf("unsupported file extension: %s", extension)
		os.Exit(1)
	}
	return commentChars
}

var CommentChars = map[string]string{
	"GoLang": "//",
	"JS": "//",
	"Bash":   "#",
	"C++/C" : "//",
	"Java" : "//",
	"Pyhton" : "#",
	"Ruby" : "#",
	"Perl" : "#",
	"PHP" : "//",
	"Swift" : "//",
	"Kotlin" : "//",
	"R" : "#",
	"Haskell" : "--",
	"SQL" : "--",
	"Rust" : "//",
	"Scala" : "//",
	"Dart" : "//",
	"Objective-C" : "//",
	"MATLAB" : "%",
	"Lua" : "--",
	"Erlang" : "%",
	"Elixir" : "#",
	"TS" : "//",
	"VHDL" : "--",
	"Verilog" : "//",
}