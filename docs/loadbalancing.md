# Load balancing

## Random

Chooses some backend from the provided list at random.

### Why use it?

If you want to introduce some chaos.

## First

Chooses the first backend from the provided list.

### Why use it?

If you explicitly want to use just one of the backends.

## Round robin

## Sticky round robin

### Why use it?

> Other services might employ caches to keep a user’s state in RAM. This might be accomplished through hard or soft stickiness between reverse proxies and service frontends.

https://sre.google/sre-book/addressing-cascading-failures/#slow-startup-and-cold-caching

This is a good choice if a backend uses in-memory cache that is relevant for a specific client.

## Least-loaded round robin

Inspired by the [SRE Book](https://sre.google/sre-book/load-balancing-datacenter/#least-loaded-round-robin-WEswh9CN) I was reading at the time.
