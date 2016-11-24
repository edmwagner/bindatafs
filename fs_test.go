package bindatafs_test

import (
	"net/http"
	"testing"

	"github.com/yookoala/bindatafs"
)

func TestFileSystem(t *testing.T) {
	var httpFs http.FileSystem = bindatafs.New(nil, nil, nil)
	_ = httpFs // just to prove bindatafs.FileSystem implements http.FileSystem
}
