# Reverse Proxy Load Balancer

## What is this?

With this project I wanted to have two things:

1. A simple way to load balance traffic between bare metal kubernetes cluster nodes
1. A fun Go project to work on

## How to use it?

### Docker

``` sh
docker pull ghcr.io/dalibormesaric/rplb:latest

docker run -d --rm -p 8000:8000 -p 8080:8080 -e FE=localhost,myapp -e BE=myapp,http://192.168.1.1:80 --memory="64m" --memory-reservation="64m" --cpus="1" ghcr.io/dalibormesaric/rplb:latest
```

### Reverse Proxy

http://localhost:8080

### Dashboard

http://localhost:8000

![monitor](/docs/monitor.png)

![traffic](/docs/traffic.png)

![traffic gif](/docs/traffic.gif)

## Example

``` sh
docker compose -f example/compose.yaml up rplb --build

for i in {1..10}; do curl -s localhost:8080 | grep h1; sleep 1; done;

docker compose -f example/compose.yaml down
```

## Testing data races

``` sh
docker compose -f example/compose.race.yaml up rplb --build

seq 1000 | parallel -n0 -j8 "curl -s http://localhost:8080 | grep h1"

docker compose -f example/compose.race.yaml down
```

## Misc

### Tools used

- https://coolors.co/palette/ef476f-ffd166-06d6a0
- https://cssgradient.io/

### TODO:

- [ ] Retries
- [ ] Load Balancing Strategies

### Some ideas

- [ ] https://stackoverflow.com/questions/16512840/get-domain-name-from-ip-address-in-go
- [ ] https://stackoverflow.com/questions/70442770/infinite-scrolling-carousel-css-only

### Docker
- https://github.com/docker/setup-buildx-action