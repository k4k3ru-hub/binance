# Repository Guidelines

## Development Guidelines

### SDK Architecture: Root Facade + Composition Root + Manual Dependency Injection

When implementing or extending an SDK module for a platform, use **Root Facade + Composition Root + Manual Dependency Injection** as the default architecture.

#### Root Facade

- Provide the primary SDK entry points from the module root package.
- Use purpose-specific constructors such as `NewRESTClient` and `NewWebSocketClient` instead of forcing clients with different lifecycles into one aggregate client.
- Re-export commonly used public types with Go type aliases so typical usage requires only the root package import.
- Alias client options, request parameter types, top-level response types, and public error types when they are part of the normal consumer workflow.
- Use descriptive alias names such as `AllMidsParams`; do not expose ambiguous root names such as `Params`.
- Keep advanced and low-level packages importable for consumers who need direct control.
- Do not re-export every internal or nested wire type. Add a root alias when consumers commonly name, construct, inspect, or use the type in a public API.
- Keep the root facade thin. It should delegate construction and behavior rather than duplicate business or transport logic.

#### Composition Root

- Assemble the SDK object graph in constructors at explicit composition boundaries, such as `rest.NewClient`.
- Compose API-group clients and operation clients in one discoverable place.
- Keep construction errors contextual and fail client creation when required dependencies cannot be composed.
- Do not spread dependency construction across request methods or create hidden dependencies from global state.
- When adding an API operation, update its API-group composition, root facade aliases where appropriate, documentation, and composition tests together.

#### Manual Dependency Injection

- Pass dependencies explicitly through constructors or narrowly scoped options.
- Define small consumer-owned interfaces at the point of use, such as `transport.Executor` or `HTTPClient`.
- Accept injectable transports, HTTP clients, clocks, signers, or other external dependencies when substitution improves testing or platform support.
- Prefer immutable per-request parameter values over storing request-specific mutable state in reusable clients.
- Do not introduce a dependency injection container or service locator unless the repository has a demonstrated need that explicit construction cannot reasonably satisfy.
- Do not access replaceable dependencies through package globals.

#### Public API and Compatibility

- Optimize the normal usage path for a single root package import.
- Preserve existing lower-level public packages unless a deliberate breaking release permits their removal.
- Treat root type aliases and constructors as public compatibility commitments.
- Avoid exposing third-party implementation types from the root facade when doing so would couple SDK consumers to an unstable dependency. Wrap those types behind SDK-owned abstractions when long-term compatibility matters.
- Keep REST, WebSocket, and other clients separate when their initialization requirements, lifecycles, or failure modes differ.

#### Verification

- Add a root-package test proving that each facade constructor composes its expected API groups.
- Add operation-level tests using injected fakes rather than live platform services.
- Update root-package examples and README snippets to use the facade API.
- Run `go test ./...`, `go vet ./...`, and `git diff --check` after SDK architecture changes.

This architecture is the default, not an absolute rule. Deviations are acceptable when a platform's protocol, lifecycle, performance constraints, or compatibility requirements make another structure materially clearer. Document the reason when deviating.
