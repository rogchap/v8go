# Contributing

## Releases

We use pre-release versions so that we can keep aligned with the upstream's semver. Use `-apN`:

```sh
export AIRPLANE_V8GO_TAG=v0.6.1-ap1 && \
  git tag ${AIRPLANE_V8GO_TAG} && \
  git push origin ${AIRPLANE_V8GO_TAG}
```
