/******************
* © Danil_e71 2019
******************/


#include <errno.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <sys/types.h>
#include <sys/stat.h>
#include "main.h"

void readFile(FILE *fileptr,RotekHeader header) {
  RotekBlock block;

  char fname[128];

  fread(&block.size, sizeof(block.size), 1, fileptr);

  int len = __builtin_bswap32(block.size);

  fread(&block.header_size, sizeof(block.header_size), 1, fileptr);
  
  len-=__bswap_16(block.header_size) + sizeof(block.header_size);

  fread(&block.type, sizeof(block.type), 1, fileptr);

  printf("type = %d : (%s)\n" ,block.type, stringFromType(block.type));

  fread(&block.crc32, sizeof(block.crc32), 1, fileptr);

  short size = 0;

  block.signaturesize = (short *) malloc(( header.info.signatureCount) * sizeof(short));
  block.signature = malloc(( header.info.signatureCount) * sizeof(char*));

  for(int i = 0; i < header.info.signatureCount; i++) {
  	fread(&block.signaturesize[i], sizeof(block.signaturesize[i]), 1, fileptr);
  	size = __bswap_16(block.signaturesize[i]);
  	printf("signature[%d] length %d\n",i,size);
  	block.signature[i] = (char *) malloc((size) * sizeof(char));
  	fread(block.signature[i], size , 1, fileptr);
  }

  fread(&block.unused, sizeof(block.unused), 1, fileptr);
  fread(&block.sha1, sizeof(block.sha1), 1, fileptr);

  mkdir("unpacked", S_IRWXU);  

  switch(block.type) {
	case 2:
	  strcpy(fname,"unpacked/boot.img");
	break;
	case 3:
	  strcpy(fname,"unpacked/rootfs.img");
	break;
	case 6:
	  strcpy(fname,"unpacked/backup_boot.img");
	break;
  }

  printf("read %s len %d bytes\n",fname,len);

  FILE *save_fileptr = fopen(fname, "wb");

  int count = len/blockSize;
  unsigned char tmp[blockSize];
  for(int i = 0; i < count;i++) {
  fread(&tmp, blockSize, 1, fileptr); 
  fwrite(&tmp,blockSize,1,save_fileptr);
  }

  fclose(save_fileptr);
 }


void main (int argc, char **argv) {
  FILE *fileptr;
  RotekHeader header;

  if (argc <= 1) {
    printf("Usage: rotek-upgrade-package-extract [firmware-file-name]\n");
    exit (0);
  }

  fileptr = fopen(argv[1], "rb");

  if (fileptr == NULL) {
    printf("Not found: %s\n", argv[1]);
    exit (0);
  }

  fread(&header.info, sizeof(header.info), 1, fileptr);

  printf("vendor %s\n",header.info.vendor);
  printf("device %s hw %d\n",header.info.device,__builtin_bswap32(header.info.hwRev));
  printf("version %d.%d.%d\n",__bswap_16(header.info.v1),__bswap_16(header.info.v2),__bswap_16(header.info.v4));

  unsigned short size = 0;

  header.signaturesize = (short *) malloc(( header.info.signatureCount) * sizeof(short));
  header.signature = malloc(( header.info.signatureCount) * sizeof(char*));

  for(int i = 0; i < header.info.signatureCount; i++) {
  	fread(&header.signaturesize[i], sizeof(header.signaturesize[i]), 1, fileptr);
  	size = __bswap_16(header.signaturesize[i]);
  	printf("firmware signature[%d] length = %d\n",i,size);
  	header.signature[i] = (char *) malloc((size) * sizeof(char));
  	fread(header.signature[i], size , 1, fileptr);
  }

  fread(&header.unused, sizeof(header.unused), 1, fileptr);
  fread(&header.sha1, sizeof(header.sha1), 1, fileptr);

  readFile(fileptr,header); //boot
  readFile(fileptr,header); //rootfs
  readFile(fileptr,header); //backup_boot

  fclose(fileptr);
}
