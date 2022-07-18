UID := $(shell id -u)
GID := $(shell id -g)

all: build-clifford build-services

build-clifford:
	gcc -Wall -o thermo ./clifford/clifford.c -lwiringPi

run-clifford:
	modprobe i2c-dev
	chmod +x ./thermo
	./thermo & 