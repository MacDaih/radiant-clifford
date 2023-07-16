#include <stdlib.h>
#include <stdio.h>

int main() {
    printf("utilities...\n");

    FILE * f;
    char line[256];
    f = popen("vcgencmd","");

    if (f == NULL)
        return 1;

    while (fgets(line, sizeof(line), f) != NULL) {
    printf("%s", f);
  }
    return 0;
}
