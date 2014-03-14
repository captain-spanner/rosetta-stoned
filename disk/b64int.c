#include <u.h>
#include <libc.h>

char	b64[]	= "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

int
b64int(char *s) {
	int i;
	char c;
	char *p;

	i = 0;
	while((c = *s++) != '\0') {
		p = strchr(b64, c);
		if(p == nil) {
			fprint(2, "b64: unknown char '%c'\n", c);
			exits("b64: char");
		}
		i = (i << 6) + (p - b64);
	}
	return i;
}

void
main(int argc, char **argv) {
	if(argc < 2) {
		exits("b64: args");
	}
	print("%d\n", b64int(argv[1]));
	exits("");
}
