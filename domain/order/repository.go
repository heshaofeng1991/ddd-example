package order

import (
	"context"

	"github.com/heshaofeng1991/ddd-johnny/domain/store"
	"github.com/heshaofeng1991/ddd-johnny/domain/tenant"
)

type Repository interface {
	UpdateTenant(ctx context.Context, store *store.Store, tenant *tenant.Tenant) error
}
