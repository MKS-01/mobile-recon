package cmd

import (
	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/adb"
	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/spf13/cobra"
)

var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "Send input events to device",
	Long:  "Simulate touch, keyboard, and button inputs",
}

var inputTextCmd = &cobra.Command{
	Use:   "text <text>",
	Short: "Type text on device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		text := args[0]
		output.Info("Typing: %s", text)

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "text", text)
		if err != nil {
			output.Error("Failed to send text: %v", err)
			return
		}

		output.Success("Text sent")
	},
}

var inputTapCmd = &cobra.Command{
	Use:   "tap <x> <y>",
	Short: "Tap screen at coordinates",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		x, y := args[0], args[1]
		output.Info("Tapping at (%s, %s)", x, y)

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "tap", x, y)
		if err != nil {
			output.Error("Failed to tap: %v", err)
			return
		}

		output.Success("Tap sent")
	},
}

var inputSwipeCmd = &cobra.Command{
	Use:   "swipe <x1> <y1> <x2> <y2> [duration_ms]",
	Short: "Swipe from one point to another",
	Args:  cobra.RangeArgs(4, 5),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		swipeArgs := []string{"shell", "input", "swipe"}
		swipeArgs = append(swipeArgs, args...)

		output.Info("Swiping from (%s, %s) to (%s, %s)", args[0], args[1], args[2], args[3])

		_, err = adb.ExecuteCommandWithDevice(serial, swipeArgs...)
		if err != nil {
			output.Error("Failed to swipe: %v", err)
			return
		}

		output.Success("Swipe sent")
	},
}

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
			output.Error("%v", err)
			return
		}

		keycode := args[0]
		output.Info("Sending keycode: %s", keycode)

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "keyevent", keycode)
		if err != nil {
			output.Error("Failed to send keycode: %v", err)
			return
		}

		output.Success("Keycode sent")
	},
}

var inputHomeCmd = &cobra.Command{
	Use:   "home",
	Short: "Press home button",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "keyevent", "3")
		if err != nil {
			output.Error("Failed to press home: %v", err)
			return
		}

		output.Success("Home button pressed")
	},
}

var inputBackCmd = &cobra.Command{
	Use:   "back",
	Short: "Press back button",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "input", "keyevent", "4")
		if err != nil {
			output.Error("Failed to press back: %v", err)
			return
		}

		output.Success("Back button pressed")
	},
}

func init() {
	rootCmd.AddCommand(inputCmd)
	inputCmd.AddCommand(inputTextCmd)
	inputCmd.AddCommand(inputTapCmd)
	inputCmd.AddCommand(inputSwipeCmd)
	inputCmd.AddCommand(inputKeyCmd)
	inputCmd.AddCommand(inputHomeCmd)
	inputCmd.AddCommand(inputBackCmd)
}
