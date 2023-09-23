package mysql_locks

import (
	"context"
	"database/sql"
	mariadb_storage "github.com/storage-lock/go-mariadb-storage"
	mysql_storage "github.com/storage-lock/go-mysql-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

var sqlDbStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[*sql.DB, *sql.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[*sql.DB, *sql.DB]()

func NewMariadbLockBySqlDb(ctx context.Context, db *sql.DB, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetMariadbLockFactoryBySqlDb(ctx, db)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewMariadbLockBySqlDbWithOptions(ctx context.Context, db *sql.DB, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetMariadbLockFactoryBySqlDb(ctx, db)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetMariadbLockFactoryBySqlDb(ctx context.Context, db *sql.DB) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return sqlDbStorageLockFactoryBeanFactory.GetOrInit(ctx, db, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		connectionManager := mariadb_storage.NewMariadbConnectionManagerFromSqlDb(db)
		options := mysql_storage.NewMySQLStorageOptions().SetConnectionManager(connectionManager)
		storage, err := mysql_storage.NewMysqlStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
