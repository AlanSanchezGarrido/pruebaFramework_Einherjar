package perm

import "code.nochebuena.dev/einherjar/contracts/security"

const (
	ReadCustomers = security.Permission (0)
	WriteCustomers = security.Permission (1)
	ReadShoppings = security.Permission(0)
	WriteShoppings = security.Permission(1)
)