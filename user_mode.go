package simutils

// UserMode ...
type UserMode uint

const (
	// PENDING ...
	PENDING UserMode = iota + 1
	// USER_STATUS_ACTIVE Active
	USER_STATUS_ACTIVE
	// USER_STATUS_INACTIVE Inactive
	USER_STATUS_INACTIVE
)
