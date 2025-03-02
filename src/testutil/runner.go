package testutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

/*
	runner - Run a subprocess.

Inputs:

	argCmdExec: Name of the executable (E.g. jacobin)
	argOpts: Options for the executable (E.g. "-h  -verbose")
	argDeadline: Timeout in seconds
	argVerbose: If true, some activity will be logged to stderr (debug tool).

Returns:

	Result code (RcRunner* defined below)
	Combined stdout and stderr of the completed subprocess
*/
const RcRunnerSuccess = 0
const RcRunnerFailure = 1
const RcRunnerTimeout = 2

func Runner(argCmdExec string, argOpts string, argDeadlineSecs int, argVerbose bool) (int, string) {
	var cwd string
	var err error

	// Get the current working directory path.
	cwd, err = os.Getwd()
	if err != nil {
		errStr := fmt.Sprintf("Runner: os.Getwd() failed: %v", err)
		return RcRunnerFailure, errStr
	}

	// Create a background collection.
	msgPrefix := fmt.Sprintf("Runner(argCmdExec=%s, argOpts=%s, argDeadlineSecs=%d, cwd=%s)", argCmdExec, argOpts, argDeadlineSecs, cwd)
	if argVerbose {
		fmt.Fprintf(os.Stderr, "%s: verbose beginning\n", msgPrefix)
	}

	// Set up a command context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(argDeadlineSecs)*time.Second)
	defer cancel()

	// Construct a command with the given parameters by splitting the option string into a slice of strings.
	sliceOpts := strings.Split(argOpts, " ")

	// Create the command context for cmd.CombinedOutput().
	cmd := exec.CommandContext(ctx, argCmdExec, sliceOpts[:]...)

	// Run the command. Get the combined stdout and stderr text.
	outBytes, err := cmd.CombinedOutput()
	outString := string(outBytes)

	// Error occurred?
	if err != nil { // YES

		if argVerbose {
			fmt.Fprintf(os.Stderr, "%s: cmd.CombinedOutput() returned err: %v\n", msgPrefix, err)
		}

		switch ctx.Err() {

		case context.DeadlineExceeded:
			errMsg := fmt.Sprintf("%s: Deadline exceeded!", msgPrefix)
			return RcRunnerTimeout, errMsg

		case context.Canceled:
			errMsg := fmt.Sprintf("%s: Canceled!", msgPrefix)
			return RcRunnerFailure, errMsg

		default:
			outString = CleanText(outString)
			errMsg := fmt.Sprintf("%s: cmd.CombinedOutput() indicated an error: %v, outString: [%s]", 
				msgPrefix, err, outString)
			return RcRunnerFailure, errMsg
		}

	}

	// Return cleaned outString and a normal status code to caller.
	return RcRunnerSuccess, CleanText(outString)
}

// Strip out all nonprintable characters.
func CleanText(argString string) string {
	inRunes := []rune(argString)
	var outRunes []rune
	for _, rr := range inRunes {
		if !unicode.IsPrint(rr) && rr != '\n' && rr != '\r' {
			outRunes = append(outRunes, ' ')
		} else {
			outRunes = append(outRunes, rr)
		}
	}
	return string(outRunes)
}
