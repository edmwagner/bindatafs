package bindatafs_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/go-serve/bindatafs"
	"github.com/go-serve/bindatafs/example"
)

func TestFileSystem(t *testing.T) {
	var httpFs http.FileSystem = bindatafs.New(nil, nil, nil)
	_ = httpFs // just to prove bindatafs.FileSystem implements http.FileSystem
}

func fileInfoEqual(src, target os.FileInfo) (err error) {
	if want, have := src.Name(), target.Name(); want != have {
		err = fmt.Errorf("Name(): expected %#v, got %#v", want, have)
		return
	}
	if want, have := src.IsDir(), target.IsDir(); want != have {
		err = fmt.Errorf("IsDir(): expected %#v, got %#v", want, have)
		return
	}
	if src.IsDir() {
		if want, have := int64(0), target.Size(); want != have {
			err = fmt.Errorf("Size(): expected %#v, got %#v", want, have)
			return
		}
		if want, have := os.ModeDir, target.Mode()&os.ModeType; want != have {
			err = fmt.Errorf("Mode():\nexpected %b\ngot      %b", want, have)
			return
		}
		if want, have := os.FileMode(0777), target.Mode()&os.ModePerm; want != have {
			err = fmt.Errorf("Mode():\nexpected %b\ngot      %b", want, have)
			return
		}
		if want, have := int64(0), target.ModTime().Unix(); want != have {
			err = fmt.Errorf("Modtime(): expected %#v, got %#v", want, have)
			return
		}
	} else {
		if want, have := src.Size(), target.Size(); want != have {
			err = fmt.Errorf("Size(): expected %#v, got %#v", want, have)
			return
		}
		if want, have := os.FileMode(0444), target.Mode()&os.ModePerm; want != have {
			err = fmt.Errorf("Mode():\nexpected %b\ngot      %b", want, have)
			return
		}
		if want, have := src.ModTime().Unix(), target.ModTime().Unix(); want != have {
			err = fmt.Errorf("Modtime(): expected %#v, got %#v", want, have)
			return
		}
	}
	return
}

func TestFileSystem_Open(t *testing.T) {
	fs := example.FileSystem()
	tests := []struct {
		desc string
		path string
	}{
		{
			desc: "test open file",
			path: "hello.txt",
		},
		{
			desc: "test open sub-directory file",
			path: "hello/world.txt",
		},
		{
			desc: "test open directory",
			path: "hello",
		},
	}

	for i, test := range tests {
		t.Logf("test %d: %s", i+1, test.desc)

		// get the file/dir in the bindatafs
		file, err := fs.Open(test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		// get the counter part in the source assets
		srcFile, err := os.Open("example/assets/" + test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		defer srcFile.Close()
		srcStat, err := srcFile.Stat()
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		//t.Logf("stat: %#v", stat)
		//t.Logf("srcStat: %#v", srcStat)

		if err := fileInfoEqual(srcStat, stat); err != nil {
			t.Errorf("stat not equal, %s", err.Error())
		}
		if got := stat.Sys(); got != nil {
			t.Errorf("Sys() expected nil, got %#v", got)
		}

	}
}

func TestFileSystem_Readdir(t *testing.T) {
	fs := example.FileSystem()
	tests := []struct {
		desc  string
		path  string
		files map[string]bool
	}{
		{
			desc: "open hello folder",
			path: "hello",
			files: map[string]bool{
				"bar.txt":   true,
				"world.txt": true,
			},
		},
	}

	for i, test := range tests {

		t.Logf("test %d: %s", i+1, test.desc)

		dir, err := fs.Open(test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		defer dir.Close()

		arrFileInfo, err := dir.Readdir(10)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		// get a list of files
		if want, have := len(test.files), len(arrFileInfo); want != have {
			t.Errorf("expected %d, got %d", want, have)
		}
		for j, file := range arrFileInfo {
			if _, ok := test.files[file.Name()]; !ok {
				t.Errorf("test %d: files %d: %s was not an expected file", i+1, j+1, file.Name())
			}
		}

		// try read pass limit
		_, err = dir.Readdir(10)
		if err == nil {
			t.Errorf("expected error after read pass, got nil")
		} else if err != io.EOF {
			t.Errorf("expected io.EOF, got %#v", err.Error())
		}

	}

}
