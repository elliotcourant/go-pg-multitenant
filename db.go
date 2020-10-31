package multitenant

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
)

type (
	// Tenanted interfaces are objects that belong to a single tenant within the
	// database. When performing a safe query objects must be restricted to the
	// current tenant of that query. This is done by enforcing the tenantId on
	// the records being modified.
	Tenanted interface {
		SetTenantId(tenantId int64)
	}
)

// NewSafeQuery will create a new query object that is restricted to models and
// data within a single tenant.
func NewSafeQuery(txn *pg.Tx, tenantId int64, model interface{}) (*orm.Query, error) {
	// We can only enforce a tenant constraint if the model we are working with
	// is actually one that belongs to a single tenant.
	tenanted, ok := model.(Tenanted)
	if !ok {
		// If the model does not belong to a single tenant then return an error.
		return nil, errors.Errorf("cannot use %T in a safe query, model does not belong to a single tenant")
	}

	// Make sure that the tenantId on the provided model is set to what we were
	// give here. This will make sure that if the model is being inserted that
	// we will insert it with the specific tenant.
	tenanted.SetTenantId(tenantId)

	// Create our go-pg query.
	query := txn.Model(model)

	// We need the table for our custom filter.
	tableName := query.TableModel().Table().Alias

	// Build a special WHERE clause to restrict filterable queries to this
	// tenant.
	tenantWhereExpression := fmt.Sprintf(`"%s"."tenant_id" = ?`, tableName)

	// Add the WHERE tenant_id = X to the query. This will make sure that any
	// UPDATE, DELETE or SELECT query built from here after will always include
	// this query filter.
	query = query.Where(tenantWhereExpression, tenantId)

	// Return our new safe query object to the caller.
	return query, nil
}
