#include <stdio.h>
#include <string.h>

void calibrate(char *buf) {
    printf("Input calibration parameters:\n");
    fgets(buf, 128, stdin);
    printf("Calibration parameters received:\n");
    printf(buf);
    FILE *f = fopen("params.txt", "w");
    fputs(buf, f);
    fclose(f);
}

int main(int argc, char *argv[]) {
    char buf[128];
    calibrate(buf);
    return 0;
}