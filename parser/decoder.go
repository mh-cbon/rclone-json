// Package parser decodes an rclone output to Objects.
package parser

// tbd: add test for output parsing.

import (
	"io"
	"regexp"
)

// TypeRawMessage is the type identifier for a Raw message.
var TypeRawMessage = "Raw"

// TypeGeneralStatMessage is the type identifier for a GeneralStat message.
var TypeGeneralStatMessage = "GeneralStat"

// TypeFileStatMessage is the type identifier for a FileStat message.
var TypeFileStatMessage = "FileStat"

// TypedMessage ...
type TypedMessage interface {
	GetType() string
	Compare(TypedMessage) bool
}

// Decoder ...
type Decoder struct {
	*LineReader
	line string
	gm   *GeneralStatMessage
}

// NewDecoder ...
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{LineReader: NewLineReader(r)}
}

// ReadObjects ...
func (l *Decoder) ReadObjects() ([]TypedMessage, error) {

	var ret []TypedMessage

	line, err := l.ReadLine()
	if err == io.EOF {
		if l.gm != nil {
			ret = append(ret, l.gm)
			l.gm = nil
		}
	}

	if line != "" {
		line = clean(line)
		if !isIgnored(line) {
			if raw, ok := isRawMessage(line); ok {
				if l.gm != nil {
					ret = append(ret, l.gm)
					l.gm = nil
				}
				ret = append(ret, &raw)
			}
			if ok, g := isGeneralStatMessage(line, l.gm); ok {
				l.gm = g
				return nil, nil
			}
			if f, ok := isFileStatMessage(line); ok {
				if l.gm != nil {
					ret = append(ret, l.gm)
					l.gm = nil
				}
				ret = append(ret, &f)
			}
		}
	}

	// if len(ret) == 0 {
	// 	fmt.Printf("==== %v\n", line)
	// }

	return ret, err
}

var logReg = regexp.MustCompile(`^.+ (DEBUG|INFO )\s+:\s+(.*)$`)

func clean(line string) string {
	if logReg.MatchString(line) {
		r := logReg.FindAllStringSubmatch(line, -1)
		if len(r) == 0 {
			return ""
		}
		// logLevel := r[0][1]
		return r[0][2]
	}
	return line
}

var ignoreReg = map[string]*regexp.Regexp{
	"transferring": regexp.MustCompile(`^Transferring:\s*$`),
}

func isIgnored(line string) bool {
	for _, reg := range ignoreReg {
		if reg.MatchString(line) {
			return true
		}
	}
	return false
}

var filesReg = map[string]*regexp.Regexp{
	"timethesame": regexp.MustCompile(`([^:]+):\s+Size and modification time the same.+$`),
	"notmodified": regexp.MustCompile(`([^:]+):\s+Unchanged skipping$`),
	"sizesdiffer": regexp.MustCompile(`([^:]+):\s+Sizes differ$`),
	"copied":      regexp.MustCompile(`([^:]+): Copied \(new\)$`),
}
var localReg = map[string]*regexp.Regexp{
	"waitingchecks":    regexp.MustCompile(`^Local file system at ([^:]+): Waiting for checks to finish$`),
	"waitingtransfers": regexp.MustCompile(`^Local file system at\s+([^:]+):\s+Waiting for transfers to finish$`),
	"modifywindow":     regexp.MustCompile(`^Local file system at\s+([^:]+):\s+Modify window is (.+)$`),
}

var infoReg = map[string]*regexp.Regexp{
	"bwlimit": regexp.MustCompile(`^Starting bandwidth limiter at (.+)$`),
	"version": regexp.MustCompile(`^rclone: Version "([^"]+)"`),
}

// RawMessage ...
type RawMessage struct {
	Type    string
	Rule    string
	Message string
	Info    string
}

// GetType ...
func (r RawMessage) GetType() string {
	return r.Type
}

// Compare ...
func (r RawMessage) Compare(left TypedMessage) bool {
	if x, ok := left.(RawMessage); ok {
		return r.Type == x.Type && r.Rule == x.Rule
	}
	return false
}

func isRawMessage(line string) (RawMessage, bool) {
	for name, reg := range filesReg {
		if reg.MatchString(line) {
			t := reg.FindAllStringSubmatch(line, -1)
			filePath := t[0][1]
			return RawMessage{Type: TypeRawMessage, Rule: name, Message: filePath}, true
		}
	}
	for name, reg := range localReg {
		if reg.MatchString(line) {
			t := reg.FindAllStringSubmatch(line, -1)
			filePath := t[0][1]
			msg := RawMessage{Type: TypeRawMessage, Rule: name, Message: filePath}
			if len(t[0]) > 1 {
				msg.Info = t[0][1]
			}
			return msg, true
		}
	}
	for name, reg := range infoReg {
		if reg.MatchString(line) {
			t := reg.FindAllStringSubmatch(line, -1)
			filePath := t[0][1]
			return RawMessage{Type: TypeRawMessage, Rule: name, Message: filePath}, true
		}
	}
	return RawMessage{}, false
}

var statsReg = map[string]*regexp.Regexp{
	"transferred":    regexp.MustCompile(`^Transferred:\s+([^\s]+)\s+([^\s]+)\s+\(([^)]+)\)$`),
	"errorsCnt":      regexp.MustCompile(`^Errors:\s+(\d+)$`),
	"checksCnt":      regexp.MustCompile(`^Checks:\s+(\d+)$`),
	"transferredCnt": regexp.MustCompile(`^Transferred:\s+(\d+)$`),
	"elapsedTime":    regexp.MustCompile(`^Elapsed time:\s+([^\s]+)$`),
}

// GeneralStatMessage ...
type GeneralStatMessage struct {
	Type                string
	TotalTransferredCnt string
	TotalTransferred    string
	TotalSpeed          string
	ErrorsCnt           string
	ChecksCnt           string
	TransferredCnt      string
	ElapsedTime         string
}

// GetType ...
func (g GeneralStatMessage) GetType() string {
	return g.Type
}

// Compare ...
func (g GeneralStatMessage) Compare(left TypedMessage) bool {
	if x, ok := left.(GeneralStatMessage); ok {
		return g.Type == x.Type
	}
	return false
}

func isGeneralStatMessage(line string, g *GeneralStatMessage) (bool, *GeneralStatMessage) {
	for name, reg := range statsReg {
		if reg.MatchString(line) {
			t := reg.FindAllStringSubmatch(line, -1)
			if name == "transferred" {
				if g == nil {
					g = &GeneralStatMessage{Type: TypeGeneralStatMessage}
				}
				// transferred: [][]string{[]string{"Transferred:   2.012 MBytes (630.953 kBytes/s)", "2.012", "MBytes", "630.953 kBytes/s"}}
				g.TotalTransferred = t[0][1] + " " + t[0][2]
				g.TotalSpeed = t[0][3]
			} else if name == "elapsedTime" {
				if g == nil {
					g = &GeneralStatMessage{Type: TypeGeneralStatMessage}
				}
				// elapsedTime: [][]string{[]string{"Elapsed time:        3.2s", "3.2s"}}
				g.ElapsedTime = t[0][1]
			} else if name == "transferredCnt" {
				if g == nil {
					g = &GeneralStatMessage{Type: TypeGeneralStatMessage}
				}
				// transferredCnt: [][]string{[]string{"Transferred:            0", "0"}}
				g.TotalTransferredCnt = t[0][1]
			} else if name == "checksCnt" {
				if g == nil {
					g = &GeneralStatMessage{Type: TypeGeneralStatMessage}
				}
				// checksCnt: [][]string{[]string{"Checks:                11", "11"}}
				g.ChecksCnt = t[0][1]
			} else if name == "errorsCnt" {
				if g == nil {
					g = &GeneralStatMessage{Type: TypeGeneralStatMessage}
				}
				// errorsCnt: [][]string{[]string{"Errors:                 0", "0"}}
				g.ErrorsCnt = t[0][1]
			}
			return true, g
		}
	}
	return false, g
}

var filestatsReg = map[string]*regexp.Regexp{
	"filestat": regexp.MustCompile(`^\s+\*\s+"([^"]+)":\s+([^%]+%)\s+([^,]+),\s+([\d.]+)\s+([^,]+),\s+ETA:\s+(.+)$`),
}

// FileStatMessage ...
type FileStatMessage struct {
	Type    string
	File    string
	Percent string
	Status  string
	Speed   string
	ETA     string
}

// GetType ...
func (f FileStatMessage) GetType() string {
	return f.Type
}

// Compare ...
func (f FileStatMessage) Compare(left TypedMessage) bool {
	if x, ok := left.(FileStatMessage); ok {
		return f.Type == x.Type && f.File == x.File
	}
	return false
}

func isFileStatMessage(line string) (FileStatMessage, bool) {
	ret := FileStatMessage{Type: TypeFileStatMessage}
	for _, reg := range filestatsReg {
		if reg.MatchString(line) {
			t := reg.FindAllStringSubmatch(line, -1)
			ret.File = t[0][1]
			ret.Percent = t[0][2]
			ret.Status = t[0][3]
			ret.Speed = t[0][4] + " " + t[0][5]
			ret.ETA = t[0][6]
			return ret, true
		}
	}
	return ret, false
}
