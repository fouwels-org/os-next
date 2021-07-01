# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: Apache-2.0

default: k300
	
k300: # OnLogic K300 target
	docker build \
	--build-arg CONFIG_PRIMARY=standard.yml \
	--build-arg CONFIG_MODULES=k300.txt \
	--build-arg KERNEL_CMDLINE="isolcpus=3 irqaffinity=0,1,2 rcu-nocbs=3 idle=poll console=ttyS0" \
	-t containers.fouwels.app/os-next:local .

magellis: # Schneider Magellis target
	docker build \
	--build-arg CONFIG_PRIMARY=standard.yml \
	--build-arg CONFIG_MODULES=magelis.txt \
	--build-arg KERNEL_CMDLINE="isolcpus=3 irqaffinity=0,1,2 rcu-nocbs=3 idle=poll console=ttyS0" \
	-t containers.fouwels.app/os-next:local .

all: # Generic fat target with all modules
	docker build \
	--build-arg CONFIG_PRIMARY=standard.yml \
	--build-arg CONFIG_MODULES=ALL \
	--build-arg KERNEL_CMDLINE="isolcpus=3 irqaffinity=0,1,2 rcu-nocbs=3 idle=poll console=ttyS0" \
	-t containers.fouwels.app/os-next:local .
	
run:
	docker container rm os-builder || true
	docker run -it --name os-builder -v $(PWD)/out:/out containers.fouwels.app/os-next:local