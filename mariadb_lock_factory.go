package mariadb_locks

import (
	"context"
	"database/sql"
	mariadb_storage "github.com/storage-lock/go-mariadb-storage"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

type MariadbLockFactory struct {
	*storage_lock_factory.StorageLockFactory[*sql.DB]
}

func NewMariadbLockFactory(dsn string) (*MariadbLockFactory, error) {
	connectionManager := mariadb_storage.NewMariaStorageConnectionManagerFromDSN(dsn)

	options := mariadb_storage.NewMariaStorageOptions()
	options.SetConnectionProvider(connectionManager)
	storage, err := mariadb_storage.NewMariaDbStorage(context.Background(), options)
	if err != nil {
		return nil, err
	}

	factory := storage_lock_factory.NewStorageLockFactory[*sql.DB](storage, connectionManager)

	return &MariadbLockFactory{
		StorageLockFactory: factory,
	}, nil
}
