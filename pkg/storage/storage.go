package storage

import (
	"github.com/autonomy/devise/pkg/storage/datastore"
	"github.com/autonomy/devise/pkg/storage/datastore/memory"
)

var datastores = map[string]func() datastore.Datastore{
	"memory": func() datastore.Datastore { return memory.New() },
}

// NewDatastore instantiates and returns a storage datastore
func NewDatastore(b string) datastore.Datastore {
	return datastores[b]()
}
