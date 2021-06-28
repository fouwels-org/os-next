# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
# SPDX-FileCopyrightText: 2021 K. Fouwels <k@fouwels.com>
#
# SPDX-License-Identifier: Apache-2.0

all: # Generic target with fast compression for development
	docker build --build-arg CONFIG_COMPRESSION=GZIP --build-arg CONFIG_PRIMARY=standard.yml --build-arg CONFIG_MODULES=ALL -t containers.fouwels.app/os-next:local .

all-small: # Generic target with small compression
	docker build --build-arg CONFIG_COMPRESSION=XZ --build-arg CONFIG_PRIMARY=standard.yml --build-arg CONFIG_MODULES=ALL -t containers.fouwels.app/os-next:local .
	
k300: # OnLogic K300 target
	docker build --build-arg CONFIG_COMPRESSION=GZIP --build-arg CONFIG_PRIMARY=standard.yml --build-arg CONFIG_MODULES=k300.txt -t containers.fouwels.app/os-next:local .

magellis: # Schneider Magellis target
	docker build --build-arg CONFIG_COMPRESSION=GZIP --build-arg CONFIG_PRIMARY=standard.yml --build-arg CONFIG_MODULES=magelis.txt -t containers.fouwels.app/os-next:local .
	
run:
	docker container rm os-builder || true
	docker run -it --name os-builder -v $(PWD)/out:/out containers.fouwels.app/os-next:local