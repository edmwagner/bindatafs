package bindatafs

import (
	"bytes"
	"io"
	"net/http"
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
	path   string
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
	var childFile http.File

	exists := false
	for _, name := range names {
		if i >= f.dirPos {
			exists = true
			childFile, err = f.assets.Open(path.Join(f.path, f.name, string(name)))
			if err != nil {
				return
			}

			fi, err = childFile.Stat()
			if err != nil {
				return
			}
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

	// return standard dirInfo with given name
	if f.t == TypeDir {
		fi = &dirInfo{
			name: f.name,
		}
		return
	}

	// return fileInfo in wrapper
	wrapper := &fileInfo{
		name: f.name,
	}
	fi = wrapper

	wrapper.FileInfo, err = f.assets.assetInfo(path.Join(f.path, f.name))
	return
}

// fileInfo implements FileInfo
type fileInfo struct {
	name string
	os.FileInfo
}

// Name implements os.FileInfo
func (fi *fileInfo) Name() string {
	return fi.name
}

// Size gives length in bytes for regular files;
// system-dependent for others
func (fi *fileInfo) Size() int64 {
	return fi.FileInfo.Size()
}

// Mode gives file mode bits
func (fi *fileInfo) Mode() os.FileMode {
	return fi.FileInfo.Mode()&os.ModeType | 0444
}

// ModTime gives modification time
func (fi *fileInfo) ModTime() (t time.Time) {
	return fi.FileInfo.ModTime()
}

// IsDir is abbreviation for Mode().IsDir()
func (fi *fileInfo) IsDir() bool {
	return fi.Mode().IsDir()
}

// Sys gives underlying data source (can return nil)
func (fi *fileInfo) Sys() interface{} {
	return nil
}

// dirInfo implements FileInfo for directory in the assets
type dirInfo struct {
	name string
}

// Name gives base name of the file
func (fi *dirInfo) Name() string {
	return fi.name
}

// Size gives length in bytes for regular files;
// system-dependent for others
func (fi *dirInfo) Size() int64 {
	return 0 // hard code 0 for now (originally system-dependent)
}

// Mode gives file mode bits
func (fi *dirInfo) Mode() os.FileMode {
	return os.ModeDir | 0777
}

// ModTime gives modification time
func (fi *dirInfo) ModTime() (t time.Time) {
	return time.Unix(0, 0)
}

// IsDir is abbreviation for Mode().IsDir()
func (fi *dirInfo) IsDir() bool {
	return fi.Mode().IsDir()
}

// Sys gives underlying data source (can return nil)
func (fi *dirInfo) Sys() interface{} {
	return nil
}
