# Load balancing

Here are short descriptions of implemented load balancing algorithms and why to use them.

## Random

Chooses at random some backend from the provided list.

### Why use it?

If you want to introduce some chaos.

![random](./random.gif)

## First

Chooses the first backend from the provided list.

### Why use it?

If you explicitly want to use just one of the backends.

![first](./first.gif)

## Round robin

Distributes requests across all backends starting from first, going to the last and then starting again with the first backend from the provided list.

### Why use it?

If you want to distribute amount of requests equally across all available backends.

![roundrobin](./roundrobin.gif)

## Sticky round robin

For each client it distributes requests across all backends starting from first, going to the last and then starting again with the first backend from the provided list.

If you have only one client, you will get the same experience as using [First](#first).

### Why use it?

> Other services might employ caches to keep a user’s state in RAM. This might be accomplished through hard or soft stickiness between reverse proxies and service frontends.

https://sre.google/sre-book/addressing-cascading-failures/#slow-startup-and-cold-caching

This is a good choice if a backend uses in-memory cache that is relevant for a specific client.

## Least-loaded round robin

Inspired by [this section](https://sre.google/sre-book/load-balancing-datacenter/#least-loaded-round-robin-WEswh9CN) in the SRE book I was reading at the time.

### Why use it?

This is an alternative to round robin when some requests take longer time to resolve. This prevents the scenario when all heavy requests would hit the same backend. Instead the load is spread out across all other available backends.

### Screenshot

First three backends have 200ms delay, second three have 100ms and last three have no delay.

![leastloaded](./leastloaded.gif)
