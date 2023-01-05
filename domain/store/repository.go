package store

import (
	"context"

	"github.com/heshaofeng1991/ddd-johnny/domain/tenant"
)

type Repository interface {
	GetByStoreCode(ctx context.Context, storeCode string) (*Store, error)
	UpdateTenantAndClearStoreCode(ctx context.Context, store *Store, tenant *tenant.Tenant) error
}
