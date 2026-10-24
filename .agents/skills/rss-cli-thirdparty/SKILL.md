---
name: rss-cli-thirdparty
description: Use when implementing or reviewing an external-service provider under pkg/thirdparty in rss-cli. Use when adding external API methods, wire types, provider transport tests, or deciding whether protocol work belongs in cli or Handler.
---

# rss-cli Thirdparty

## Quick Start

Start from the provider package and keep a new protocol operation internal until its application capability is confirmed:

```go
func (c *cli) getGroups(ctx context.Context, req *groupsRequest) (*groupsResponse, error)
```

## Responsibility Map

| Surface        | Owns                                                          | Does not own                              |
| -------------- | ------------------------------------------------------------- | ----------------------------------------- |
| Caller         | Business sequencing and command policy                        | Raw endpoint assembly or wire types       |
| Handler        | Scenario capability, conversion, and provider-facing errors   | Duplicate HTTP plumbing or command policy |
| Raw API method | One request, path/query/body, protocol status and envelope    | Cross-page iteration or business workflow |
| Types          | Named protocol request/response mapping and field semantics   | Caller policy                             |
| Shared client  | Transport, context, decoding, and existing logging mechanisms | Provider business decisions               |

## Why

`pkg/thirdparty` owns external protocol behavior. Keep third-party concepts inside the provider package and expose scenario-oriented capabilities through the Handler when the application contract requires them.

## Workflow

1. Read the nearest `AGENTS.md`, provider docs, existing callers, and 2-3 comparable methods.
2. Confirm the external method, URI, auth, input, output, success status, and pagination from the provider reference.
3. Confirm the semantic operation name, request, response, and Handler capability before adding a new operation.
4. Trace `caller -> Handler -> conversion -> raw API -> shared client`; record validation ownership before editing.
5. Add named request/response mapping types in the provider package. Do not export raw protocol structures.
6. Implement one raw API method per protocol request with `ctx + req -> resp + err`; keep page retrieval separate from aggregation.
7. Reuse the existing REST client for transport, context, auth, decoding, and logging before adding helpers.
8. Add transport-stub tests for method, URL, query/form fields, sensitive field placement, response parsing, non-2xx, and invalid JSON.
9. Decide separately whether a Handler or command-layer change is approved. Run the scoped tests, build, lint, and diff checks.

## API Method Contract

New raw methods use one predictable shape:

```go
func (c *cli) operate(
    ctx context.Context,
    req *operateRequest,
) (*operateResponse, error)
```

- Use one named request pointer and one named response pointer plus error.
- Keep path, query, and body inputs in the request type.
- Use a named response even for an empty successful response.
- Do not use maps, anonymous structs, extra scalar parameters, or generic wrappers as the method contract.
- Preserve stable existing methods unless their migration is explicitly in scope.

## Protocol Checks

| Question                    | Required decision                                                                           |
| --------------------------- | ------------------------------------------------------------------------------------------- |
| Input validation            | Follow the owning layer; do not add validation to every layer                               |
| Page or complete collection | Raw method gets one page; Handler aggregates only when its capability promises completeness |
| Omitted, empty, or null     | Preserve distinctions required by the protocol                                              |
| Success or absence          | Follow documented statuses and provider envelopes                                           |
| Existing client sufficient  | Inspect and reuse `pkg/rest/client` before adding transport helpers                         |

## Anti-patterns

WRONG: expose raw endpoints or wire DTOs to callers, add a shared abstraction for one adapter, or hide protocol semantics behind generic HTTP helpers.

RIGHT: keep one named raw operation per request, convert at the Handler boundary, and reuse the existing shared client.

## Checklist

- [ ] Existing Handler, CLI, and public command behavior is unchanged unless explicitly requested.
- [ ] New raw methods have `context.Context`, one named request, one named response, and error.
- [ ] No new exported protocol method or cross-provider abstraction was introduced without approval.
- [ ] Request and response wire types remain package-private.
- [ ] Tests use transport stubs and do not require Docker or real credentials.
- [ ] Formatting, tests, build, lint, and diff checks pass.

## See Also

- `rss-cli-how-to`
- `rss-cli-testsuite-support`
- `rss-cli-audience-boundary`
