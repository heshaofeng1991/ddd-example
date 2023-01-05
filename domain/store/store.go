package store

import (
	"github.com/heshaofeng1991/common/util"
	domainTenant "github.com/heshaofeng1991/ddd-johnny/domain/tenant"
	"github.com/pkg/errors"
)

var ErrStoreNotFound = errors.New("store not found")

type Store struct {
	storeID     int64
	storeTenant int64
	storeName   string
	storeURL    string
	platform    string
	email       string
	storeCode   string

	tenant *domainTenant.Tenant
}

func (s Store) StoreTenant() int64 {
	return s.storeTenant
}

func (s Store) StoreName() string {
	return s.storeName
}

func (s Store) StoreURL() string {
	return s.storeURL
}

func (s Store) Platform() string {
	return s.platform
}

func (s Store) Email() string {
	return s.email
}

func (s Store) Tenant() *domainTenant.Tenant {
	return s.tenant
}

func (s Store) StoreID() int64 {
	return s.storeID
}

func (s Store) StoreCode() string {
	return s.storeCode
}

func NewStore(storeURL string, storeTenant int64, tenant *domainTenant.Tenant) (*Store, error) {
	if storeURL == "" {
		return nil, errors.New("store url is empty")
	}

	return &Store{
		storeURL:    storeURL,
		tenant:      tenant,
		storeTenant: storeTenant,
	}, nil
}

func UnmarshalStoreFromDatabase(storeID int64, storeTenant int64, storeURL string, tenant *domainTenant.Tenant) (
	*Store,
	error,
) {
	store, err := NewStore(storeURL, storeTenant, tenant)
	if err != nil {
		return nil, err
	}

	store.storeID = storeID

	return store, nil
}

func (s *Store) CanUpdateTenant() error {
	if s.StoreTenant() == util.SystemTenantID {
		return nil
	}

	return errors.New("can't update tenant")
}
