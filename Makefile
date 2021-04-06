build:

	docker build -t registry.lagoni.co.uk/os_build_env:local .
sata: 
	docker run -i -t --rm --privileged=true \
	-v build_data:/build/src \
	-v $(PWD)/out:/build/out \
	--name toolchain \
	registry.lagoni.co.uk/os_build_env:local /build.sh sata ALL

nvme: 
	docker run -i -t --rm --privileged=true \
	-v build_data:/build/src \
	-v $(PWD)/out:/build/out \
	--name toolchain \
	registry.lagoni.co.uk/os_build_env:local /build.sh nvme ALL