Tham khảo [Learn Docker Swarm with Vagrant](https://levelup.gitconnected.com/learn-docker-swarm-with-vagrant-47dd52b57bcc)
## Khởi tạo 3 máy ảo
```
$ vagrant up
```

## Thực hành trên `manager01`
1. SSH vào manager01
    ```
    $ vagrant ssh manager01
    ```

2. Hãy kiểm tra ip của nó có đúng là `192.168.33.2` không
    ```
    $ ifconfig
    eth1: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
      inet 192.168.33.2  netmask 255.255.255.0  broadcast 192.168.33.255
      inet6 fe80::20c:29ff:fe4e:3e01  prefixlen 64  scopeid 0x20<link>
      ether 00:0c:29:4e:3e:01  txqueuelen 1000  (Ethernet)
      RX packets 20362  bytes 4847911 (4.8 MB)
      RX errors 0  dropped 0  overruns 0  frame 0
      TX packets 19272  bytes 4830880 (4.8 MB)
      TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
    ```
3. Khởi tạo Docker Swarm
    ```
    $ docker swarm init --listen-addr 192.168.33.2:2377 --advertise-addr 192.168.33.2:2377

    To add a worker to this swarm, run the following command:
    docker swarm join --token SWMTKN-1-1t4mi73lltrxg3ipxxoxh1v0ofe6knlmis81qezkngp62gyzmk-8zeuc96e7ah5lh8f8kbfi0nzn 192.168.33.2:2377

    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
    ```

5. Gõ lệnh dưới đấy để lấy token tạo manager tiếp theo
    ```
    $ docker swarm join-token manager

    To add a manager to this swarm, run the following command:

    docker swarm join --token SWMTKN-1-1t4mi73lltrxg3ipxxoxh1v0ofe6knlmis81qezkngp62gyzmk-7a0jkvzldoigbtyl1vezzh0xp 192.168.33.2:2377
    ```

## Thực hành trên manager02

1. SSH vào manager02
   ```
   $ vagrant ssh manager02
   ```
2. Gõ lệnh để join docker swarm dưới role manager
   ```
   docker swarm join --token SWMTKN-1-1t4mi73lltrxg3ipxxoxh1v0ofe6knlmis81qezkngp62gyzmk-7a0jkvzldoigbtyl1vezzh0xp 192.168.33.2:2377
   ```

## Thực hành trên worker01
1. SSH vào worker01
   ```
   $ vagrant ssh worker01
   ```
2. Gõ lệnh để join docker swarm dưới role worker
   ```
   docker swarm join --token SWMTKN-1-1t4mi73lltrxg3ipxxoxh1v0ofe6knlmis81qezkngp62gyzmk-8zeuc96e7ah5lh8f8kbfi0nzn 192.168.33.2:2377
   ```


## Liệt kê các node trong docker swarm
```
$ docker node ls
ID                            HOSTNAME    STATUS    AVAILABILITY   MANAGER STATUS   ENGINE VERSION
6dijd48m3lcmk99835342n7om *   manager01   Ready     Active         Leader           20.10.7
76cyxv9mp1wmpa2bhn97mkn6k     manager02   Ready     Active         Reachable        20.10.7
ihiwerjs86p7uv3czx7x91obl     worker01    Ready     Active                          20.10.7
```

## Quay lại manager01 cài Portainer

1. SSH vào manager01
    ```
    $ vagrant ssh manager01
    ```
2. Gõ lệnh này sau đây để lấy file `portainer-agent-stack.yml`
   ```
   $ curl -L https://downloads.portainer.io/portainer-agent-stack.yml -o portainer-agent-stack.yml
   ```
3. Gõ lệnh này tiếp theo để triển khai stack
   ```
   $ docker stack deploy -c portainer-agent-stack.yml portainer
   ```

## Điều khiển portainer qua http://localhost:9000

Do file [Vagrantfile](Vagrantfile) cấu hình forwarded_port '9000'
```ruby
if i == 1  # Server Manager đầu tiên của Docker Swarm
  #Only configure port to host for Manager01
  manager.vm.network :forwarded_port, guest: 8080, host: 8080
  manager.vm.network :forwarded_port, guest: 5000, host: 5000
  manager.vm.network :forwarded_port, guest: 9000, host: 9000
end
```

Giờ bạn chỉ cần vào `http://localhost:9000` là thấy được giao diện của Portainer

## Thực hành xong thì tạm thời dừng

Gõ lệnh
```
$ vagrant halt
```

Khi nào cần dùng lại
```
$ vagrant up
```

## Deploy Nginx lên Portainer

```yaml
version: "3.8"

services:
  web:
    image: nginx:alpine
    ports:
    - "8001:80"
```

![](img/ngix.jpg)

Xem Swarm Cluster
![](img/swarm_cluster.jpg)

Mở trình duyệt tới địa chỉ http://localhost:8001

## Deploy WhoAmI

```yaml
version: "3.8"

services:
   whoami:
    image: containous/whoami:dev
    ports:
    - "5000:80"
```

## Deploy Redis
```yaml
version: "3.8"
services:
 redis:
    image: "redis:alpine"
    command: redis-server --requirepass abc123
    ports:
     - "6379:6379"
    environment:
     - REDIS_REPLICATION_MODE=master
```
