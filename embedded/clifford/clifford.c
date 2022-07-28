#include <wiringPi.h>
#include <wiringPiI2C.h>
#include <math.h>
#include <stdio.h>
#include <string.h>
#include <sys/ioctl.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <netdb.h>
#include <netinet/in.h>
#include <stdlib.h>
#include"AD.h"
#include "LPS22HB.h"
#include "SHTC3.h"

#define MAX 80
#define PORT 8080
#define SA struct sockaddr

struct ain_data {
    float temp;
    float lux;
};

int fd,fp,fad;

int Config_Set;

float TH_Value,RH_Value, PRESS_DATA;
unsigned char u8Buf[3];
char checksum;
char SDA = 8;
char SCL = 9;

char temphum[80];
struct ain_data ad = {0};

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
}

void SHTC3_WAKEUP() {     
    SHTC3_WriteCommand(SHTC3_WakeUp);                  // write wake_up command  
    delayMicroseconds(300);                          //Delay 300us
}
void SHTC3_SLEEP() {    
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
    read(fd, buf, sizeof(buf));

    checksum=buf[2];
    if(!SHTC3_CheckCrc(buf,2,checksum))
        TH_DATA=(buf[0]<<8|buf[1]);
        
    SHTC3_WriteCommand(SHTC3_NM_CD_ReadRH);                 //Read temperature first,clock streching disabled (polling)
    delay(20);
    read(fd, buf, 3);

    checksum=buf[2];
    if(!SHTC3_CheckCrc(buf,2,checksum))
        RH_DATA=(buf[0]<<8|buf[1]);

    TH_Value=175 * ((float)TH_DATA / 65536.0f) - 45.0f;       //Calculate temperature value
    RH_Value=100 * (float)RH_DATA / 65536.0f;              //Calculate humidity value     
}

char I2C_readByte(int reg)
{
	return wiringPiI2CReadReg8(fp, reg);
}

unsigned short I2C_readU16(int reg)
{
	return wiringPiI2CReadReg16(fp, reg);
}

int I2C_AD_readU16(int reg)
{   int val;
    unsigned char Val_L,Val_H;
    val=wiringPiI2CReadReg16(fad,reg);                    //High and low bytes are the opposite       
    Val_H=val&0xff;
    Val_L=val>>8;
    val=(Val_H<<8)|Val_L;                               //Correct byte order
    return val;
}

void I2C_writeWord(int reg, int val)
{   unsigned char Val_L,Val_H;
    Val_H=val&0xff;
    Val_L=val>>8;
    val=(Val_H<<8)|Val_L;                               ////Correct byte order
	wiringPiI2CWriteReg16(fad,reg,val);
}

unsigned int ADS1015_INIT()
{   unsigned int state;
    state=I2C_AD_readU16(ADS_POINTER_CONFIG) & 0x8000  ;
    return state;
}

unsigned int ADS1015_SINGLE_READ(unsigned int channel)           //Read single channel data
{   
    unsigned int data;
    Config_Set = ADS_CONFIG_MODE_NOCONTINUOUS        |   //mode：Single-shot mode or power-down state    (default)
                 ADS_CONFIG_PGA_4096                 |   //Gain= +/- 4.096V                              (default)
                 ADS_CONFIG_COMP_QUE_NON             |   //Disable comparator                            (default)
                 ADS_CONFIG_COMP_NONLAT              |   //Nonlatching comparator                        (default)
                 ADS_CONFIG_COMP_POL_LOW             |   //Comparator polarity：Active low               (default)
                 ADS_CONFIG_COMP_MODE_TRADITIONAL    |   //Traditional comparator                        (default)
                 ADS_CONFIG_DR_RATE_1600             ;   //Data rate=1600SPS                             (default)
    switch (channel)
    {
        case (0):
            Config_Set |= ADS_CONFIG_MUX_SINGLE_0;
            break;
        case (1):
            Config_Set |= ADS_CONFIG_MUX_SINGLE_1;
            break;
        case (2):
            Config_Set |= ADS_CONFIG_MUX_SINGLE_2;
            break;
        case (3):
            Config_Set |= ADS_CONFIG_MUX_SINGLE_3;
            break;
    }
    Config_Set |=ADS_CONFIG_OS_SINGLE_CONVERT;
    I2C_writeWord(ADS_POINTER_CONFIG,Config_Set);
    delay(2);
    data=I2C_AD_readU16(ADS_POINTER_CONVERT)>>4;
    return data;
}

void I2C_writeByte(int reg, int val)
{
	wiringPiI2CWriteReg8(fp, reg, val);
}

void LPS22HB_RESET()
{   unsigned char Buf;
    Buf=I2C_readU16(LPS_CTRL_REG2);
    Buf|=0x04;                                         
    I2C_writeByte(LPS_CTRL_REG2,Buf);                  //SWRESET Set 1
    while(Buf)
    {
        Buf=I2C_readU16(LPS_CTRL_REG2);
        Buf&=0x04;
    }
}
void LPS22HB_START_ONESHOT()
{
    unsigned char Buf;
    Buf=I2C_readU16(LPS_CTRL_REG2);
    Buf|=0x01;                                         //ONE_SHOT Set 1
    I2C_writeByte(LPS_CTRL_REG2,Buf);
}
unsigned char LPS22HB_INIT()
{
    fp=wiringPiI2CSetup(LPS22HB_I2C_ADDRESS);
    if(I2C_readByte(LPS_WHO_AM_I)!=LPS_ID) return 0;    //Check device ID 
    LPS22HB_RESET();                                    //Wait for reset to complete
    I2C_writeByte(LPS_CTRL_REG1 ,   0x02);              //Low-pass filter disabled , output registers not updated until MSB and LSB have been read , Enable Block Data Update , Set Output Data Rate to 0 
    return 1;
}

void read_interface() {

    // device voltage
    int const AREF = 3.3;

    int  AIN0_DATA,AIN1_DATA;
    fad=wiringPiI2CSetup(ADS_I2C_ADDRESS);
    if(ADS1015_INIT()!=0x8000) {	
		return 1;
	}
    // AIN0 is dedicated to temp. sensor
    AIN0_DATA=ADS1015_SINGLE_READ(0);
    ad.temp = (float)(AIN0_DATA * 2) / 10;

    // AIN0 is dedicated to lux / ambient light sensor
    AIN1_DATA=ADS1015_SINGLE_READ(1);
    float amps = (AIN1_DATA * AREF / 1024.0) / 10000.0;
    float microamps = amps * 1000000.0;
    ad.lux = microamps * 2.0;
}

void serve_tcp() {
    int sockfd = -1;
    int conn, len;
    struct sockaddr_in addr, cli;
   
    sockfd = socket(AF_INET, SOCK_STREAM, 0);

    if (sockfd == -1) {
        perror("failed to create socket \n");
        exit(0);
    }

    bzero(&addr, sizeof(addr));

    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(PORT);

    if (bind(sockfd,(SA*)&addr,sizeof(addr)) != 0) {
        perror("failed to bind socket\n");
        exit(0);
    }

    if (listen(sockfd, 1) != 0) {
        perror("failed to listen\n");
        exit(0);
    }

    len = sizeof(cli);

    char buffer[MAX];
    conn = accept(sockfd, (SA*)&cli,&len);
    if (conn < 0) {
        perror("accepting connection failed");
        exit(0);
    }

    while(1) {

        SHTC3_Read_DATA();
        SHTC3_SLEEP();
        SHTC3_WAKEUP();

        LPS22HB_START_ONESHOT();        //Trigger one shot data acquisition
        if((I2C_readByte(LPS_STATUS)&0x01)==0x01)   //a new pressure data is generated
        {
            u8Buf[0]=I2C_readByte(LPS_PRESS_OUT_XL);
            u8Buf[1]=I2C_readByte(LPS_PRESS_OUT_L);
            u8Buf[2]=I2C_readByte(LPS_PRESS_OUT_H);
            PRESS_DATA=(float)((u8Buf[2]<<16)+(u8Buf[1]<<8)+u8Buf[0])/4096.0f;
        }
        
        read_interface();
        char str[256];
        sprintf(str, "{\"press\":%6.2f,\"inner_temp\":%6.2f,\"hum\":%6.2f, \"temp\": %6.2f, \"light\": %6.2f}",  PRESS_DATA, TH_Value, RH_Value, ad.temp, ad.lux);
        bzero(buffer, MAX);
        int r = recv(conn,buffer,sizeof(buffer),0);

        if (r < 0 ) {
            perror("read message error");
            break;
        }

        if (send(conn, str, sizeof(str),0) < 0) {
            perror("sending message error");
            break;
        }
    }
    close(sockfd);
}

int main() {  
    if (wiringPiSetup() < 0) return 1;

    if(!LPS22HB_INIT())
    {
        printf("\nPressure Sensor Error\n");
        return 0;
    }

    if(ADS1015_INIT()!=0x8000) {	
        printf("\nADS1015 interfacing Error\n");
		return 0;
	}
    
    fd=wiringPiI2CSetup(SHTC3_I2C_ADDRESS);

    SHTC_SOFT_RESET();

    serve_tcp();
    return 1;
}