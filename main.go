package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create the layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.Box.SetTitle("Tmux Launcher").SetBorder(true).SetBackgroundColor(tcell.ColorDefault)

	// Create the selection window
	list := tview.NewList().ShowSecondaryText(false)
	list.Box.SetTitle("Sessions").SetBorder(true).SetBackgroundColor(tcell.ColorDefault)

	// Create the preview window
	preview := tview.NewTextView().SetDynamicColors(true).SetWrap(false).SetScrollable(true)
	preview.Box.SetBorder(true).SetTitle("Preview").SetBackgroundColor(tcell.ColorDefault)

	// Function to update the preview
	updatePreview := func(sessionName string) {
		preview.Clear()
		cmd := exec.Command("tmux", "capture-pane", "-ep", "-t", sessionName)

		output, err := cmd.Output()
		if err != nil {
			preview.SetText(fmt.Sprintf("Error retrieving pane content: %v", err))
			return
		}

		// Manually parse ANSI escape sequences and apply colors
		coloredOutput := parseANSISequences(string(output))
		preview.SetBackgroundColor(tcell.ColorDefault)

		preview.SetText(coloredOutput)
	}

	// Function to update the list of sessions
	updateSessionList := func() {
		list.Clear()
		sessions, err := getTmuxSessions()
		if err != nil {
			list.AddItem("Error retrieving sessions", err.Error(), 0, nil)
			return
		}
		for i, session := range sessions {
			r := rune('0' + i)

			// Sets the initial preview to the first session option
			if r == '0' {
				updatePreview(session)
			}

			list.AddItem(session, "", r, nil)
		}
	}

	// Set up the layout
	flex.AddItem(preview, 0, 4, false).AddItem(list, 0, 1, true)

	// Update the list of sessions
	updateSessionList()

	// Handle user input
	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// Update the preview
		updatePreview(mainText)
	})

	// Handle user input
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// Print the selected session name
		fmt.Printf("Selected session: %s\n", mainText)

		openTerminalAndAttach(mainText)

		// Close the original terminal (Linux and macOS)
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			fmt.Println("closing...")
			cmd := exec.Command("exit")
			cmd.Run()
		}

		// Exit the application
		app.Stop()
	})

	// Set keybinding to exit on "Enter"
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func getTmuxSessions() ([]string, error) {
	cmd := exec.Command("tmux", "ls", "-F", "#{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	sessions := strings.Split(strings.TrimSpace(string(output)), "\n")
	return sessions, nil
}

// Function to parse ANSI escape sequences and apply colors
func parseANSISequences(input string) string {
	translated := tview.TranslateANSI(input)

	return translated
}

// Function to open a new terminal window and attach to the selected tmux session
func openTerminalAndAttach(sessionName string) error {
	var cmd *exec.Cmd

	// Check the operating system to determine how to open a new terminal
	switch runtime.GOOS {
	case "darwin": // macOS
		// Open a new Terminal window and attach to the tmux session
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "tmux attach-session -t %s"`, sessionName))
	case "linux":
		// Open a new Alacritty window and attach to the tmux session
		cmd = exec.Command("alacritty", "-e", "tmux", "attach-session", "-t", sessionName)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Redirect standard input, output, and error streams to /dev/null
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Set Detach flag to true
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	// Execute the command
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
