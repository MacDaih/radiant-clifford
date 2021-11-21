all: build-clifford build-services

build-clifford:
	gcc -Wall -o thermo ~/embedded/clifford/clifford.c -lwiringPi
	chmod +x ./thermo 
	./thermo > /dev/null 2>&1

build-services:
	docker-compose up --build