#include <u.h>
#include <libc.h>
#include <bio.h>

void
main(int argc, char **argv) {
	int i;

	quotefmtinstall();
	if(argc < 2) {
		int l;
		char *s;
		Biobuf *b;

		b = Bfdopen(0, O_RDONLY);
		while((s = Brdstr(b, '\n', 1)) != nil) {
			l = Blinelen(b);
			if(s[l-1] == '\n')
				s[l-1] = '\0';
			print("%#q\n", s);
		}
		Bterm(b);
		exits("");
	}
	for(i = 1; i < argc; i++)
		print("%#q\n", argv[i]);
}
