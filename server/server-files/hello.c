#include <stdio.h>
#include <unistd.h>

int main() {
   printf("Hello from a compiled C binary!\n");
   printf("I will close in 5 seconds...\n");
   fflush(stdout);
   sleep(5);
   return 0;
}
