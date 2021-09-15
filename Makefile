# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: Apache-2.0

IMAGE = containers.fouwels.app/fouwels/os-next
TAG = local

standard: 	# Standard build configuration
	docker build -t $(IMAGE):$(TAG) .

turbo: 		# fast target for development/qemu
	docker build --build-arg COMPRESSION_LEVEL=9 -t $(IMAGE):$(TAG) .

fat: 		# fat target with all modules packed and available to be loaded
	docker build --build-arg CONFIG_MODULES=ALL -t $(IMAGE):$(TAG) .

DRPC-230: 	standard # IMI DRPC-230 target
k300: 		standard # OnLogic K300 target
k700: 		standard # OnLogic K700 target
magellis: 	standard # Schneider Magellis target
	
run:
	docker container rm os-builder || true
	docker run -it --name os-builder -v $(PWD)/out:/out $(IMAGE):$(TAG)