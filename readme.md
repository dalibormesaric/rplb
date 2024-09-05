# RPLB – Reverse Proxy Load Balancer

A simple application that can load balance requests based on configured hostname accross configured backends.

## Features

- ⚙️ Simple configuration
   - In-line 
- ⚡️ Static Load Balancing
   - Sticky Round Robin
   - Round Robin
   - Random
   - First
- Resilient
   - Retry
- 📈 Dashboard
   - Monitor
   - Traffic

![monitor](/docs/monitor.png)

![traffic](/docs/traffic.png)

## Getting started

You can run `RPLB` with these commands:

``` sh
docker pull ghcr.io/dalibormesaric/rplb:latest

docker run -d --rm -p 8000:8000 -p 8080:8080 -e FE=localhost,myapp -e BE=myapp,http://10.0.0.1:1234,myapp,http://10.0.0.2:1234,myapp,http://10.0.0.3:1234 --memory="64m" --memory-reservation="64m" --cpus="1" ghcr.io/dalibormesaric/rplb:latest
```

### Configuration

> FE=localhost,myapp
> - `localhost` is hostname where `RPLB` is running, so in this case you would access your backend via `http://localhost:8080` and the Dashboard via `http://localhost:8000`
> - `myapp` is name of the Backend Pool that this hostname is connected to

> BE=myapp,http://10.0.0.1:1234,myapp,http://10.0.0.2:1234,myapp,http://10.0.0.3:1234
> - `myapp` is name of the Backend Pool to which the URL is assigned
> - `http://10.0.0.1:1234`, `http://10.0.0.1:1234` and `http://10.0.0.3:1234` are the URLs of your application

### Home Assistant

## Try it out

There is an `/example` folder in this repository containing already configured `RPLB` with tree backends. You can try it our by running:

``` sh
docker compose -f example/compose.yaml up rplb --build

for i in {1..10}; do curl -s localhost:8080 | grep \<h1; sleep 1; done;

docker compose -f example/compose.yaml down
```

![traffic gif](/docs/traffic.gif)

## Development

``` sh
go run -ldflags "-X github.com/dalibormesaric/rplb/internal/config.Version=$(git describe --tags --abbrev=0)" cmd/rplb/main.go
```

``` sh
docker build --build-arg="VERSION=$(git describe --tags --abbrev=0)" -t rplb .
```

### Build Pipeline

- https://github.com/docker/setup-buildx-action

### Testing

### Testing data races

``` sh
docker compose -f example/compose.race.yaml up rplb --build

seq 1000 | parallel -n0 -j8 "curl -s http://localhost:8080 | grep \<h1"

docker compose -f example/compose.race.yaml down
```

## What is this?

With this project I wanted to have two things:

1. A simple way to load balance traffic between bare metal kubernetes cluster nodes
1. A fun Go project to work on

### Reverse Proxy

http://localhost:8080

### Dashboard

http://localhost:8000


## Misc

### Tools used

- https://coolors.co/palette/ef476f-ffd166-06d6a0
- https://cssgradient.io/

### TODO:

- [ ] docs
- [ ] /metrics

### Some ideas

- [ ] https://bazel-contrib.github.io/SIG-rules-authors/go-tutorial.html
- [ ] https://stackoverflow.com/questions/16512840/get-domain-name-from-ip-address-in-go
- [ ] https://stackoverflow.com/questions/70442770/infinite-scrolling-carousel-css-only

### Docker