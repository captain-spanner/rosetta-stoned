#include <u.h>
#include <libc.h>
#include <bio.h>

enum {
	Nprimes	= 8,
	Hrot	= 5,
	HASHP	= 16777619,
};

uint	mapz;

// 8 random primes. could use one without degredation?
int	primes[Nprimes] = {
	73,
	167,
	379,
	661,
	1129,
	2039,
	3797,
	6563,
};

int
hash(char *s) {
	uint c, h;

	for(h = (uint)*s++; (c = (uint)*s) != '\0'; s++) {
		h *= HASHP;
		h ^= c;
	}

	return h;
}

uint
rot(uint32 h) {
	return (h << Hrot) | (h >> (32 - Hrot));
}

void
map_setup(void) {
	char *s;
	Biobuf *b;

	b = Bopen("«count»", O_RDONLY);
	s = Brdstr(b, '\n', 1);
	Bterm(b);
	mapz = atol(s);
	if(mapz == 0 ){
		fprint(2, "map: bad «count» %s\n", s);
		exits("map: count");
	}
}

void
map(char *s) {
	int h, i, x;

	print("%s:\n", s);
	h = hash(s);
	for(i = 0; i < Nprimes; i++) {
		h *= primes[i];
		x = h % mapz;
		print("setmap(%d, %d)\n", x, i);
		h = rot(h);
	}
}
		
void
main(int argc, char *argv[]) {
	int i;

	quotefmtinstall();
	map_setup();
	for(i = 1; i < argc; i++) {
		map(argv[i]);
	}
	exits(nil);
}
