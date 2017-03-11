package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mh-cbon/rclone-json/parser"
	"github.com/mh-cbon/rclone-json/rclone"
)

type rCloneTests struct {
	Cmd              *rclone.Command
	ExpectedPaths    []string
	ExpectedMessages []parser.TypedMessage
}

func TestMain(t *testing.T) {

	base, _ := os.Getwd()

	tests := []rCloneTests{
		rCloneTests{
			Cmd: &rclone.Command{
				Stdout:        os.Stderr,
				Src:           "test/source",
				Dst:           "test/dest",
				Verbose:       true,
				Stats:         "1s",
				BwLimit:       "3M",
				Checkers:      "2",
				TransferLimit: "4",
			},
			ExpectedPaths: []string{
				"test/dest/file1",
				"test/dest/file2",
				"test/dest/file3",
				"test/dest/file4",
				"test/dest/folder/emptyfile1",
				"test/dest/folder/file",
				"test/dest/folder2/some/other/emptyfile2",
				"test/dest/folder2/some/other/what",
			},
			ExpectedMessages: []parser.TypedMessage{
				parser.RawMessage{Type: "RawMessage", Rule: "bwlimit", Message: "5MBytes/s"},
				parser.RawMessage{Type: "RawMessage", Rule: "version", Message: "v1.35-DEV"},
				parser.RawMessage{Type: "RawMessage", Rule: "modifywindow", Message: filepath.Join(base, "/test/dest")},
				parser.RawMessage{Type: "RawMessage", Rule: "waitingchecks", Message: filepath.Join(base, "/test/dest")},
				parser.RawMessage{Type: "RawMessage", Rule: "waitingtransfers", Message: filepath.Join(base, "/test/dest")},
				parser.GeneralStatMessage{Type: "GeneralStat"},
				parser.FileStatMessage{Type: "FileStat", File: "file1"},
				parser.FileStatMessage{Type: "FileStat", File: "file2"},
				parser.FileStatMessage{Type: "FileStat", File: "file3"},
				parser.FileStatMessage{Type: "FileStat", File: "file4"},
				parser.GeneralStatMessage{Type: "GeneralStat"},
				parser.FileStatMessage{Type: "FileStat", File: "file1"},
				parser.FileStatMessage{Type: "FileStat", File: "file2"},
				parser.FileStatMessage{Type: "FileStat", File: "file3"},
				parser.FileStatMessage{Type: "FileStat", File: "file4"},
				parser.RawMessage{Type: "RawMessage", Rule: "copied", Message: "file1"},
				parser.RawMessage{Type: "RawMessage", Rule: "copied", Message: "folder/emptyfile1"},
				parser.GeneralStatMessage{Type: "GeneralStat"},
				parser.FileStatMessage{Type: "FileStat", File: "file2"},
				parser.FileStatMessage{Type: "FileStat", File: "file3"},
				parser.FileStatMessage{Type: "FileStat", File: "file4"},
				parser.FileStatMessage{Type: "FileStat", File: "folder/file"},
				parser.RawMessage{Type: "RawMessage", Rule: "copied", Message: "file2"},
				parser.RawMessage{Type: "RawMessage", Rule: "copied", Message: "folder2/some/other/emptyfile2"},
				parser.GeneralStatMessage{Type: "GeneralStat"},
				parser.FileStatMessage{Type: "FileStat", File: "file3"},
				parser.FileStatMessage{Type: "FileStat", File: "file4"},
				parser.FileStatMessage{Type: "FileStat", File: "folder/file"},
				parser.FileStatMessage{Type: "FileStat", File: "folder2/some/other/what"},
				parser.GeneralStatMessage{Type: "GeneralStat"},
				parser.FileStatMessage{Type: "FileStat", File: "file3"},
				parser.FileStatMessage{Type: "FileStat", File: "file4"},
				parser.FileStatMessage{Type: "FileStat", File: "folder/file"},
				parser.FileStatMessage{Type: "FileStat", File: "folder2/some/other/what"},
				parser.RawMessage{Type: "RawMessage", Rule: "copied", Message: "folder/file"},
				parser.GeneralStatMessage{Type: "GeneralStat"},
				parser.FileStatMessage{Type: "FileStat", File: "file3"},
				parser.FileStatMessage{Type: "FileStat", File: "file4"},
				parser.FileStatMessage{Type: "FileStat", File: "folder2/some/other/what"},
				parser.RawMessage{Type: "RawMessage", Rule: "copied", Message: "file3"},
				parser.GeneralStatMessage{Type: "GeneralStat"},
				parser.FileStatMessage{Type: "FileStat", File: "file4"},
				parser.FileStatMessage{Type: "FileStat", File: "folder2/some/other/what"},
				parser.RawMessage{Type: "RawMessage", Rule: "copied", Message: "folder2/some/other/what"},
				parser.RawMessage{Type: "RawMessage", Rule: "copied", Message: "file4"},
				parser.GeneralStatMessage{Type: "GeneralStat"},
				parser.RawMessage{Type: "RawMessage", Rule: "version", Message: "v1.35-DEV"},
			},
		},
	}

	for _, test := range tests {

		mustNotErr(os.RemoveAll(test.Cmd.Dst))
		mustNotErr(os.MkdirAll(test.Cmd.Dst, os.ModePerm))

		mustNotErr(test.Cmd.Start())
		var objects []parser.TypedMessage
		for {
			list, err := test.Cmd.Read()
			objects = append(objects, list...)
			if io.EOF == err {
				break
			}
			mustNotErr(err)
		}
		mustNotErr(test.Cmd.Wait())

		notFound := []parser.TypedMessage{}
		for _, expected := range test.ExpectedMessages {
			// this might be tricky, under some circumstances,
			// it might returns different messages
			if mustContain(objects, expected) == false {
				notFound = append(notFound, expected)
			}
		}

		if len(notFound) > 0 {
			t.Log("Did not find those objects \n")
			for _, n := range notFound {
				t.Logf("%#v\n", n)
			}
			t.Log(" \n Got those obects \n")
			for _, n := range objects {
				t.Logf("%#v\n", n)
			}
			t.Error("FAILED")
			break
		}

		for _, path := range test.ExpectedPaths {
			if !fileMustExist(t, path) {
				t.Errorf("Dest file is missing %q\n", path)
				break
			}
		}
	}
}

func mustContain(list []parser.TypedMessage, search parser.TypedMessage) bool {
	for _, item := range list {
		if item.Compare(search) {
			return true
		}
	}
	return false
}

func compareTypedMessage(left parser.TypedMessage, right parser.TypedMessage) bool {
	switch left.(type) {
	case *parser.RawMessage:
		switch right.(type) {
		case *parser.RawMessage:
			return left.Compare(right)
		}
	case *parser.FileStatMessage:
		switch right.(type) {
		case *parser.FileStatMessage:
			return left.Compare(right)
		}
	case *parser.GeneralStatMessage:
		switch right.(type) {
		case *parser.GeneralStatMessage:
			return left.Compare(right)
		}
	}
	return false
}

func fileMustExist(t *testing.T, path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Dest file is missing %q\n", path)
		return false
	}
	return true
}
