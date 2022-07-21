# rdp-socks5
Remote Desktop over Socks5 proxy

# usage
## Run socks5 proxy
![image](https://user-images.githubusercontent.com/8877695/180183968-3fbc99f2-9a5f-4cbe-9d68-9a189c06969a.png)

## Run proxy rdp
```
rdp-socks5.exe -l 127.0.0.1:3388 -r 10.0.9.11:3389 -x "socks5://127.0.0.1:1081?timeout=5m"
```

## Using rdp over proxy
![image](https://user-images.githubusercontent.com/8877695/180182226-711c0833-57c4-4e4f-9102-778f826e79fa.png)
