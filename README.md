# mkcert
自签名 nginx 证书
```
mkcert -ip=127.0.0.1 -domain=localhost
```

```
hellojukay@local mkcert (main) $ mkcert -h
mkcert [options]  generate TLS cert.
    -h
    --help        Print help message

    --domain      你的域名
    --ip          你的ip地址

    -v
    --verbose     Print debug log
```

使用已经存在的 CA 证书
```
./mkcert --domain=localhost --ip='127.0.0.1' --root-crt=ca.crt --root-key=ca.key
```

# 通过接口
```
go build
./mkcert-server
# open localhost:8080
```