fast: # Generic target with fast compression for development
	docker build --build-arg CONFIG_COMPRESSION=GZIP --build-arg CONFIG_PRIMARY=standard.json --build-arg CONFIG_MODULES=ALL -t registry2.lagoni.co.uk/os-builder:local .

small: # Generic target with small compression
	docker build --build-arg CONFIG_COMPRESSION=XZ --build-arg CONFIG_PRIMARY=standard.json --build-arg CONFIG_MODULES=ALL -t registry2.lagoni.co.uk/os-builder:local .
	
k300: # OnLogic K300 target
	docker build --build-arg CONFIG_COMPRESSION=XZ --build-arg CONFIG_PRIMARY=standard.json --build-arg CONFIG_MODULES=k300.txt -t registry2.lagoni.co.uk/os-builder:local .

magellis: # Schneider Magellis target
	docker build --build-arg CONFIG_COMPRESSION=XZ --build-arg CONFIG_PRIMARY=standard.json --build-arg CONFIG_MODULES=magelis.txt -t registry2.lagoni.co.uk/oos-builder:local .
	
run:
	docker container rm os-builder || true
	docker run -it --name os-builder -v $(PWD)/out:/out registry2.lagoni.co.uk/os-builder:local