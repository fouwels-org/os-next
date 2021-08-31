# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: Apache-2.0

default: 	standard
DRPC-230: 	standard # IMI DRPC-230 target
k300: 		standard # OnLogic K300 target
k700: 		standard # OnLogic K700 target
magellis: 	standard # Schneider Magellis target

standard: # Standard build configuration
	docker build \
	--build-arg CONFIG_PRIMARY=standard.yml \
	--build-arg CONFIG_MODULES=standard.mod \
	-t containers.fouwels.app/fouwels/os-next:local .

turbo: # fast target for development/qemu
	docker build \
	--build-arg COMPRESSION_LEVEL=9 \
	--build-arg CONFIG_PRIMARY=standard.yml \
	--build-arg CONFIG_MODULES=standard.mod \
	-t containers.fouwels.app/fouwels/os-next:local .

all: # Generic fat target with all modules available to be loaded
	docker build \
	--build-arg CONFIG_PRIMARY=standard.yml \
	--build-arg CONFIG_MODULES=ALL \
	-t containers.fouwels.app/fouwels/os-next:local .
	
run:
	docker container rm os-builder || true
	docker run -it --name os-builder -v $(PWD)/out:/out containers.fouwels.app/fouwels/os-next:local