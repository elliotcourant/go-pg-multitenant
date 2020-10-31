package multitenant

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
)

type (
	Tenant struct {
		TenantId int64  `pg:"tenant_id,pk,type:'bigserial'"`
		Name     string `pg:"name,notnull"`
	}

	Product struct {
		TenantId    int64   `pg:"tenant_id,pk,notnull,on_delete:RESTRICT"`
		Tenant      *Tenant // Add the foreign key for tenants.
		ProductId   int64   `pg:"product_id,pk,type:'bigserial'"`
		SKU         string  `pg:"sku,unique,notnull"`
		Description string  `pg:"description"`
	}
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func GivenMySchemaIsSetup(t *testing.T) {

}

func GivenMyDatabaseIsSeeded(t *testing.T) {

}

func GivenIHaveATenant(t *testing.T) (tenantId int64) {
	return 1
}

func GivenIHaveADatabaseTransaction(
	t *testing.T,
	callback func(t *testing.T, txn *pg.Tx),
) {
	callback(t, nil)
}
