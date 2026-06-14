# Backend DDD / Clean Architecture Boundaries

This document defines where backend domain rules belong. The goal is not to
increase the number of DDD-shaped types, but to keep business rules cohesive
while preserving the current Clean Architecture dependency direction.

## Dependency Direction

The backend keeps this dependency direction:

```text
handler -> usecase -> domain <- infra
```

`domain` must not import handler DTOs, OpenAPI generated types, PostgreSQL
adapters, sqlc generated code, Firebase clients, or other infrastructure
details. Infrastructure code implements domain interfaces.

## Layer Responsibilities

### Value Objects

Use `internal/domain/value` for constraints that can be validated from one
value alone:

- required / optional semantics
- length limits
- enum membership
- formatting
- normalization when it is part of the domain meaning

Examples: `URL`, `Source`, `Stage`, `TaskTitle`, `InboxClipTitle`.

### Entities

Use `internal/domain/entity` for identity, lifecycle state, and invariants that
an entity can enforce from its own state.

Entity methods should use domain language, such as `Complete`, `UpdateStage`,
or `Rename`. Constructors should require already-valid value objects rather
than raw request strings.

### Domain Services

Use `internal/domain/service` when a rule is domain behavior but does not fit a
single entity or value object.

Typical triggers:

- a rule needs more than one entity
- a rule needs a domain repository query
- duplicate handling is part of the domain behavior
- a transition policy coordinates multiple domain objects

Domain service names should be specific to the business operation. Prefer
`InboxClipRegistrationService` or `StageProgressionService` over broad names
such as `EntryService`.

Domain service methods should use domain verbs such as `Register`, `Advance`,
`Link`, or `Resolve`. Do not use `Execute` for domain services.

### Use Cases

Use `internal/usecase` for application flow:

- convert input primitives into value objects
- perform ownership checks required by the application boundary
- call entities, domain services, and repositories
- manage transaction or unit-of-work boundaries
- return output DTOs used by handlers or app adapters

Use case methods may continue to use `Execute`. In this codebase,
`Execute` means "run this application command/query"; it is not a domain verb.

### Repositories And Ports

Domain repositories belong under `internal/domain/repository` when the
interface represents persistence or lookup of domain objects and can be used by
domain rules.

Usecase ports belong under `internal/usecase/...` when the interface represents
an application operation, transaction boundary, external side effect, or read
model that is not itself a domain repository.

For example, a repository that loads `InboxClip` by `(userID, URL)` is a domain
repository. A unit of work that saves multiple aggregates in one transaction
for one application command should be considered a usecase port unless it is a
stable domain concept.

## Database Constraints

If a rule is a true domain invariant, enforce it in both domain/application code
and the database.

Application-side checks improve intent and error handling, but they do not
protect against concurrent requests. Unique rules such as "one user cannot
register the same inbox URL twice" require a database unique constraint or
equivalent atomic write.

Repository implementations must translate database constraint errors to domain
repository errors such as `repository.ErrAlreadyExists`.

## Current Policy Decisions

### InboxClip Registration

`InboxClip` registration is a domain operation. A user's repeated capture of
the same URL should return the existing clip instead of creating another one.

This rule belongs in `domain/service.InboxClipRegistrationService`, backed by
`repository.InboxClipRepository`. The database also enforces
`UNIQUE(user_id, url)` so concurrent captures cannot create duplicates.

### CompanyAlias Uniqueness

Company aliases are a user-owned dictionary. Within one user's company, the
same alias should not be stored more than once. This is enforced with
`UNIQUE(user_id, company_id, alias)` and by mapping unique violations to
`repository.ErrAlreadyExists`.

If future matching requires an alias to identify at most one company per user,
tighten this to `UNIQUE(user_id, alias)` in a separate migration with explicit
data migration.

### Stage History

Current behavior treats stage history as an explicit record operation. Updating
an `Entry` stage does not automatically create `StageHistory`.

If the product requirement changes to "every stage change must have history",
introduce `StageProgressionService` and a usecase transaction boundary that
saves `Entry` and creates `StageHistory` atomically.

### Account Linking

Authentication currently supports only Google. The find/link/create flow can
remain in the `user.Authenticate` usecase while the policy is simple.

If additional providers, link rejection rules, or email collision policies are
added, move the policy into an `AccountLinkingService`.
