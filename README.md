# SDK

### The SDK has moved to a public repo at [github.com/chainguard-dev/sdk](https://github.com/chainguard-dev/sdk)!

For guidance working with the SDK, post questions to
[#eng-squad-experience](https://chainguard-dev.slack.com/archives/C03JLFATBST) in Slack.

## Recommended workflow for making SDK changes

The workflow will be different since changes to SDK are not immediately available in `mono`. This recommended workflow
is still a work in progress and subject to change as we discover new pain points in the process.

1. Determine if your change is appropriate for `chainguard.dev/sdk` (i.e. is it okay to be public), or would
`mono/api-internal` be a better fit?
2. Make changes to your local copy of `chainguard.dev/sdk`. If you are developing something in `mono` that relies
on those changes, it may be helpful to add a temporary `replace chainguard.dev/sdk => [path to your local SDK repo]`
in the relevant `go.mod` for faster iteration.
3. Push `chainguard.dev/sdk` changes and open a PR in [github.com/chainguard-dev/sdk](https://github.com/chainguard-dev/sdk)
and post in [#eng-squad-experience](https://chainguard-dev.slack.com/archives/C03JLFATBST).
or tag relevant reviewers.
4. Cut a new release, update necessary `go.mod` of dependencies. The release cadence for SDK is still TBD, at the moment
ad hoc releases as necessary are fine. If you're unsure, ask in [#eng-squad-experience](https://chainguard-dev.slack.com/archives/C03JLFATBST).
