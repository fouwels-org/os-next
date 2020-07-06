COMPOSE=docker-compose
BUILDFILE=build.yml
DOCKER=docker

.PHONY: build push up down up-d sftp

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

sftp:
	sftp root@10.0.10.203:/boot/efi/EFI/ <<< $'put ./out/uinit-custom'


