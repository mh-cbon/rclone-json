// Package rclone embeds rclone binary invocation and processing.
package rclone

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// SizeRes contains result data of the size command.
type SizeRes struct {
	ObjectsCount string
	Size         string
	Bytes        string
}

var oCntRegexp = regexp.MustCompile(`Total objects:\s+([0-9]+)`)
var sizeRegexp = regexp.MustCompile(`Total size:\s+([^(]+)\(([0-9]+)[^)]+\)`)

//Size command
type Size struct {
	cmd     *exec.Cmd
	BinPath string
	Dst     string
	Res     SizeRes
}

// NewSize returns a new rclone size Command.
func NewSize(dst string) *Size {
	return &Size{Dst: dst}
}

func (r *Size) args() []string {
	ret := []string{"size"}
	return append(ret, r.Dst)
}

//Exec ...
func (r *Size) Exec() error {
	if r.cmd != nil {
		return fmt.Errorf("command already %v", "started")
	}
	if r.BinPath == "" {
		r.BinPath = rclonePath
	}
	r.cmd = exec.Command(r.BinPath, r.args()...)
	out, err := r.cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	/*
	   Total objects: 8
	   Total size: 43.389 MBytes (45497088 Bytes)
	*/
	s := oCntRegexp.FindAllStringSubmatch(string(out), -1)
	r.Res.ObjectsCount = s[0][1]
	h := sizeRegexp.FindAllStringSubmatch(string(out), -1)
	r.Res.Size = strings.TrimSpace(h[0][1])
	r.Res.Bytes = h[0][2]
	return nil
}

//GetSize ...
func (r *Size) GetSize() string {
	return r.Res.Size
}

//Kill ...
func (r *Size) Kill() error {
	if r.cmd == nil {
		return fmt.Errorf("command not %v", "started")
	}
	if r.cmd.Process == nil {
		return fmt.Errorf("command already %v", "finished")
	}
	return r.cmd.Process.Kill()
}
