package adapters

import (
	"context"

	domainStore "github.com/heshaofeng1991/ddd-johnny/domain/store"
	domainTenant "github.com/heshaofeng1991/ddd-johnny/domain/tenant"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/gen/platformproduct"
	"github.com/pkg/errors"
)

type ListingRepository struct {
	entClient *ent.Client
}

func (l ListingRepository) UpdateTenant(
	ctx context.Context,
	store *domainStore.Store,
	tenant *domainTenant.Tenant,
) error {
	_, err := l.entClient.PlatformProduct.
		Update().
		Where(platformproduct.StoreID(store.StoreID())).
		SetTenantID(tenant.TenantID()).
		Save(ctx)

	return errors.Wrap(err, "")
}

func NewListingRepository(entClient *ent.Client) *ListingRepository {
	return &ListingRepository{entClient: entClient}
}
