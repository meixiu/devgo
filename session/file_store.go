package session

import (
	"github.com/gorilla/sessions"
)

type (
	fileStore struct {
		*sessions.FilesystemStore
	}
)

func NewFileStore(path string, keyPairs ...[]byte) Store {
	return &fileStore{sessions.NewFilesystemStore(path, keyPairs...)}
}

func (this *fileStore) Options(options Options) {
	this.FilesystemStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
