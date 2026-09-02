package store

import (
	"log"
	"sync"
)

var (
	globalDB   *DB
	globalOnce sync.Once
	globalErr  error
)

// Init opens the database and runs migrations + seed data.
// It is safe to call multiple times; only the first call has effect.
func Init(dataDir string) error {
	globalOnce.Do(func() {
		globalDB, globalErr = Open(dataDir)
		if globalErr != nil {
			return
		}
		if globalErr = globalDB.SeedTemplates(); globalErr != nil {
			return
		}
		var created bool
		if created, globalErr = globalDB.EnsureDefaultAdmin(); globalErr != nil {
			return
		}
		if created {
			log.Printf("created bootstrap admin account %q with the default password; change it after first login",
				DefaultAdminUsername)
		}
	})
	return globalErr
}

// Get returns the singleton DB instance. Returns nil if Init was not called.
func Get() *DB { return globalDB }
