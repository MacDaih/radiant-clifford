#include <dht11.h>
#include <string.h>

#define DHT11_PIN 6
#define A3 3

dht11 DHT;

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
}

void loop() {
  int light = analogRead(A3);
  
  int chk = DHT.read(DHT11_PIN);
  float temp = float(DHT.temperature);
  float hum = float(DHT.humidity);
  char buffer[512];
  sprintf(buffer,"{\"t\":%d.%02d,\"h\":%d.%02d,\"l\":%i}",(int)temp,(int)(temp*100/100),(int)hum,(int)(hum*100/100),light);
  Serial.println(buffer);
  delay(5000);
}
