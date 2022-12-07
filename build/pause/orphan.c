/* Creates a zombie to be reaped by init. Useful for testing. */

#include <stdio.h>
#include <unistd.h>

int main() {
  pid_t pid;
  pid = fork();
  if (pid == 0) {
    while (getppid() > 1)
      ;
    printf("Child exiting: pid=%d ppid=%d\n", getpid(), getppid());
    return 0;
  } else if (pid > 0) {
    printf("Parent exiting: pid=%d ppid=%d\n", getpid(), getppid());
    return 0;
  }
  perror("Could not create child");
  return 1;
}
