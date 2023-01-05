package tenant

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context) (*Tenant, error)
	Save(ctx context.Context, tenant *Tenant) (*Tenant, error)
}
