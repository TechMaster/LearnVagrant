build:
	cd main; docker build . -t manager02:5000/main:latest; docker push manager02:5000/main:latest
	cd admin; docker build . -t manager02:5000/admin:latest; docker push manager02:5000/admin:latest
	cd teacher; docker build . -t manager02:5000/teacher:latest; docker push manager02:5000/teacher:latest
	cd user; docker build . -t manager02:5000/user:latest; docker push manager02:5000/user:latest
	cd video; docker build . -t manager02:5000/video:latest; docker push manager02:5000/video:latest
	cd media; docker build . -t manager02:5000/media:latest; docker push manager02:5000/media:latest

push:
	docker push manager02:5000/main:latest
	docker push manager02:5000/admin:latest
	docker push manager02:5000/teacher:latest
	docker push manager02:5000/user:latest
	docker push manager02:5000/video:latest
	docker push manager02:5000/media:latest
