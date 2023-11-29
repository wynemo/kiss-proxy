# kiss-proxy
this a simple and stupid http proxy server.
it supports [HTTP tunnel](https://en.wikipedia.org/wiki/HTTP_tunnel), and [HTTP keep-alive](https://en.wikipedia.org/wiki/HTTP_tunnel)


```bash
go get github.com/wynemo/kiss-proxy/httpproxy
$GOPATH/bin/httpproxy 0.0.0.0:7000&
export http_proxy=http://127.0.0.1:7000
export https_proxy=http://127.0.0.1:7000
```

and run this in python REPL:

```python
import requests
session = requests.session()
r = session.get('http://baidu.com')
r = session.get('https://douban.com')
```


TODO: handle http 100 status


## compile locally

`go build -o output/httpproxy httpproxy/httpproxy.go`
