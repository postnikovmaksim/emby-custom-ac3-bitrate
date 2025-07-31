package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func checkForAc3(args []string) bool {
	for i, arg := range args {
		if strings.HasPrefix(arg, "-c:a:") {
			if args[i+1] == "ac3" {
				return true
			}
		}
	}
	return false
}

func modifyAudioBitrate(args []string, bitrate string) {
	for i, arg := range args {
		if strings.HasPrefix(arg, "-ab:a:") {
			args[i+1] = bitrate
			return
		}
	}
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	} else {
		return fallback
	}
}

func main() {
	// go build -ldflags="-s -w" ffmpeg.go
	// ffmpeg.exe -y -i "its an [aac].mp4" -c:a:0 ac3 -ab:a:0 224000 -ar:a:0 48000 -ac:a:0 6 "its an [ac3].mp4"

	args := os.Args
	fmt.Println("Arguments:", args)

	if checkForAc3(args) {
		fmt.Println("AC3 audio detected. Modifying bitrate...")
		modifyAudioBitrate(args, getEnv("EMBY_CUSTOM_AC3_BITRATE", "640000"))
	} else {
		fmt.Println("No AC3 audio detected. Skipping bitrate modification.")
	}

	command := &exec.Cmd{
		Path:   getEnv("EMBY_CUSTOM_FFMPEG_PATH", "/bin/_ffmpeg"),
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	fmt.Println("Executing command:", command.Path, command.Args)

	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Command failed with exit code: %d\n", exitError.ExitCode())
			os.Exit(exitError.ExitCode())
		} else {
			fmt.Println("Command execution failed:", err)
		}
	} else {
		fmt.Println("Command executed successfully.")
	}
}
