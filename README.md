This repository is a reproduction for [Chromium issue 1380505](https://bugs.chromium.org/p/chromium/issues/detail?id=1380505).
The issue is that it is not possible to catch both HTML and JSON responses for the same URL. Requesting both resources
invalidate any cached version of each other.

To run, you should use [Go](https://go.dev/):

```
go run main.go
```

And then open [http://localhost:8000/](http://localhost:8000/) in Chromium. Open network tab in DevTools.

## Expected

The first time you should see three requests made, one for `/` for `text/html`, another for `/data.json`
for `application/json` and another for `/data.json` for `text/html`. In the server's terminal you should see
something like:

```
2023/07/20 11:51:04 Listening on :8000
2023/07/20 11:51:12 Serving HTML at /
2023/07/20 11:51:12 Serving JSON at /data.json
2023/07/20 11:51:12 Serving HTML at /data.json
```

Loading the page again should ideally not make any requests (handled all from cache) given that all responses
are set to `immutable` `Cache-Control` response header, or at least there should be only request for `/`
(as it is the main requested page) made.

## Actual

You see every time you open the page three requests made to the server.

## Discussion

Etags do not help. `immutable` `Cache-Control` header does not help. A browser still make requests.

Tested on Chromium 114.0.5735.198. Firefox 115.0.2 behaves the same.
