package tenant

import "github.com/pkg/errors"

type Tenant struct {
	tenantID int64
	code     string
}

func NewTenant(tenantID int64, code string) (*Tenant, error) {
	if tenantID <= 0 {
		return nil, errors.New("customer id is invalid")
	}

	return &Tenant{tenantID: tenantID, code: code}, nil
}

func (t Tenant) TenantID() int64 {
	return t.tenantID
}

func (t Tenant) Code() string {
	return t.code
}

func UnmarshalTenantFromDatabase(tenantID int64, code string) (*Tenant, error) {
	if tenantID == 0 {
		return nil, errors.New("tenantID is required")
	}

	if code == "" {
		return nil, errors.New("code is required")
	}

	return &Tenant{tenantID: tenantID, code: code}, nil
}
