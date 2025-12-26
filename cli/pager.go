package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// shouldUsePager determines if content should be paginated based on terminal size and content length
func shouldUsePager(content string, noPager bool) bool {
	if noPager {
		return false
	}

	// Check if output is being piped
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	// Get terminal height
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// If we can't get terminal size, don't use pager
		return false
	}

	// Count lines in content
	lines := strings.Count(content, "\n") + 1

	// Use pager if content is longer than terminal height minus some buffer
	return lines > height-3
}

// showInPager displays content using an external pager (less, more, etc.)
// Based on https://stackoverflow.com/a/54198703
func showInPager(content string) error {
	pagerCmd := os.Getenv("PAGER")
	var cmd *exec.Cmd

	if pagerCmd == "" {
		// Default to less with flags:
		// -R: render ANSI color codes correctly
		// -F: quit if content fits on one screen
		// -X: don't clear screen on exit
		cmd = exec.Command("less", "-RFX")
	} else {
		// If user has custom PAGER, use it as-is
		cmd = exec.Command("sh", "-c", pagerCmd)
	}

	// Get stdin pipe for writing content to pager
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdin pipe: %w", err)
	}

	// Connect pager's stdout and stderr to terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the pager
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting pager: %w", err)
	}

	// Write content to pager
	_, writeErr := io.WriteString(stdin, content)

	// Close stdin to signal end of input
	closeErr := stdin.Close()

	// Wait for pager to finish
	waitErr := cmd.Wait()

	// Check for errors in order of occurrence
	if writeErr != nil {
		// Ignore broken pipe errors - they occur when user quits pager early (e.g., pressing 'q')
		if !isBrokenPipe(writeErr) {
			return fmt.Errorf("error writing to pager: %w", writeErr)
		}
	}
	if closeErr != nil {
		// Also ignore broken pipe on close
		if !isBrokenPipe(closeErr) {
			return fmt.Errorf("error closing pager stdin: %w", closeErr)
		}
	}
	if waitErr != nil {
		return fmt.Errorf("error waiting for pager: %w", waitErr)
	}

	return nil
}

// isBrokenPipe checks if an error is a broken pipe error (EPIPE)
func isBrokenPipe(err error) bool {
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return errors.Is(pathErr.Err, syscall.EPIPE)
	}
	return errors.Is(err, syscall.EPIPE)
}
