# Triển khai Redis High Availability trong Docker Swarm

## Replicate  2 Redis
Creditss: Bùi Hiên và Phạm Hùng

Triển khai redis cluster trên 2 node:
1. Node master ở máy ảo manager01
2. Node slave ở máy áo manager02

Đầu tiên cần tạo folder trên để mount data redis trên 2 máy áo này:
```
$ vagrant ssh manager01
```
Trong máy ảo `manager01` tạo thư mục
```
$ mkdir /home/vagrant/redis
```

Vào `manager02`
```
$ vagrant ssh manager02
```
Trong máy ảo `manager02` tạo thư mục
```
$ mkdir /home/vagrant/redis
```

Trên ứng dụng Portainer vào phần Stack > Add Stack triển khai file Docker-compose.yml như dưới

```yaml
version: "3.8"
networks:
  my-network:
    external: false
services:
  redis-master: #tên service
    image: redis:alpine #tên image
    command: redis-server --requirepass 123 # set password cho redis
    ports:
      - "6379:6379" # ánh xạ cổng 6379 của container ra ngoài cổng 6379 trên máy host
    volumes:
      - /home/vagrant/redis:/data # mount volume từ thư mục /data của container ra ngoài thư mục /home/vagrant/redis trên máy host
    networks:
      - my-network # network overlay để các container trong Network này có thể giao tiếp được với nhau
    deploy:
      placement:
        constraints: # Chỉ định node quản lý
          - node.role == manager
          - node.hostname == manager01
  redis-slave-1: #tên service
    image: redis:alpine #tên image
    command: redis-server --masterauth 123 --slaveof redis-master 6379
    depends_on:
      - redis-master
    ports:
      - "6380:6379" # ánh xạ cổng 6379 của container ra ngoài cổng 6380 trên máy host
    volumes:
      - /home/vagrant/redis:/data # mount volume từ thư mục /data của container ra ngoài thư mục /home/vagrant/redis trên máy host
    networks:
      - my-network # Network overlay để các container trong Network này có thể giao tiếp được với nhau
    deploy:
      placement:
        constraints: # Chỉ định node quản lý
          - node.role == manager
          - node.hostname == manager02
```

Cải tiến

```yaml
version: "3.8"
services:
  redis-master: #tên service
    image: redis:alpine #tên image
    command: redis-server --requirepass 123 # set password cho redis
    ports:
      - "6379:6379" # ánh xạ cổng 6379 của container ra ngoài cổng 6379 trên máy host
    volumes:
      - /home/vagrant/redis:/data # mount volume từ thư mục /data của container ra ngoài thư mục /home/vagrant/redis trên máy host    
    deploy:
      placement:
        constraints: # Chỉ định node quản lý
          - node.role == manager
          - node.hostname == manager01
  redis-slave-1: #tên service
    image: redis:alpine #tên image
    command: redis-server --masterauth 123 --slaveof redis-master 6379
    depends_on:
      - redis-master
    ports:
      - "6380:6379" # ánh xạ cổng 6379 của container ra ngoài cổng 6380 trên máy host
    volumes:
      - /home/vagrant/redis:/data # mount volume từ thư mục /data của container ra ngoài thư mục /home/vagrant/redis trên máy host    
    deploy:
      placement:
        constraints: # Chỉ định node quản lý
          - node.role == manager
          - node.hostname == manager02
  redis-slave-2: #tên service
    image: redis:alpine #tên image
    command: redis-server --masterauth 123 --slaveof redis-master 6379
    depends_on:
      - redis-master
    ports:
      - "6381:6379" # ánh xạ cổng 6379 của container ra ngoài cổng 6380 trên máy host
    volumes:
      - /home/vagrant/redis:/data # mount volume từ thư mục /data của container ra ngoài thư mục /home/vagrant/redis trên máy host    
    deploy:
      placement:
        constraints: # Chỉ định node quản lý
          - node.hostname == worker01
```


#### Cài đặt redis-cli trên manager01
```
$ sudo -i
# add-apt-repository ppa:redislabs/redis
# apt update
# apt install redis-tools
```

#### Tạo key ở master Redis
Kết nối vào redis ở `manager01:6379`
```
$ redis-cli -h manager01 -p 6379 -a 123
```

Tạo key `foo` có giá trị `Bar`
```
manager01:6379> SET foo "Bar"
```

#### Tạo SSH session khác ở `manager01` để đọc key từ slave Redis
Kết nối vào redis ở `manager02:6380`
```
$ redis-cli -h manager02 -p 6380 -a 123
```

Lấy giá trị `foo`
```
manager02:6380> GET foo
"Bar"
```

Như vậy khi set key ở Master Redis, chúng ta đọc được key ở Slave Redis

> Mặc dùng chúng ta cấu hình thành công Redis Master - Slave configuration. Nhưng thực tế ứng dụng luôn chỉ kết nối vào Redis Master Node bởi trong connection chúng ta ghi rõ port `6379`. Mô hình này không hỗ trợ High Availability - Fail Over

## Docker Swarm + Redis Sentinel

Redis Cluster dùng để phân phối, chia dữ liệu (partioning) ra các node khác nhau theo địa lý hoặc theo năm.

Còn Redis + Sentinel tập trung vào mục đích High Availability: Failed Over

Tham khảo:
- [Redis — High Availability with Docker Swarm](https://medium.com/@emmano3h/redis-high-availability-with-docker-swarm-2142a4d80b49)
- [Redis Sentinel — High Availability: Everything you need to know from DEV to PROD: Complete Guide](https://medium.com/@amila922/redis-sentinel-high-availability-everything-you-need-to-know-from-dev-to-prod-complete-guide-deb198e70ea6)
- [bitnami/redis-sentinel](https://hub.docker.com/r/bitnami/redis-sentinel/)

```yaml
version: '3.8'

networks:
  app-tier:
    driver: bridge

services:
  redis:
    image: 'bitnami/redis:latest'
    environment:
      - REDIS_REPLICATION_MODE=master
      - REDIS_PASSWORD=123
    networks:
      - app-tier
    ports:
      - '6379'
  redis-slave:
    image: 'bitnami/redis:latest'
    environment:
      - REDIS_REPLICATION_MODE=slave
      - REDIS_MASTER_HOST=redis
      - REDIS_MASTER_PASSWORD=123
      - REDIS_PASSWORD=123
    ports:
      - '6379'
    depends_on:
      - redis
    networks:
      - app-tier
  redis-sentinel:
    image: 'bitnami/redis-sentinel:latest'
    environment:
      - REDIS_MASTER_PASSWORD=123
    depends_on:
      - redis
      - redis-slave
    ports:
      - '26379-26381:26379'
    networks:
      - app-tier
```

## Replicate Postgresql


## Triển khai Traefik

## Cấu hình Docker Secret