COMPOSE=docker-compose
BUILDFILE=build.yml
DOCKER=docker

.PHONY: build push up down up-d

#Docker
build: Dockerfile
	$(COMPOSE) -f $(BUILDFILE) build
push: Dockerfile
	$(COMPOSE) -f $(BUILDFILE) push
up: Dockerfile
	$(COMPOSE) -f $(BUILDFILE) up
up-d: Dockerfile
	$(COMPOSE) -f $(BUILDFILE) up -d
down: Dockerfile
	$(COMPOSE) -f $(BUILDFILE) down
down-v: Dockerfile
	$(COMPOSE) -f $(BUILDFILE) down -v


