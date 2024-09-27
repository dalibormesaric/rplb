# Load-balancing

## Random

## First

Chooses the first backend from the provided list.

## Round-robin

## Sticky round-robin

This is a good choice if a backend uses in-memory cache that is relevant for a specific client.

## Least-loaded round-robin

Inspired by the [SRE Book](https://sre.google/sre-book/load-balancing-datacenter/#least-loaded-round-robin-WEswh9CN) I was reading at the time.
