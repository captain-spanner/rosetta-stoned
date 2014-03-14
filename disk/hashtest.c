#include <u.h>
#include <libc.h>
#include <bio.h>

enum {
	HASHP	= 16777619,
};

uint
hash(char *s) {
	uint c, h;

	for(h = (uint)*s++; (c = (uint)*s) != '\0'; s++) {
		h *= HASHP;
		h ^= c;
	}

	return h;
}

void
main(int argc, char *argv[]) {
	int i;

	quotefmtinstall();
	for(i = 1; i < argc; i++) {
		print("%q: %ud\n", argv[i], hash(argv[i]));
	}
	exits(nil);
}
