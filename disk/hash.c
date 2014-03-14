#include <u.h>
#include <libc.h>
#include <bio.h>

enum {
	Hlimit	= 253,
	Nprimes	= 4,
	HASHP	= 16777619,
};

char	pre[]	= "Hash-";
uint	hashsz	= Hlimit;
int	primes[Nprimes] = {
	31,
	61,
	127,
	253,
};

extern int stat(char *, uchar *, int);

int
hash(char *s) {
	uint c, h;

	for(h = (uint)*s++; (c = (uint)*s) != '\0'; s++) {
		h *= HASHP;
		h ^= c;
	}

	return h % hashsz;
}

void
mkdirs() {
	int i;

	for(i = 0; i < hashsz; i++) {
		print("mkdir -p %s%02X\n", pre, i);
		print("echo %02X >> «dirs»\n", i);
	}
}

void
hashfiles(void) {
	int l;
	char *s;
	Biobuf *b;

	b = Bopen("«files»", O_RDONLY);
	while((s = Brdstr(b, '\n', 1)) != nil) {
		l = Blinelen(b);
		if(s[l-1] == '\n')
			s[l-1] = '\0';
		fprint(2, "%#q -> %s%02X\n", s, pre, hash(s));
		print("mv -- %#q %s%02X\n", s, pre, hash(s));
	}
	Bterm(b);
}

void
mvdirs() {
	int i;

	for(i = 0; i < hashsz; i++)
		print("mv %s%02X %02X\n", pre, i, i);
}

int
readint(char *file) {
	char *s;
	Biobuf *b;

	b = Bopen(file, O_RDONLY);
	s = Brdstr(b, '\n', 1);
	Bterm(b);
	return atoi(s);
}

void
analyze(int n) {
	int i, h, p, r;

	if(n <= Hlimit) {
		fprint(2, "literal %d\n", n);
		print("echo literal %d > «format»\n", n);
		print("echo . > «dirs»\n");
		return;
	}
	r = sqrt(n);
	fprint(2, "sqrt = %d\n", r);
	h = primes[0];
	for(i = 1; i < Nprimes; i++) {
		p = primes[i];
		if(p < r)
			h = p;
		else
			break;
	}
	fprint(2, "chose %d * %f\n", h, (double)n/(double)h);
	hashsz = h;
	print("echo hash %d %d > «format»\n", h, n);
	mkdirs();
	hashfiles();
	mvdirs();
}

void
checkfmt(void) {
	uchar buff[1024];

	if(stat("«format»", buff, sizeof(buff)) >= 0) {
		fprint(2, "already formatted\n");
		exits("formatted");
	}
}

void
main(void) {
	int count;

	quotefmtinstall();
	checkfmt();
	count = readint("«count»");
	fprint(2, "count %d\n", count);
	analyze(count);
	exits(nil);
}
