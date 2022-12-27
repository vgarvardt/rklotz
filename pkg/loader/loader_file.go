package loader

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/vgarvardt/rklotz/pkg/formatter"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// FileLoader is the Loader implementation for local file system
type FileLoader struct {
	path   string
	f      formatter.Formatter
	logger *zap.Logger
}

// NewFileLoader creates new FileLoader instance
func NewFileLoader(path string, f formatter.Formatter, logger *zap.Logger) (*FileLoader, error) {
	return &FileLoader{path, f, logger}, nil
}

// Load loads posts and saves them one by one in the storage
func (l *FileLoader) Load(storage storage.Storage) error {
	err := filepath.Walk(l.path, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if nil != err {
				return err
			}

			l.logger.Debug("Loading post from file", zap.String("path", path))
			post, err := model.NewPostFromFile(l.path, path, l.f)
			if err != nil {
				return err
			}

			l.logger.Debug("Saving post to storage", zap.String("path", post.Path), zap.String("title", post.Title))
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
