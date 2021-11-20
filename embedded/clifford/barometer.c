#include <wiringPi.h>
#include <wiringPiI2C.h>
#include <math.h>
#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include "LPS22HB.h"

#define BUFFER_LENGTH   6

int fd;
unsigned char u8Buf[3];
float PRESS_DATA=0;
float TEMP_DATA=0;

char I2C_readByte(int reg) {
	return wiringPiI2CReadReg8(fd, reg);
}

unsigned short I2C_readU16(int reg) {
	return wiringPiI2CReadReg16(fd, reg);
}

void I2C_writeByte(int reg, int val) {
	wiringPiI2CWriteReg8(fd, reg, val);
}
void LPS22HB_RESET() {   unsigned char Buf;
    Buf=I2C_readU16(LPS_CTRL_REG2);
    Buf|=0x04;                                         
    I2C_writeByte(LPS_CTRL_REG2,Buf);                  //SWRESET Set 1
    while(Buf)
    {
        Buf=I2C_readU16(LPS_CTRL_REG2);
        Buf&=0x04;
    }
}

void LPS22HB_START_ONESHOT() {
    unsigned char Buf;
    Buf=I2C_readU16(LPS_CTRL_REG2);
    Buf|=0x01;                                         //ONE_SHOT Set 1
    I2C_writeByte(LPS_CTRL_REG2,Buf);
}

unsigned char LPS22HB_INIT() {
    fd=wiringPiI2CSetup(LPS22HB_I2C_ADDRESS);
    if(I2C_readByte(LPS_WHO_AM_I)!=LPS_ID) return 0;    //Check device ID 
    LPS22HB_RESET();                                    //Wait for reset to complete
    I2C_writeByte(LPS_CTRL_REG1 ,   0x02);              //Low-pass filter disabled , output registers not updated until MSB and LSB have been read , Enable Block Data Update , Set Output Data Rate to 0 
    return 1;
}

void serve_sock(char *sock, char *key) {
    int    sd=-1, sd2=-1;
    int    rc, length;
    char   buffer[BUFFER_LENGTH];
    struct sockaddr_un serveraddr;

    sd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (sd < 0)
        perror("socket() failed");

    memset(&serveraddr, 0, sizeof(serveraddr));
    serveraddr.sun_family = AF_UNIX;
    strcpy(serveraddr.sun_path, sock);

    rc = bind(sd, (struct sockaddr *)&serveraddr, SUN_LEN(&serveraddr));
    if (rc < 0)
        perror("bind() failed");

    rc = listen(sd, 10);

    if (rc < 0)
        perror("listen() failed");

    sd2 = accept(sd, NULL, NULL);
    if (sd2 < 0)
        perror("accept() failed");

    length = BUFFER_LENGTH;
    rc = setsockopt(sd2, SOL_SOCKET, SO_RCVLOWAT,(char *)&length, sizeof(length));
    if (rc < 0)
        perror("setsockopt(SO_RCVLOWAT) failed");

    for (;;) {
        LPS22HB_START_ONESHOT();        //Trigger one shot data acquisition
        if((I2C_readByte(LPS_STATUS)&0x01)==0x01)   //a new pressure data is generated
        {
            u8Buf[0]=I2C_readByte(LPS_PRESS_OUT_XL);
            u8Buf[1]=I2C_readByte(LPS_PRESS_OUT_L);
            u8Buf[2]=I2C_readByte(LPS_PRESS_OUT_H);
            PRESS_DATA=(float)((u8Buf[2]<<16)+(u8Buf[1]<<8)+u8Buf[0])/4096.0f;
        }
        if((I2C_readByte(LPS_STATUS)&0x02)==0x02) {
            u8Buf[0]=I2C_readByte(LPS_TEMP_OUT_L);
            u8Buf[1]=I2C_readByte(LPS_TEMP_OUT_H);
            TEMP_DATA=(float)((u8Buf[1]<<8)+u8Buf[0])/100.0f;
        }

        char str[80];
        sprintf(str, "{\"p\":%6.2f,\"dt\":%6.2f}",  PRESS_DATA, TEMP_DATA);
        
        rc = recv(sd2, key, sizeof(key), 0);
        if (rc < 0){
            perror("recv() failed");
            break;
        } 

        if (rc == 0 || rc < sizeof(buffer)) {
            printf("Connection closed with client");
            break;
        }

        rc = send(sd2, str, sizeof(str), 0);
        if (rc < 0) {
            perror("send() failed");
            break;
        }
    }
    if (sd != -1)
        close(sd);

    if (sd2 != -1)
        close(sd2);

    unlink(sock);
}

int main() {  

    if (wiringPiSetup() < 0) return 1;

    if(!LPS22HB_INIT())
    {
        printf("\nPressure Sensor Error\n");
        return 0;
    }
    char socket[] = "/tmp/barom.sock";
    char key = "barom";
    serve_sock(socket, key);
    return 0;
}