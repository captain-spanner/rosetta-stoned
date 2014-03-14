#include <u.h>
#include <libc.h>

void
main(int argc, char **argv) {
	int fd;
	char *p;
	vlong off;
	long size, r;

	if(argc != 4) {
		fprint(2, "usage: getoff file offset size\n");
		exits("usage");
	}
	if(strcmp(argv[1], "-") == 0) {
		fd = 0;
	} else {
		fd = open(argv[1], O_RDONLY);
		if(fd < 0) {
			fprint(2, "getoff: open %s: %r\n", argv[1]);
			exits("open");
		}
	}
	off = atoll(argv[2]);
	size =  atol(argv[3]);
	if(off < 0) {
		fprint(2, "getoff: bad offset: %s\n", argv[2]);
		exits("offset");
	}
	if(size <= 0) {
		fprint(2, "getoff: bad size: %s\n", argv[3]);
		exits("size");
	}
	p = malloc(size);
	r = pread(fd, p, size, off);
	if(r != size) {
		if(r < 0)
	 		fprint(2, "getoff: read %s: %r\n", argv[1]);
		else
	 		fprint(2, "getoff: %s: short read\n", argv[1]);
		exits("read");
	}
	r = write(1, p, size);
	if(r < 0) {
	 	fprint(2, "getoff: write %s: %r\n", argv[1]);
		exits("write");
	}
	exits("");
}
