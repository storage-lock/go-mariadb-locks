package mysql_locks

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	mysql_storage "github.com/storage-lock/go-mysql-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

func TestNewMariadbLockBySqlDb(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	// Init 建表
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()
	factory, err := GetMariadbLockFactoryBySqlDb(ctx, mockDB)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetMariadbLockFactoryBySqlDb(ctx, mockDB)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewMariadbLockBySqlDb(ctx, mockDB, "test-mariadb-lock-sqldb")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-mariadb-lock-sqldb-opt")
	lock2, err := NewMariadbLockBySqlDbWithOptions(ctx, mockDB, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestNewMariadbLockBySqlDb_Errors(t *testing.T) {
	old := newMysqlStorage
	newMysqlStorage = func(ctx context.Context, options *mysql_storage.MysqlStorageOptions) (*mysql_storage.MysqlStorage, error) {
		return nil, errors.New("boom")
	}
	defer func() { newMysqlStorage = old }()

	mockDB, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	ctx := context.Background()
	_, err = GetMariadbLockFactoryBySqlDb(ctx, mockDB)
	assert.NotNil(t, err)
	_, err = NewMariadbLockBySqlDb(ctx, mockDB, "x")
	assert.NotNil(t, err)
	_, err = NewMariadbLockBySqlDbWithOptions(ctx, mockDB, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}

func TestNewMariadbLockByDsn_Local(t *testing.T) {
	// MariaDB 协议兼容 MySQL，用本机 sl-mysql 验证成功路径
	dsn := "root:rootpass@tcp(127.0.0.1:13307)/storage_lock_test"
	ctx := context.Background()

	probe, err := GetMariadbLockFactoryByDsn(ctx, dsn)
	if err != nil {
		t.Skipf("本机 sl-mysql 不可用，跳过: %v", err)
	}
	assert.NotNil(t, probe)

	// 缓存命中
	factory2, err := GetMariadbLockFactoryByDsn(ctx, dsn)
	assert.Nil(t, err)
	assert.Equal(t, probe, factory2)

	lock, err := NewMariadbLockByDsn(ctx, dsn, "test-mariadb-lock-dsn")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-mariadb-lock-dsn-opt")
	lock2, err := NewMariadbLockByDsnWithOptions(ctx, dsn, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}
