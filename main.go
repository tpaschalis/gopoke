package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/0xAX/notificator"
)

var notify *notificator.Notificator

func main() {
	args := os.Args
	cmd := exec.Command(args[1], args[2:]...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println(cmd.Run().Error())
	}

	notify = notificator.New(notificator.Options{
		AppName: strings.Join(args[1:], " "),
	})

	start, exitCode := time.Now(), -1
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	} else {
		exitCode = 0
	}

	switch exitCode {
	case 0:
		notify.Push("✅ Success", "Exited in "+fmtDuration(time.Since(start)), "", notificator.UR_CRITICAL)
	default:
		notify.Push("💥 Failure", "Exited in "+fmtDuration(time.Since(start)), "", notificator.UR_CRITICAL)
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
