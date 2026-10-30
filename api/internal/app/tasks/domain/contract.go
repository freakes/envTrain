package domain

import "context"

// The domain owns ports; adapters depend on these contracts.
type Repository interface {
	Create(context.Context, Task) error
	Get(context.Context, string) (Task, error)
	List(context.Context) ([]Task, error)
}
type Usecase interface {
	Create(context.Context, CreateRequest) (Task, error)
	Get(context.Context, string) (Task, error)
	List(context.Context) ([]Task, error)
}
type Delivery interface{ Expose() }
