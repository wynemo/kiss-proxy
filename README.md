# kiss-proxy
this a simple and stupid http proxy server.
it supports [HTTP tunnel](https://en.wikipedia.org/wiki/HTTP_tunnel), and [HTTP keep-alive](https://en.wikipedia.org/wiki/HTTP_tunnel)


```bash
go get github.com/wynemo/kiss-proxy/httpproxy
$GOPATH/bin/httpproxy
export http_proxy=http://127.0.0.1:8118
```

and run this in python REPL:

```python
import requests
session = requests.session()
r = session.get('http://dabin.info/static/')
r = session.get('http://dabin.info/')
r = session.get('http://dabin.info/?page=2')
```
