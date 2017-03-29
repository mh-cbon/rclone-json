// Package rclone embeds rclone binary invocation and processing.
package rclone

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// CheckRes contains result data of the size command.
type CheckRes struct {
	SrcChange  string
	DstChange  string
	AllChanges string
}

//Check command
type Check struct {
	cmd     *exec.Cmd
	BinPath string
	Src     string
	Dst     string
	Res     CheckRes
}

// NewCheck returns a new rclone check Command.
func NewCheck(src, dst string) *Check {
	return &Check{Dst: dst, Src: src}
}

func (r *Check) args() []string {
	ret := []string{"check"}
	return append(ret, r.Src, r.Dst)
}

//Exec ...
func (r *Check) Exec() error {
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
		2017/03/17 17:01:38 NOTICE: Local file system at /home/mh-cbon/gow/src/github.com/mh-cbon/rclone-json/test/dest: 0 files not in Local file system at /home/mh-cbon/gow/src/github.com/mh-cbon/rclone-json/test/source
		2017/03/17 17:01:38 NOTICE: Local file system at /home/mh-cbon/gow/src/github.com/mh-cbon/rclone-json/test/source: 1 files not in Local file system at /home/mh-cbon/gow/src/github.com/mh-cbon/rclone-json/test/dest
		2017/03/17 17:01:38 ERROR : gg: File not in Local file system at /home/mh-cbon/gow/src/github.com/mh-cbon/rclone-json/test/dest
		2017/03/17 17:01:38 ERROR : file1: Sizes differ
		2017/03/17 17:01:38 ERROR : file3: Sizes differ
		2017/03/17 17:01:38 NOTICE: Local file system at /home/mh-cbon/gow/src/github.com/mh-cbon/rclone-json/test/dest: 3 differences found
		2017/03/17 17:01:38 Failed to check: 3 differences found
	*/
	src, err := filepath.Abs(r.Src)
	if err != nil {
		return err
	}
	dst, err := filepath.Abs(r.Dst)
	if err != nil {
		return err
	}

	sout := strings.Split(string(out), "\n")
	for _, line := range sout {
		t := filesNotFound.FindAllStringSubmatch(line, -1)
		if len(t) > 0 {
			if strings.Index(line, src) > -1 {
				r.Res.SrcChange = t[0][1]
			} else if strings.Index(line, dst) > -1 {
				r.Res.DstChange = t[0][1]
			}
		} else {
			t = changesFound.FindAllStringSubmatch(line, -1)
			if len(t) > 0 {
				r.Res.AllChanges = t[0][1]
			}
		}
	}
	return nil
}

//GetChangesCnt ...
func (r *Check) GetChangesCnt() string {
	return r.Res.AllChanges
}

//Kill ...
func (r *Check) Kill() error {
	if r.cmd == nil {
		return fmt.Errorf("command not %v", "started")
	}
	if r.cmd.Process == nil {
		return fmt.Errorf("command already %v", "finished")
	}
	return r.cmd.Process.Kill()
}

var filesNotFound = regexp.MustCompile(`([0-9]+)\s+files not`)
var changesFound = regexp.MustCompile(`([0-9]+)\s+differences found$`)
