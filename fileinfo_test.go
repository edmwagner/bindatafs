package bindatafs_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/yookoala/bindatafs"
)

func Test_DirInfo(t *testing.T) {
	var i os.FileInfo = &bindatafs.DirInfo{}
	_ = i
	t.Log("*bindatafs.DirInfo{} implements os.FileInfo interface")
}

func TestFile(t *testing.T) {
	var f http.File = &bindatafs.File{}
	_ = f
	t.Logf("*bindatafs.File implements http.File interface")
}
