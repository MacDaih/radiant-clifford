UID := $(shell id -u)
GID := $(shell id -g)

all: build-clifford build-services

build-clifford:
	gcc -Wall -o thermo ./clifford/clifford.c -lwiringPi
	chmod +x ./thermo
	./thermo & 