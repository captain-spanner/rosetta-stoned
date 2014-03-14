#include <u.h>
#include <libc.h>

char	buff[64*1024*1024];

void
main(void) {
	int fd, i, n;

	i = 0;
	for(;;) {
		n = read(0, &buff[i], 1);
		if(n < 1 || buff[i] == '\n') {
			n = i+1;
			break;
		}
		i++;
	}
	fd = create("«line»", O_WRONLY, 0664);
	write(fd, buff, n);
}
