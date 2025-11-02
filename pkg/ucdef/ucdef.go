// Package ucdef defines use case definitions that is used across the application.
package ucdef

import "context"

type AuditInfo struct {
	Tags          []string
	AggregateID   *string
	AggregateName *string
}

type UserReadAction[I, O, S any] interface {
	OperationID() string

	Validate(context.Context, *I) (S, error)
	Execute(context.Context, S) (*O, error)
}

type UserWriteAction[I, O, S any] interface {
	OperationID() string

	Validate(context.Context, *I) (S, error)
	Execute(context.Context, S) (*O, error)

	AuditInfo(S) AuditInfo
}

type EventConsumer[E, S any] interface {
	OperationID() string

	EnsureIntegrity(context.Context, E) (S, error)
	Execute(context.Context, S) error

	AuditInfo(S) AuditInfo
}

type BackgroundTask[T, S any] interface {
	OperationID() string

	EnsureIntegrity(context.Context, T) (S, error)
	Execute(context.Context, S) error

	AuditInfo(S) AuditInfo
}
