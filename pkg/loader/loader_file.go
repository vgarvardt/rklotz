package loader

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// FileLoader is the Loader implementation for local file system
type FileLoader struct {
	path string
}

// NewFileLoader creates new FileLoader instance
func NewFileLoader(path string) (*FileLoader, error) {
	return &FileLoader{path}, nil
}

// Load loads posts and saves them one by one in the storage
func (l *FileLoader) Load(storage storage.Storage) error {
	err := filepath.Walk(l.path, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if nil != err {
				return err
			}

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
	if err != nil {
		return err
	}

	return storage.Finalize()
}
