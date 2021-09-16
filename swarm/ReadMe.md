# Thực hành  Vagrant + Docker Swarm + Traefik + Portainer + Registry


### Thực hành các bước căn bản
Xem [Basic.md](Basic.md)

### Các bước triển khai nâng cao
Xem [Advanced.md](Advanced.md)

### Chạy luôn hệ thống đầy đủ
File [Vagrantfile](Vagrantfile) được bổ xung script shell để tự động hoá hầu hết cấu hình ban đầu ngoại trừ việc khởi tạo Docker Swarm.
```
$ vagrant up
$ vagrant ssh manager01
$ docker swarm init --listen-addr 192.168.33.2:2377 --advertise-addr 192.168.33.2:2377
$ docker swarm join-token manager
```

Mở terminal tab mới rồi ssh vào `manager02`, `manager03` rồi copy mã để join swarm

Quay lại terminal ssh vào manager01, chạy script tự động triển khai

```
$ cd src
$ make deploy
```

[src/Makefile](src/Makefile) sẽ gọi tiếp vào [src/websites/Makefile](src/websites/Makefile)

**Nếu mọi thứ chạy đúng thì sẽ có những đường dẫn như sau**
1. [http://dashboard.techmaster.com](http://dashboard.techmaster.com): Traefik dashboard. User/pass đăng nhập cuong/minh009-
2. [http://registry.techmaster.com](http://registry.techmaster.com): Danh sách docker image ở registry manager02:5000. User/pass đăng nhập cuong/minh009-
3. [http://potainer.techmaster.com](http://potainer.techmaster.com): Quản lý Docker bằng Portainer. User/pass tự tạo lúc đầu
4. [http://techmaster.com](http://techmaster.com): trang chủ
5. [http://techmaster.com/blog](http://techmaster.com/blog): blog ở trang chủ
6. [http://techmaster.com/admin](http://techmaster.com/admin): web site quản trị
7. [http://techmaster.com/teacher](http://techmaster.com/teacher): trang đành riêng cho giảng viên
8. [http://techmaster.com/user](http://techmaster.com/user): trang đành riêng cho sinh viên
9. [http://techmaster.com/video](http://techmaster.com/video): lưu trữ video
10. [http://techmaster.com/media](http://techmaster.com/media): lưu trữ media

### Xử lý sự cố

Hãy thử dùng lệnh sau đây
Liệt kê tất cả service
```
$ docker service ls
```

Xem một service cụ thể
```
$ docker service ps --no-trunc service_name
```

Xem các container đang chạy trên máy ảo đang ssh vào
```
$ docker ps
```

Thử vào  [http://dashboard.techmaster.com](http://dashboard.techmaster.com) để xem router và service có hiện lên trong danh sách hay không.

Thử vào [http://potainer.techmaster.com](http://potainer.techmaster.com)

hoặc bật logs ở `service.gateway.command '--log.level=DEBUG'` trong file [src/dc_traefik.yml](src/dc_traefik.yml). Khi vào xem logs của Traefik Gateway sẽ biết được routing có vấn đề gì không.

Việc định tuyến chập chờn hoặc chậm, hãy kiểm tra xem cấu hình DNS ở /etc/hosts máy chủ ngoài cùng, 3 máy ảo bên trong đã ok chưa.

Khi không push được image lên registry `manager02:5000` kiểm tra ở các máy ảo đã có file `/etc/docker/daemon.json` cho phép dùng HTTP chưa.