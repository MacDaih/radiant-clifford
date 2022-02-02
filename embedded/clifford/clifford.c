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
#include "SHTC3.h"

#define BUFFER_LENGTH 6
#define THERMO "thermo.sock"

int fd;
unsigned char u8Buf[3];
float PRESS_DATA=0;
float TEMP_DATA=0;

float TH_Value,RH_Value;
char checksum;
char SDA = 8;
char SCL = 9;

char temphum[80];

char SHTC3_CheckCrc(char data[],unsigned char len,unsigned char checksum) {
  unsigned char bit;        // bit mask
  unsigned char crc = 0xFF; // calculated checksum
  unsigned char byteCtr;    // byte counter
  // calculates 8-Bit checksum with given polynomial
  for(byteCtr = 0; byteCtr < len; byteCtr++) {
    crc ^= (data[byteCtr]);
    for(bit = 8; bit > 0; --bit) {
      if(crc & 0x80) {
        crc = (crc << 1) ^ CRC_POLYNOMIAL;
      } else {
        crc = (crc << 1);
      }
    }
  }
  // verify checksum
  if(crc != checksum) {                 
    return 1;                       //Error
  } else {
    return 0;                       //No error
  }       
}

void SHTC3_WriteCommand(unsigned short cmd) {   
    char buf[] = { (cmd>>8) ,cmd};
    wiringPiI2CWriteReg8(fd,buf[0],buf[1]);          
                                                 //1:error 0:No error
}

void SHTC3_WAKEUP() {     
    SHTC3_WriteCommand(SHTC3_WakeUp);                  // write wake_up command  
    delayMicroseconds(300);                          //Delay 300us
      
}
void SHTC3_SLEEP() {    
 //   bcm2835_i2c_begin();
    SHTC3_WriteCommand(SHTC3_Sleep);                        // Write sleep command  
}

void SHTC_SOFT_RESET() {   
    SHTC3_WriteCommand(SHTC3_Software_RES);                 // Write reset command
    delayMicroseconds(300);                                 //Delay 300us
}

void SHTC3_Read_DATA() {   
    unsigned short TH_DATA,RH_DATA;
    char buf[3];
   SHTC3_WriteCommand(SHTC3_NM_CD_ReadTH);                 //Read temperature first,clock streching disabled (polling)
    delay(20);
    read(fd, buf, 3);

   checksum=buf[2];
   if(!SHTC3_CheckCrc(buf,2,checksum))
        TH_DATA=(buf[0]<<8|buf[1]);
    
    SHTC3_WriteCommand(SHTC3_NM_CD_ReadRH);                 //Read temperature first,clock streching disabled (polling)
    delay(20);
    read(fd, buf, 3);

    checksum=buf[2];
    if(!SHTC3_CheckCrc(buf,2,checksum))
        RH_DATA=(buf[0]<<8|buf[1]);
    
    TH_Value=175 * (float)TH_DATA / 65536.0f - 45.0f;       //Calculate temperature value
    RH_Value=100 * (float)RH_DATA / 65536.0f;              //Calculate humidity value     
}

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

    chmod(sock, 0777);
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

        SHTC3_Read_DATA();
        SHTC3_SLEEP();
        SHTC3_WAKEUP();

        char str[256];
        sprintf(str, "{\"p\":%6.2f,\"dt\":%6.2f,\"t\":%6.2f,\"h\":%6.2f}",  PRESS_DATA, TEMP_DATA, TH_Value, RH_Value);

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
    fd=wiringPiI2CSetup(SHTC3_I2C_ADDRESS);
    SHTC_SOFT_RESET();

    char socket[] = THERMO;
    char key[] = "thermo";
    serve_sock(socket, key);
    return 0;
}