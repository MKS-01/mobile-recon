// Package cmd/input implements input simulation commands.
// This file provides commands to send touch, keyboard, and button events
// to Android devices, enabling automation and testing scenarios.
//
// Commands:
//   - input text:   Type text on the device keyboard
//   - input tap:    Tap screen at specific coordinates
//   - input swipe:  Swipe from one point to another with optional duration
//   - input key:    Send hardware key events by keycode
//   - input home:   Press the home button (keycode 3)
//   - input back:   Press the back button (keycode 4)
//
// Use Cases:
//   - UI automation and testing
//   - Scripted interactions for demos
//   - Bypassing manual input for repetitive tasks
//   - Remote device control
package cmd

import (
	"github.com/mks/adb-toolkit/pkg/adb"
	"github.com/mks/adb-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

// inputCmd is the parent command for all input simulation operations.
var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "Send input events to device",
	Long:  "Simulate touch, keyboard, and button inputs",
}

// inputTextCmd simulates typing text on the device.
// Note: Spaces in text should be handled by the shell or escaped.
// Example: adb-toolkit input text "Hello"
var inputTextCmd = &cobra.Command{
	Use:   "text <text>",
	Short: "Type text on device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		text := args[0]
		utils.PrintInfo("Typing: %s", text)

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "text", text)
		if err != nil {
			utils.PrintError("Failed to send text: %v", err)
			return
		}

		utils.PrintSuccess("Text sent")
	},
}

// inputTapCmd simulates a tap/touch event at specific screen coordinates.
// Coordinates are in pixels relative to the screen resolution.
// Example: adb-toolkit input tap 500 1000
var inputTapCmd = &cobra.Command{
	Use:   "tap <x> <y>",
	Short: "Tap screen at coordinates",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		x, y := args[0], args[1]
		utils.PrintInfo("Tapping at (%s, %s)", x, y)

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "tap", x, y)
		if err != nil {
			utils.PrintError("Failed to tap: %v", err)
			return
		}

		utils.PrintSuccess("Tap sent")
	},
}

// inputSwipeCmd simulates a swipe gesture from one point to another.
// Optional duration parameter controls the swipe speed (in milliseconds).
// Example: adb-toolkit input swipe 500 1000 500 200 300
var inputSwipeCmd = &cobra.Command{
	Use:   "swipe <x1> <y1> <x2> <y2> [duration_ms]",
	Short: "Swipe from one point to another",
	Args:  cobra.RangeArgs(4, 5),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		swipeArgs := []string{"shell", "input", "swipe"}
		swipeArgs = append(swipeArgs, args...)

		utils.PrintInfo("Swiping from (%s, %s) to (%s, %s)", args[0], args[1], args[2], args[3])

		_, err = adb.ExecuteCommandWithDevice(serial, swipeArgs...)
		if err != nil {
			utils.PrintError("Failed to swipe: %v", err)
			return
		}

		utils.PrintSuccess("Swipe sent")
	},
}

// inputKeyCmd sends hardware key events using Android keycodes.
// Full keycode list: https://developer.android.com/reference/android/view/KeyEvent
var inputKeyCmd = &cobra.Command{
	Use:   "key <keycode>",
	Short: "Send keycode event",
	Long: `Send keycode event. Common codes:
  3  - Home
  4  - Back
  26 - Power
  27 - Camera
  66 - Enter
  67 - Backspace
  82 - Menu`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		keycode := args[0]
		utils.PrintInfo("Sending keycode: %s", keycode)

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "keyevent", keycode)
		if err != nil {
			utils.PrintError("Failed to send keycode: %v", err)
			return
		}

		utils.PrintSuccess("Keycode sent")
	},
}

// inputHomeCmd is a convenience command to press the home button (keycode 3).
// Equivalent to: adb-toolkit input key 3
var inputHomeCmd = &cobra.Command{
	Use:   "home",
	Short: "Press home button",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "keyevent", "3")
		if err != nil {
			utils.PrintError("Failed to press home: %v", err)
			return
		}

		utils.PrintSuccess("Home button pressed")
	},
}

// inputBackCmd is a convenience command to press the back button (keycode 4).
// Equivalent to: adb-toolkit input key 4
var inputBackCmd = &cobra.Command{
	Use:   "back",
	Short: "Press back button",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "keyevent", "4")
		if err != nil {
			utils.PrintError("Failed to press back: %v", err)
			return
		}

		utils.PrintSuccess("Back button pressed")
	},
}

// init registers all input simulation subcommands.
// This function is automatically called when the package is imported.
func init() {
	rootCmd.AddCommand(inputCmd)
	inputCmd.AddCommand(inputTextCmd)   // Type text
	inputCmd.AddCommand(inputTapCmd)    // Tap screen
	inputCmd.AddCommand(inputSwipeCmd)  // Swipe gesture
	inputCmd.AddCommand(inputKeyCmd)    // Generic keycode
	inputCmd.AddCommand(inputHomeCmd)   // Home button shortcut
	inputCmd.AddCommand(inputBackCmd)   // Back button shortcut
}
