build:
	docker build -t registry.lagoni.co.uk/os_build_env:local .
run: 
	docker run -it --rm --privileged=true \
	-v build_data:/build \
	-v $(PWD)/out:/build/out \
	--name toolchain \
	registry.lagoni.co.uk/os_build_env:local /build.sh FACTORY