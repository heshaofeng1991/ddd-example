package adapters

import (
	"context"

	domainStore "github.com/heshaofeng1991/ddd-johnny/domain/store"
	domainTenant "github.com/heshaofeng1991/ddd-johnny/domain/tenant"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/gen/store"
	"github.com/pkg/errors"
)

type StoreRepository struct {
	entClient *ent.Client
}

func NewStoreRepository(entClient *ent.Client) *StoreRepository {
	return &StoreRepository{entClient: entClient}
}

func (s StoreRepository) GetByStoreCode(ctx context.Context, storeCode string) (*domainStore.Store, error) {
	found, err := s.entClient.Store.Query().
		Where(store.StoreCode(storeCode)).
		WithTenant().
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, domainStore.ErrStoreNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get store")
	}

	result, err := domainStore.UnmarshalStoreFromDatabase(found.ID, found.Edges.Tenant.ID, found.StoreURL, nil)

	return result, errors.Wrap(err, "")
}

func (s StoreRepository) UpdateTenantAndClearStoreCode(
	ctx context.Context,
	store *domainStore.Store,
	tenant *domainTenant.Tenant,
) error {
	_, err := s.entClient.Store.
		UpdateOneID(store.StoreID()).
		SetTenantID(tenant.TenantID()).
		SetStoreCode("").
		Save(ctx)

	return errors.Wrap(err, "")
}
