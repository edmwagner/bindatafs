package bindatafs_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/yookoala/bindatafs"
	"github.com/yookoala/bindatafs/example"
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
