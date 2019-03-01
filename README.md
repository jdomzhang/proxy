# proxy
fiddler-like proxy

# install
> go get -v github.com/jdomzhang/proxy

# start the proxy
> proxy

```
Proxy listening on :7777 ...
```

# use it

## in Bash:

```
export HTTP_PROXY="http://localhost:7777"
```

## in Go code

```
os.Setenv("HTTP_PROXY", "http://localhost:7777")
```
