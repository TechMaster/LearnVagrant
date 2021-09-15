docker build . -t iam:latest
docker run -d --name iam -p 8001:8001 iam:latest


docker image tag iam:latest manager02:5000/iam:latest
docker image push manager02:5000/iam:latest

Hiện đang gặp lỗi vì chưa hỗ trợ https
```
The push refers to repository [manager02:5000/iam]
Get https://manager02:5000/v2/: http: server gave HTTP response to HTTPS client
```

```
$ docker build . -t iam:0.1.0
$ docker image tag iam:0.1.1 manager02:5000/iam:0.1.1
$ docker image push manager02:5000/iam:0.1.1

$ curl -X GET https://manager02:5000/v2/_catalog
{"repositories":["iam"]}

$ curl -X GET http://manager02:5000/v2/iam/tags/list
{"name":"iam","tags":["latest","0.1.1"]}

docker stack deploy --compose-file docker-compose.iam.yml iam


docker service create \
  --name iam \
  --publish published=8001,target=8001 \
  --replicas 2 \
  manager02:5000/iam:0.1.1

  docker stack deploy --compose-file redis-ha.yml redis-ha