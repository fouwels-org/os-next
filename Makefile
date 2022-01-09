# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: Apache-2.0

IMAGE = ghcr.io/fouwels/os-next
TAG = local
COMPOSE=docker compose
BUILDFILE=compose.yml

default: fast

# Targets

DRPC-230: 	standard # IMI DRPC-230 target
k300: 		standard # OnLogic K300 target
k700: 		standard # OnLogic K700 target
magellis: 	standard # Schneider Magellis target

# Profiles

build: 
	$(COMPOSE) -f $(BUILDFILE) build

fast: 		# fast target for development/qemu
	$(COMPOSE) -f $(BUILDFILE) build --build-arg COMPRESSION_LEVEL=9 

fat: 		# fat target with all modules packed and available to be loaded
	$(COMPOSE) -f $(BUILDFILE) build --build-arg COMPRESSION_LEVEL=ALL 

# Output image
run:
	$(COMPOSE) -f $(BUILDFILE) up

clean:
	rm -rf ./out