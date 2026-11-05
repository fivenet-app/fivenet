// Package properties owns the shared dynamic-property contract for user and
// vehicle resources: definitions, typed values, lifecycle state, and catalog
// resolution.
//
// Definitions are schemas. System definitions are code-owned; custom
// definitions are generic and persisted by scope. Scope, key, and value type
// identify a custom definition and are not changed by updates. Incompatible
// changes require a new key and an explicit value migration.
//
// Disabled definitions reject new values and lifecycle patches but permit
// scalar clears. Soft deletion hides a definition and retains its values and
// state; restoration exposes them again. Neither operation purges data.
//
// State policy controls reason and expiry support. State patches change only
// supplied fields; a state-only mutation requires an existing scalar value,
// and clearing a scalar value removes its state. Policy changes do not migrate
// existing state.
//
// Generic expiry cleanup considers enabled custom definitions with a
// CLEAR_VALUE expiry policy. Managed properties use dedicated handlers.
package properties
