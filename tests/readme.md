# Tests

## Integration tests

Integration tests are used to test RPLB end to end. This is also black-box testing of the application.

### Running in GitHub actions

Locally we can use `host.docker.internal`, which is a Docker Desktop thing. To be able to reach backends in GitHub actions, suggestion is to use `172.17.0.1` which is docker's default network gateway. This way we can still have frontend resolve to localhost, and have a different backend host.

- https://forums.docker.com/t/host-docker-internal-seems-doesnt-work-with-ci-cd-github-action-linux/119558/2
- https://forums.docker.com/t/how-to-reach-localhost-on-host-from-docker-container/113321/2