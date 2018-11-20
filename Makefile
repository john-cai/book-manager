default: docker

docker:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bookmanager .
	docker build -t bookmanager .

docker-compose: docker
	docker-compose stop
	docker-compose rm -f
	docker-compose up -d
	
test:
	dropdb --if-exists bookmanager_test
	createdb bookmanager_test
	psql -U postgres -d bookmanager_test -a -f database/migration/*
	go test -v ./...