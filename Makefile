all: rotek-upgrade-package-extract

rotek-upgrade-package-extract: main.o
	gcc main.o -o rotek-upgrade-package-extract

main.o: main.c
	gcc -c main.c


clean:
	rm -rf *.o rotek-upgrade-package-extract
