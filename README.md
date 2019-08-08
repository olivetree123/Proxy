# TCP Proxy

TCP Proxy 示例，基于 [go-tcp-proxy](https://github.com/jpillora/go-tcp-proxy)，做了一些简化。


## Usage

```
$ tcp-proxy --help
Usage of tcp-proxy:
  -l="localhost:9999": local address
  -n: disable nagles algorithm
  -r="localhost:80": remote address
```

*Note: Regex match and replace*
**only works on text strings**
*and does NOT work across packet boundaries*

### Simple Example

Since HTTP runs over TCP, we can also use `tcp-proxy` as a primitive HTTP proxy:

```
$ tcp-proxy -r echo.jpillora.com:80
Proxying from localhost:9999 to echo.jpillora.com:80
```

Then test with `curl`:

```
$ curl -H 'Host: echo.jpillora.com' localhost:9999/foo
{
  "method": "GET",
  "url": "/foo"
  ...
}
```

### Match Example

```
$ tcp-proxy -r echo.jpillora.com:80 -match 'Host: (.+)'
Proxying from localhost:9999 to echo.jpillora.com:80
Matching Host: (.+)

#run curl again...

Connection #001 Match #1: Host: echo.jpillora.com
```

### Replace Example

```
$ tcp-proxy -r echo.jpillora.com:80 -replace '"ip": "([^"]+)"~"ip": "REDACTED"'
Proxying from localhost:9999 to echo.jpillora.com:80
Replacing "ip": "([^"]+)" with "ip": "REDACTED"
```

```
#run curl again...
{
  "ip": "REDACTED",
  ...
```

*Note: The `-replace` option is in the form `regex~replacer`. Where `replacer` may contain `$N` to substitute in group `N`.*

### Todo

* Implement `tcpproxy.Conn` which provides accounting and hooks into the underlying `net.Conn`
* Verify wire protocols by providing `encoding.BinaryUnmarshaler` to a `tcpproxy.Conn`
* Modify wire protocols by also providing a map function
* Implement [SOCKS v5](https://www.ietf.org/rfc/rfc1928.txt) to allow for user-decided remote addresses
