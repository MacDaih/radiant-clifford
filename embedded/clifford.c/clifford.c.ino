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
  String temp = String(DHT.temperature);
  String hum = String(DHT.humidity);
  char buffer[512];
  sprintf(buffer,"{\"t\":\"%s\",\"h\":\"%s\",\"l\":%i}",temp.c_str(),hum.c_str(),light);
  Serial.println(buffer);
  delay(5000);
}
