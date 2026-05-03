package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

const (
	reset  = "\033[0m"
	green  = "\033[32m"
	cyan   = "\033[36m"
	blue   = "\033[34m"
	yellow = "\033[33m"
)

var builtins = map[string]bool{
	"cd": true, "exit": true, "cls": true, "clear": true, "neofetch": true,
}

var aliases = map[string]string{
	"ls": "dir", "cat": "type", "man": "help",
}

// Parse arguments respecting quotes
func parseArgs(input string) []string {
	var args []string
	var cur strings.Builder
	inQuote := false
	var quoteChar byte

	for i := 0; i < len(input); i++ {
		c := input[i]
		if !inQuote && (c == '"' || c == '\'') {
			inQuote = true
			quoteChar = c
			continue
		}
		if inQuote && c == quoteChar {
			inQuote = false
			quoteChar = 0
			continue
		}
		if !inQuote && c == ' ' {
			if cur.Len() > 0 {
				args = append(args, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteByte(c)
	}
	if cur.Len() > 0 {
		args = append(args, cur.String())
	}
	return args
}

// Built-in: cd
func builtinCD(args []string) error {
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return os.Chdir(home)
	}
	path := strings.Join(args, " ")
	return os.Chdir(path)
}

// Built-in: clear screen
func builtinClear() {
	fmt.Print("\033[2J\033[H")
}

// Built-in: neofetch
func builtinNeofetch() {
	host, _ := os.Hostname()
	wd, _ := os.Getwd()
	usr, _ := user.Current()
	fmt.Printf("%s┌─[ %s%s@%s%s ]%s\n", green, cyan, usr.Username, host, green, reset)
	fmt.Printf("%s└─[ %s%s%s ]%s\n", green, blue, wd, green, reset)
	fmt.Printf("  %s➜  Shell: %sOnigma Shell%s (%s)\n", yellow, cyan, reset, runtime.GOOS)
	fmt.Println()
}

// Execute external command using cmd /C
func execExternal(rawInput string) error {
	fields := strings.Fields(rawInput)
	if len(fields) > 0 {
		if alias, ok := aliases[fields[0]]; ok {
			rawInput = alias + " " + strings.Join(fields[1:], " ")
		}
	}
	cmd := exec.Command("cmd", "/C", rawInput)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Dispatch command (built-in or external)
func executeCommand(line string) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	args := parseArgs(line)
	if len(args) == 0 {
		return nil
	}
	switch strings.ToLower(args[0]) {
	case "cd":
		return builtinCD(args[1:])
	case "exit":
		os.Exit(0)
		return nil // satisfy compiler (unreachable)
	case "cls", "clear":
		builtinClear()
		return nil
	case "neofetch":
		builtinNeofetch()
		return nil
	default:
		return execExternal(line)
	}
}

// Build prompt string
func getPrompt() string {
	wd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(wd, home) {
		wd = "~" + strings.TrimPrefix(wd, home)
	}
	return fmt.Sprintf("%s➜ %sONIGMA%s %s(%s)%s > ",
		green, cyan, reset, blue, wd, reset)
}

// Print startup banner
func printBanner() {
	fmt.Print(`
   ██████╗ ███╗   ██╗██╗ ██████╗ ███╗   ███╗ █████╗ 
  ██╔═══██╗████╗  ██║██║██╔════╝ ████╗ ████║██╔══██╗
  ██║   ██║██╔██╗ ██║██║██║  ███╗██╔████╔██║███████║
  ██║   ██║██║╚██╗██║██║██║   ██║██║╚██╔╝██║██╔══██║
  ╚██████╔╝██║ ╚████║██║╚██████╔╝██║ ╚═╝ ██║██║  ██║
   ╚═════╝ ╚═╝  ╚═══╝╚═╝ ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝
                                                      
`)
	fmt.Printf("%s          Onigma Shell – Windows Terminal Experience%s\n\n", cyan, reset)
}

// Main entry point
func main() {
	printBanner()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(getPrompt())
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if err := executeCommand(line); err != nil {
			fmt.Fprintf(os.Stderr, "%sError:%s %v\n", yellow, reset, err)
		}
	}
}