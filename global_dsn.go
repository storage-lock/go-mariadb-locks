package mysql_locks

import (
	"context"
	"database/sql"
	mariadb_storage "github.com/storage-lock/go-mariadb-storage"
	mysql_storage "github.com/storage-lock/go-mysql-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

var dsnStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[string, *sql.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[string, *sql.DB]()

func NewMariadbLockByDsn(ctx context.Context, dsn string, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetMariadbLockFactoryByDsn(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewMariadbLockByDsnWithOptions(ctx context.Context, uri string, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetMariadbLockFactoryByDsn(ctx, uri)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetMariadbLockFactoryByDsn(ctx context.Context, dsn string) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return dsnStorageLockFactoryBeanFactory.GetOrInit(ctx, dsn, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		connectionManager := mariadb_storage.NewMariadbConnectionManagerFromDsn(dsn)
		options := mysql_storage.NewMySQLStorageOptions().SetConnectionManager(connectionManager)
		storage, err := mysql_storage.NewMysqlStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
