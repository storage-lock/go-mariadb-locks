package mysql_locks

import (
	"context"
	storage_lock_test_helper "github.com/storage-lock/go-storage-lock-test-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewMysqlLockByDsn(t *testing.T) {
	mysqlDsn := os.Getenv("STORAGE_LOCK_MARIA_DSN")
	assert.NotEmpty(t, mysqlDsn)

	factory, err := GetMariadbLockFactoryByDsn(context.Background(), mysqlDsn)
	assert.Nil(t, err)

	storage_lock_test_helper.PlayerNum = 10
	storage_lock_test_helper.EveryOnePlayTimes = 100
	storage_lock_test_helper.TestStorageLock(t, factory)
}
