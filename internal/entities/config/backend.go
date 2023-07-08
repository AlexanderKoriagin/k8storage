package config

import (
	"sync"
)

// StorageBackend is a struct that holds the configuration for the storage backend
type StorageBackend struct {
	SourceCfg  *Redis        // short term storage
	StorageCfg []*Redis      // long term storage nodes
	ChanStop   chan struct{} // channel to stop the storage backend worker
	Wg         *sync.WaitGroup
}
