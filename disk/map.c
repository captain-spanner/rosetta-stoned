#include <u.h>
#include <libc.h>
#include <bio.h>
#include <mp.h>
#include <libsec.h>

enum {
	Nprimes	= 8,
	Hrot	= 5,
	HASHP	= 16777619,
};

uchar	*map;
uint	mapz;
double	phi	= 1.618033988749894848204586834;

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

uint32
primege(uint32 n) {
	mpint *m;
	if((n & 1) == 0)
		n++;
	for(;;) {
		m = mpnew(0);
		m = uitomp(n,m);
		if(smallprimetest(m)) {
			mpfree(m);
			return n;
		}
		mpfree(m);
		n += 2;
	}
}

void
mapsetup(void) {
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
	mapz = primege((uint32)(2 * mapz * phi));
	map = malloc(mapz);
}

uint
rot(uint32 h) {
	return (h << Hrot) | (h >> (32 - Hrot));
}

void
setmap(uint x, int i) {
	if(x >= mapz) {
		fprint(2, "wtf: x = %ud, mapsz = $ud\n", x, mapz);
		exits("map: setmap");
	}
	map[x] |= 1 << i;
}

void
mapfiles(void) {
	uint h;
	char *s;
	Biobuf *b;
	int i, l, x;

	b = Bopen("«files»", O_RDONLY);
	while((s = Brdstr(b, '\n', 1)) != nil) {
		l = Blinelen(b);
		if(s[l-1] == '\n')
			s[l-1] = '\0';
		h = hash(s);
		for(i = 0; i < Nprimes; i++) {
			h *= primes[i];
			x = h % mapz;
			setmap(x, i);
			h = rot(h);
		}
	}
	Bterm(b);
}

void
writemap(void) {
	int n, fd;

	fd = create("«map»", O_WRONLY, 0664);
	n = write(fd, map, mapz);
	if(n != mapz) {
		fprint(2, "wrote %d to map, return %d\n", mapz, n);
		exits("map: writemap");
	}
	close(fd);
}
		
void
main(void) {
	quotefmtinstall();
	mapsetup();
	mapfiles();
	writemap();
	exits(nil);
}
