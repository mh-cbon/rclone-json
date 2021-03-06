package rclone

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/mh-cbon/rclone-json/parser"
)

//Sync ...
type Sync struct {
	cmd           *exec.Cmd
	decoder       *parser.Decoder
	stderr        io.Closer
	BinPath       string
	Stdout        io.Writer
	Src           string
	Dst           string
	Verbose       bool
	Stats         string
	BwLimit       string
	Checkers      string
	TransferLimit string
}

// New returns a new rclone Command.
func New(dst, src string) *Sync {
	return &Sync{
		Verbose: true,
		Dst:     dst,
		Src:     src,
	}
}

func (r *Sync) args() []string {
	ret := []string{"sync"}
	if r.Verbose {
		ret = append(ret, "-vv")
	}
	if r.Stats != "" {
		ret = append(ret, "--stats", r.Stats)
	}
	if r.BwLimit != "" {
		ret = append(ret, "--bwlimit", r.BwLimit)
	}
	if r.Checkers != "" {
		ret = append(ret, "--checkers", r.Checkers)
	}
	if r.TransferLimit != "" {
		ret = append(ret, "--transfers", r.TransferLimit)
	}
	return append(ret, r.Src, r.Dst)
}

//Start ...
func (r *Sync) Start() error {
	if r.cmd != nil {
		return fmt.Errorf("command already %v", "started")
	}
	if r.BinPath == "" {
		r.BinPath = rclonePath
	}
	r.cmd = exec.Command(r.BinPath, r.args()...)
	r.cmd.Stdout = r.Stdout
	stderr, err := r.cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := r.cmd.Start(); err != nil {
		return err
	}

	r.decoder = parser.NewDecoder(stderr)

	return nil
}

//Read ...
func (r *Sync) Read() ([]parser.TypedMessage, error) {
	return r.decoder.ReadObjects()
}

// encoder ...
type encoder interface {
	Encode(s interface{}) error
}

//ConvertTo ...
func (r *Sync) ConvertTo(w encoder) error {
	for {
		olist, err := r.Read()

		for _, o := range olist {
			if err2 := w.Encode(o); err2 != nil {
				return err
			}
		}

		if io.EOF == err {
			break
		} else {
			return err
		}
	}
	return nil
}

//Wait ...
func (r *Sync) Wait() error {
	if r.cmd == nil {
		return fmt.Errorf("command not %v", "started")
	}
	err := r.cmd.Wait()
	if r.stderr != nil {
		r.stderr.Close()
	}
	// r.cmd = nil
	// r.decoder = nil
	return err
}

//Kill ...
func (r *Sync) Kill() error {
	if r.cmd == nil {
		return fmt.Errorf("command not %v", "started")
	}
	if r.cmd.Process == nil {
		return fmt.Errorf("command already %v", "finished")
	}
	err := r.cmd.Process.Kill()
	// r.cmd = nil
	// r.decoder = nil
	return err
}
