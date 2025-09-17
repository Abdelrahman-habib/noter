package fileserver

import (
	"io/fs"
	"net/http"
	"path/filepath"
)

type NeuteredFileSystem struct {
	fs fs.FS
}

func NewFileServer(fs fs.FS) http.Handler {
	return http.FileServerFS(NeuteredFileSystem{fs})
}

func (nfs NeuteredFileSystem) Open(path string) (fs.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
