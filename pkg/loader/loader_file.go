package loader

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/vgarvardt/rklotz/pkg/formatter"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// FileLoader is the Loader implementation for local file system
type FileLoader struct {
	path   string
	f      formatter.Formatter
	logger *slog.Logger
}

// NewFileLoader creates new FileLoader instance
func NewFileLoader(path string, f formatter.Formatter, logger *slog.Logger) (*FileLoader, error) {
	return &FileLoader{path, f, logger}, nil
}

// Load loads posts and saves them one by one in the storage
func (l *FileLoader) Load(s storage.Storage) error {
	if err := filepath.Walk(l.path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !f.IsDir() {
			l.logger.Debug("Loading post from file", slog.String("path", path))
			post, err := model.NewPostFromFile(l.path, path, l.f)
			if err != nil {
				return err
			}

			l.logger.Debug("Saving post to storage", slog.String("path", post.Path), slog.String("title", post.Title))
			err = s.Save(post)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return s.Finalize()
}
