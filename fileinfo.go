package bindatafs

import (
	"bytes"
	"io"
	"os"
	"path"
	"time"
)

// Type of an asset provided by FileSystem
type Type int

const (

	// TypeFile is the type "file"
	TypeFile Type = iota

	// TypeDir is the type "dir"
	TypeDir
)

// File implments http.File
type File struct {
	*bytes.Reader
	assets *generatedAssets
	name   string
	t      Type
	dirPos int
}

// Close is a dummy method to implment io.Closer
func (f *File) Close() error {
	return nil
}

// Readdir reads the contents of the directory associated with
// file and returns a slice of up to n FileInfo values, as would
// be returned by Lstat, in directory order. Subsequent calls on
// the same file will yield further FileInfos.
func (f *File) Readdir(count int) (lfi []os.FileInfo, err error) {
	names, err := f.assets.assetDir(f.name)
	lfi = make([]os.FileInfo, 0)

	var fi os.FileInfo
	var i, j int

	exists := false
	for _, name := range names {
		if i >= f.dirPos {
			exists = true
			fi, err = f.assets.assetInfo(path.Join(f.name, string(name)))
			lfi = append(lfi, fi)
			j++
		}

		i++
		if j == count {
			break
		}
	}

	f.dirPos += j
	if !exists {
		err = io.EOF
	}

	return
}

// Stat returns the FileInfo structure describing file. If
// there is an error, it will be of type *PathError.
func (f *File) Stat() (fi os.FileInfo, err error) {
	if f.t == TypeDir {
		fi = &DirInfo{}
		return
	}
	return f.assets.assetInfo(f.name)
}

// DirInfo implements FileInfo for directory in the assets
type DirInfo struct {
	name string
	size int64
}

// Name gives base name of the file
func (fi *DirInfo) Name() string {
	return fi.name
}

// Size gives length in bytes for regular files;
// system-dependent for others
func (fi *DirInfo) Size() int64 {
	return fi.size
}

// Mode gives file mode bits
func (fi *DirInfo) Mode() os.FileMode {
	return os.ModeDir
}

// ModTime gives modification time
func (fi *DirInfo) ModTime() (t time.Time) {
	return t
}

// IsDir is abbreviation for Mode().IsDir()
func (fi *DirInfo) IsDir() bool {
	return fi.Mode().IsDir()
}

// Sys gives underlying data source (can return nil)
func (fi *DirInfo) Sys() interface{} {
	return nil
}
