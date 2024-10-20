# Load balancing

Here are short descriptions of implemented load balancing algorithms and why to use them.

## Random

Chooses at random some backend from the provided list.

### Why use it?

If you want to introduce some chaos.

## First

Chooses the first backend from the provided list.

### Why use it?

If you explicitly want to use just one of the backends.

## Round robin

Distributes requests across all backends starting from first, going to the last and then starting again with the first backend from the provided list.

### Why use it?

If you want to distribute amount of requests equally across all available backends.

## Sticky round robin

For each client it distributes requests across all backends starting from first, going to the last and then starting again with the first backend from the provided list.

If you have only one client, you will get the same experience as using [First](#first).

### Why use it?

> Other services might employ caches to keep a user’s state in RAM. This might be accomplished through hard or soft stickiness between reverse proxies and service frontends.

https://sre.google/sre-book/addressing-cascading-failures/#slow-startup-and-cold-caching

This is a good choice if a backend uses in-memory cache that is relevant for a specific client.

## Least-loaded round robin

Inspired by the [SRE Book](https://sre.google/sre-book/load-balancing-datacenter/#least-loaded-round-robin-WEswh9CN) I was reading at the time.

### Why use it?
