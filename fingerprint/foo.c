#include <stdio.h>
#include <fcntl.h>
#include <sys/ioctl.h>
#include <mtd/mtd-user.h>
#include <errno.h>
#include <time.h>

int main( void )
{
        int fd;
        char buf[4]="abc";

        fd = open("/dev/ftdi", O_RDWR);
        write(fd, &buf, 4);
	

        sleep(1);
				char buf_r[4] = {0};


				read(fd,&buf_r, 1);

				int i = 0;
				for(; i < 4; i++){
					puts(buf_r[i]);
				}
  



	      close(fd);



        return 0;
}
