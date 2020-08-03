// Copyright 1999-2020. Plesk International GmbH.

package upload

import (
	"os"
)

type UploadData interface {
	UploadFiles(file os.FileInfo, overwrite bool) error
}

type UploadItems struct {
	ClientRoot string
	ServerRoot string
	Items      []string
}
