package repository

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/vgarvardt/rklotz/pkg/model"
)

type FileLoader struct {
	path string
}

func NewFileLoader(path string) (*FileLoader, error) {
	return &FileLoader{path}, nil
}

func (l *FileLoader) Load(storage Storage) error {
	err := filepath.Walk(l.path, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			log.WithField("path", path).Debug("Loading post from file")
			post, err := model.NewPostFromFile(l.path, path)
			if err != nil {
				return err
			}

			log.WithFields(log.Fields{"path": post.Path, "title": post.Title}).
				Debug("Saving post to storage")
			err = storage.Save(post)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
