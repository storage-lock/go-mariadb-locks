package mariadb_locks

import (
	storage_lock "github.com/storage-lock/go-storage-lock"
	"sync"
)

var GlobalMariadbLockFactory *MariadbLockFactory
var globalMariadbLockFactoryOnce sync.Once
var globalMariadbLockFactoryErr error

func InitGlobalMariadbLockFactory(dsn string) error {
	factory, err := NewMariadbLockFactory(dsn)
	if err != nil {
		return err
	}
	GlobalMariadbLockFactory = factory
	return nil
}

func NewMariadbLock(dsn string, lockId string) (*storage_lock.StorageLock, error) {
	globalMariadbLockFactoryOnce.Do(func() {
		globalMariadbLockFactoryErr = InitGlobalMariadbLockFactory(dsn)
	})
	if globalMariadbLockFactoryErr != nil {
		return nil, globalMariadbLockFactoryErr
	}
	return GlobalMariadbLockFactory.CreateLock(lockId)
}
