package example

import "github.com/yookoala/bindatafs"

// FileSystem returns a Filesystem implementation for the given assets
func FileSystem() bindatafs.FileSystem {
	return bindatafs.New(AssetDir, AssetInfo, Asset)
}
