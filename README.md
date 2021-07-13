Callbacks service
=================

Why:
---

Sometimes you need to publicly expose one endpoint so it can be the entrypoint for a service callback.
One example is Oauth2 flow, after the consent the service redirects to a pre-configured url.
Image having many environments (i.e. ephemeral dev environments) instead of registering a new callback url for your Oauth2 integration
is better to have a common entrypoint for all, and then based on the state query param do a final redirection to the real callback.

How:
---

This service is stateless and only have one endpoint, `/redirect` which is the one that should be specified as the callback url.
The idea is to encode the real redirection in the state param as a base64 encoded string, i.e.:

```bash

# First we build the application

$ make build

# And run it, use `-port` and `-host` to define host and port

$ ./app.exe -port 8080 -host 0.0.0.0

# in another terminal:

$ echo -n https://google.com | base64

# this returns: aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbQ==

# And now act as the service doing the redirection

$ curl -i http://localhost:8080/redirect\?state\=aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbQ==


HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://www.google.com
Date: Tue, 13 Jul 2021 21:00:59 GMT
Content-Length: 45

<a href="https://www.google.com">Found</a>.

```

- You can also run this as a docker container: `docker run -p 8080:8080 dmanchon/callbacks:latest`

- Implements a `/metrics` endpoint for prometheus to scrape and `/health` endpoint for liveness test.
