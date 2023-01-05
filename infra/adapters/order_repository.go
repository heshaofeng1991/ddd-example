package adapters

import (
	"context"

	domainStore "github.com/heshaofeng1991/ddd-johnny/domain/store"
	domainTenant "github.com/heshaofeng1991/ddd-johnny/domain/tenant"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/gen/order"
	"github.com/pkg/errors"
)

type OrderRepository struct {
	entClient *ent.Client
}

func (o OrderRepository) UpdateTenant(
	ctx context.Context,
	store *domainStore.Store,
	tenant *domainTenant.Tenant,
) error {
	_, err := o.entClient.Order.
		Update().
		Where(order.StoreID(store.StoreID())).
		SetTenantID(tenant.TenantID()).
		Save(ctx)

	return errors.Wrap(err, "")
}

func NewOrderRepository(entClient *ent.Client) *OrderRepository {
	return &OrderRepository{entClient: entClient}
}
