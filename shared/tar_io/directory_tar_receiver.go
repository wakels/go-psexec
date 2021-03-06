package tar_io

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/golang-devops/go-psexec/shared/io_throttler"
)

type dirTarReceiver struct {
	dir string
}

func (d *dirTarReceiver) OnEntry(tarHeader *tar.Header, tarFileReader io.Reader) error {
	relativePath := tarHeader.Name

	if tarHeader.FileInfo().IsDir() {
		fullDestPath := filepath.Join(d.dir, relativePath)
		defer os.Chtimes(fullDestPath, tarHeader.AccessTime, tarHeader.ModTime)

		return os.MkdirAll(fullDestPath, os.FileMode(tarHeader.Mode))
	} else {
		fullDestPath := filepath.Join(d.dir, relativePath)
		if val, ok := tarHeader.Xattrs["SINGLE_FILE_ONLY"]; ok && val == "1" {
			fullDestPath = d.dir
		}

		parentDir := filepath.Dir(fullDestPath)
		err := os.MkdirAll(parentDir, os.FileMode(tarHeader.Mode))
		if err != nil {
			return fmt.Errorf("Unable to MkDirAll '%s', error: %s", parentDir, err.Error())
		}

		file, err := os.OpenFile(fullDestPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(tarHeader.Mode))
		if err != nil {
			return fmt.Errorf("Unable to open file '%s', error: %s", fullDestPath, err.Error())
		}

		defer func() {
			file.Close()
			os.Chtimes(fullDestPath, tarHeader.AccessTime, tarHeader.ModTime)
		}()

		_, err = io_throttler.CopyThrottled(io_throttler.DefaultIOThrottlingBandwidth, file, tarFileReader)
		if err != nil {
			return fmt.Errorf("Unable to copy stream to file '%s', error: %s", fullDestPath, err.Error())
		}

		return nil
	}
}
