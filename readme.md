# Reverse Proxy Load Balancer

## What is this?

With this project I wanted to have two things:

1. A simple way to load balance traffic between my bare metal kubernetes cluster nodes
1. A fun go project to work on

## How to use it?

```
docker pull ghcr.io/dalibormesaric/rplb:latest

docker run -d --rm -p 8000:8000 -p 8080:8080 -e FE=localhost,myapp -e BE=myapp,http://192.168.1.1:80 --memory="64m" --memory-reservation="64m" --cpus=".1" ghcr.io/dalibormesaric/rplb:latest
```

dashboard
localhost:8000

reverse proxy
localhost:8080

## Misc

### Tools used

- https://coolors.co/palette/ef476f-ffd166-06d6a0
- https://cssgradient.io/

### TODO:

- [ ] Fix container stop after homeassistant upgrade?
- [ ] Fix Monitor pings when backend is not reachable

### Some ideas

- monitor -> uptime?
- https://stackoverflow.com/questions/16512840/get-domain-name-from-ip-address-in-go
- [ ] https://stackoverflow.com/questions/70442770/infinite-scrolling-carousel-css-only
- [ ] 503 Service Unavailable - when backends are unavailable
- [ ] 404 Page Not Found - when no frontends

### Docker
- https://github.com/docker/setup-buildx-action