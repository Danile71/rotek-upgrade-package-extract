all:  program

program: main.o
	gcc obj/main.o -o bin/rotek-upgrade-package-extract

main.o: src/main.c
	gcc -c src/main.c -o obj/main.o

clean:
	rm -rf bin obj

$(shell mkdir -p bin)
$(shell mkdir -p obj)
