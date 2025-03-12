package storage

import (
	"fmt"
	"notex/pkg/types"
)

// StorageFactory 存储工厂
type StorageFactory struct {
	providers map[StorageType]StorageProvider
}

// NewStorageFactory 创建存储工厂
func NewStorageFactory() *StorageFactory {
	return &StorageFactory{
		providers: make(map[StorageType]StorageProvider),
	}
}

// RegisterProvider 注册存储提供者
func (f *StorageFactory) RegisterProvider(storageType StorageType, provider StorageProvider) {
	f.providers[storageType] = provider
}

// CreateStorage 创建存储实例
func (f *StorageFactory) CreateStorage(config *types.StorageConfig) (Storage, error) {
	provider, ok := f.providers[StorageType(config.Type)]
	if !ok {
		return nil, fmt.Errorf("storage provider not found for type: %s", config.Type)
	}

	return provider.NewStorage(config)
}

// LocalStorageProvider 本地存储提供者
type LocalStorageProvider struct{}

func (p *LocalStorageProvider) NewStorage(config *types.StorageConfig) (Storage, error) {
	return NewLocalStorage(config), nil
}

// OSSStorageProvider OSS存储提供者
type OSSStorageProvider struct{}

func (p *OSSStorageProvider) NewStorage(config *types.StorageConfig) (Storage, error) {
	return NewOSSStorage(config)
}

// DefaultFactory 默认存储工厂实例
var DefaultFactory = NewStorageFactory()

func init() {
	// 注册默认的存储提供者
	DefaultFactory.RegisterProvider(StorageTypeLocal, &LocalStorageProvider{})
	DefaultFactory.RegisterProvider(StorageTypeOSS, &OSSStorageProvider{})
}
