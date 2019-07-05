/******************
* © Danil_e71 2019
******************/




#define blockSize 128

#pragma pack(push, 4)

typedef struct {
	unsigned int size;
	unsigned short header_size;
    unsigned short type;
	unsigned int crc32;
	unsigned short *signaturesize;
	unsigned char **signature;
    unsigned short unused;
    unsigned char sha1[20];
}RotekBlock;

typedef struct {
	unsigned char vendor[32];
	unsigned char device[32];
	unsigned short v1;
	unsigned short v2;
	unsigned short v3;
	unsigned short v4;
	unsigned char unk[12];
	unsigned int hwRev;
	unsigned char signatureCount;
	unsigned char unk1[3];
}RotekInfoHeader;

typedef struct {
	RotekInfoHeader info;
	unsigned short *signaturesize;
	unsigned char **signature;
	unsigned short unused;
	unsigned char sha1[20];
}RotekHeader;


#pragma pack(pop)

enum TYPE {unk0,unk1,kernel,rootfs,unk4,unk5,backup_kernel};

inline char *stringFromType(enum TYPE f)
{
 char *strings[] = { "unk0", "unk1", "Linux Kernel Image", "Root FS Image","unk4","Branding Image","Backup Linux Kernel Image"};
    return strings[f];
}

char *stringFromType(enum TYPE f);


