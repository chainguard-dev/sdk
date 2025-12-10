# Registry

This directory contains examples demonstrating the use of the `chainguard.dev/sdk/auth` and `chainguard.dev/sdk/auth/ggcr` packages for interacting with the Chainguard registry (GGCR).

## Examples

### exchange

The `exchange` example demonstrates how to exchange a token for an assumed identity to access the registry.

### chainctl

The `chainctl` example shows how to use a token source backed by `chainctl` to access the registry using local credentials.

## Notes

- Theses examples are somewhat contrived, since users can generally access `cgr.dev/chainguard` repos without auth.
  However, if a token is presented auth checks are still done, so it still serves as a useful guide.
- Modify the constants (e.g., `sub`, `aud`) in the examples as needed for your specific use case.
