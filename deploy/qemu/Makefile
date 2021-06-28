# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
#
# SPDX-License-Identifier: Apache-2.0

deploy: deploy-kernel

deploy-kernel: # Deploy and do not wipe disk
	./qemu.sh kernel persist ./base.img ../../out/BOOTx64.EFI
deploy-kernel-clean: # Deploy with wiped disk
	./qemu.sh kernel clean ./base.img ../../out/BOOTx64.EFI
deploy-disk:
	./qemu.sh disk persist ./base.img ../../out/BOOTx64.EFI
deploy-disk-clean:
	./qemu.sh disk clean ./base.img ../../out/BOOTx64.EFI

	
# Service commands
tpm:
	./qemu.sh startTPM