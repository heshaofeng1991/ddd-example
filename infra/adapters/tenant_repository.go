package adapters

import (
	"context"
	"strconv"

	domainTenant "github.com/heshaofeng1991/ddd-johnny/domain/tenant"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/pkg/errors"
)

type TenantRepository struct {
	entClient *ent.Client
}

func NewTenantRepository(entClient *ent.Client) *TenantRepository {
	return &TenantRepository{entClient: entClient}
}

func (t TenantRepository) Create(ctx context.Context) (*domainTenant.Tenant, error) {
	tenant, err := t.entClient.Tenant.Create().Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tenant")
	}

	tenant, err = t.entClient.Tenant.UpdateOneID(tenant.ID).SetCode("NSS" + strconv.Itoa(int(tenant.ID))).Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update tenant code")
	}

	result, err := domainTenant.UnmarshalTenantFromDatabase(tenant.ID, tenant.Code)

	return result, errors.Wrap(err, "")
}

func (t TenantRepository) Save(ctx context.Context, tenant *domainTenant.Tenant) (*domainTenant.Tenant, error) {
	// implement me.
	panic("implement me")
}
