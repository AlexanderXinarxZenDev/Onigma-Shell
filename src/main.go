package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"runtime"
	"strings"
	"syscall"
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
	"ls": "dir", // won't be used on Unix, but kept for cross‑platform compatibility
	"cat": "type",
	"man": "help",
}

// -------------------------------------------------------------------
// Signal handling for Ctrl+C (will not exit the shell)
// -------------------------------------------------------------------
func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		for range c {
			fmt.Println()
			fmt.Print(getPrompt())
		}
	}()
}

// -------------------------------------------------------------------
// Parse arguments respecting quotes
// -------------------------------------------------------------------
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

// -------------------------------------------------------------------
// Built-in: cd
// -------------------------------------------------------------------
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

// -------------------------------------------------------------------
// Built-in: clear screen
// -------------------------------------------------------------------
func builtinClear() {
	fmt.Print("\033[2J\033[H")
}

// -------------------------------------------------------------------
// Built-in: neofetch
// -------------------------------------------------------------------
func builtinNeofetch() {
	host, _ := os.Hostname()
	wd, _ := os.Getwd()
	usr, _ := user.Current()
	fmt.Printf("%s┌─[ %s%s@%s%s ]%s\n", green, cyan, usr.Username, host, green, reset)
	fmt.Printf("%s└─[ %s%s%s ]%s\n", green, blue, wd, green, reset)
	fmt.Printf("  %s➜  Shell: %sOnigma Shell%s (%s)\n", yellow, cyan, reset, runtime.GOOS)
	fmt.Println()
}

// -------------------------------------------------------------------
// Get current Git branch (if inside a git repo)
// -------------------------------------------------------------------
func getGitBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir, _ = os.Getwd()
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	branch := strings.TrimSpace(string(out))
	if branch == "" {
		return ""
	}
	return branch
}

// -------------------------------------------------------------------
// Get active virtual environment name (conda, venv, etc.)
// -------------------------------------------------------------------
func getVirtualEnv() string {
	if env := os.Getenv("VIRTUAL_ENV"); env != "" {
		parts := strings.Split(env, string(os.PathSeparator))
		return parts[len(parts)-1]
	}
	if env := os.Getenv("CONDA_DEFAULT_ENV"); env != "" {
		return env
	}
	return ""
}

// -------------------------------------------------------------------
// Execute external command using OS‑native shell
// -------------------------------------------------------------------
func execExternal(rawInput string) error {
	// Handle .ong script files
	if strings.HasSuffix(strings.ToLower(rawInput), ".ong") {
		return runOngScript(rawInput)
	}

	// Apply aliases (only on the first word)
	fields := strings.Fields(rawInput)
	if len(fields) > 0 {
		if alias, ok := aliases[fields[0]]; ok {
			rawInput = alias + " " + strings.Join(fields[1:], " ")
		}
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", rawInput)
	default: // Linux, macOS, BSD...
		cmd = exec.Command("sh", "-c", rawInput)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// -------------------------------------------------------------------
// Run .ong script (Onigma script – cross‑platform)
// -------------------------------------------------------------------
func runOngScript(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read .ong file: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := executeCommand(line); err != nil {
			return fmt.Errorf("error in .ong script at line %d: %v", i+1, err)
		}
	}
	return nil
}

// -------------------------------------------------------------------
// Dispatch command (built-in or external)
// -------------------------------------------------------------------
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
		return nil // unreachable
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

// -------------------------------------------------------------------
// Build prompt: ONIGMA + Git branch + Virtual env + current directory
// -------------------------------------------------------------------
func getPrompt() string {
	wd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(wd, home) {
		wd = "~" + strings.TrimPrefix(wd, home)
	}

	prompt := fmt.Sprintf("%s➜ %sONIGMA%s", green, cyan, reset)

	gitBranch := getGitBranch()
	if gitBranch != "" {
		prompt += fmt.Sprintf(" %s[%s]%s", blue, gitBranch, reset)
	}

	venv := getVirtualEnv()
	if venv != "" {
		prompt += fmt.Sprintf(" %s(%s)%s", yellow, venv, reset)
	}

	prompt += fmt.Sprintf(" %s(%s)%s > ", blue, wd, reset)
	return prompt
}

// -------------------------------------------------------------------
// Print startup banner
// -------------------------------------------------------------------
func printBanner() {
	fmt.Print(`
   ██████╗ ███╗   ██╗██╗ ██████╗ ███╗   ███╗ █████╗ 
  ██╔═══██╗████╗  ██║██║██╔════╝ ████╗ ████║██╔══██╗
  ██║   ██║██╔██╗ ██║██║██║  ███╗██╔████╔██║███████║
  ██║   ██║██║╚██╗██║██║██║   ██║██║╚██╔╝██║██╔══██║
  ╚██████╔╝██║ ╚████║██║╚██████╔╝██║ ╚═╝ ██║██║  ██║
   ╚═════╝ ╚═╝  ╚═══╝╚═╝ ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝
                                                      
`)
	fmt.Printf("%s          Onigma Shell – Cross‑Platform Terminal Experience%s\n\n", cyan, reset)
}

// -------------------------------------------------------------------
// Main entry point
// -------------------------------------------------------------------
func main() {
	printBanner()
	setupSignalHandler()

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