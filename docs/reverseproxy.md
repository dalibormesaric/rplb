# Reverse Proxy

Describe how Router has to point to Home Assistant as DNS Server.

Describe how it works together with Home Assistant + AdGuard Home setup

DNS rewrites - Allows to easily configure custom DNS response for a specific domain name.

## How it works?

Reverse proxy extracts host name from the incoming request and tries to find frontend with the same name. On success, reverse proxy finds a backend connected to that frontend.