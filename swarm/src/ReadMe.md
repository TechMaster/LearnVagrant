
# Triển khai hệ thống gồm Traefik và một vài dịch vụ

1. Traefik Gateway
2. traefik/whoami
3. iam

## Build docker image iam
```
$ docker build . -t iam:latest
$ docker image tag iam:latest manager02:5000/iam:latest
$ docker image push manager02:5000/iam:latest
```

### Nếu gặp lỗi khi push
```
The push refers to repository [manager02:5000/iam]
Get https://manager02:5000/v2/: http: server gave HTTP response to HTTPS client
```

## Cấu hình cho phép truy cập insecured Docker Registry
Docker Registry server mặc định yêu cầu HTTPs để phục vụ. Trong môi trường thử nghiệm Vagrant trên local, không bật được HTTPS thì chúng ta cấu hình kết nối vào insecured docker registry.

Kết nối SSH vào manager01
```
$ vagrant ssh manager01
```

Chuyển sang root rồi tạo file `/etc/docker/daemon.json`
```
$ sudo -i
$ nano /etc/docker/daemon.json
```

Thêm nội dung như sau
```json
{
    "insecure-registries" : [ "manager02:5000" ]
}
```

## Kiểm tra lại docker registry
$ curl -X GET https://manager02:5000/v2/_catalog
{"repositories":["iam"]}

$ curl -X GET http://manager02:5000/v2/iam/tags/list
{"name":"iam","tags":["latest"]}

## Deploy cả stack

```
$ docker stack deploy -c docker-compose.yml traefik
```

## Thử nghiệm

Vào địa chỉ
- http://localhost
- http://iris.com
- http://localhost:8080