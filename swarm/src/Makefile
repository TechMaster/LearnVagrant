deploy:
	docker network create -d overlay --attachable techmaster
	docker network create -d overlay --attachable techmaster_backend
	docker stack deploy -c dc_traefik.yml traefik
	docker stack deploy -c dc_portainer.yml portainer
	docker stack deploy -c dc_registry.yml registry
	cd websites; make build
	docker stack deploy -c dc_websites.yml websites