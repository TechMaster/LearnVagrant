deploy:
	docker network create -d overlay --attachable techmaster
	docker network create -d overlay --attachable techmaster_backend
	docker stack deploy -c dc_traefik.yml traefik
	docker stack deploy -c dc_registry.yml registry
	docker stack deploy -c dc_websites.yml websites