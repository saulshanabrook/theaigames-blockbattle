#include <unistd.h>

int main(){
	char *argv[] = {NULL};
	execv("/src/binary/bot", argv);
	return 0;
}
