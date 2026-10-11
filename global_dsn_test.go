package mysql_locks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalDsnLockFactoryBeanFactory_InitiallyInitialized(t *testing.T) {
	assert.NotNil(t, dsnStorageLockFactoryBeanFactory)
}

func TestNewMariadbLockByDsn_InvalidDsn(t *testing.T) {
	_, err := NewMariadbLockByDsn(context.Background(), "invalid-dsn", "test-lock")
	assert.Error(t, err)
}

func TestGetMariadbLockFactoryByDsn_InvalidDsn(t *testing.T) {
	_, err := GetMariadbLockFactoryByDsn(context.Background(), "invalid-dsn")
	assert.Error(t, err)
}

func TestNewMariadbLockByDsnWithOptions_InvalidDsn(t *testing.T) {
	_, err := NewMariadbLockByDsnWithOptions(context.Background(), "invalid-dsn", nil)
	assert.Error(t, err)
}

