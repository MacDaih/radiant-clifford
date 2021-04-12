#include <dht11.h>

#define DHT11_PIN 6
#define A3 3

dht11 DHT;

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
}

void loop() {
  int val = analogRead(A3);

  
  int chk = DHT.read(DHT11_PIN);
  Serial.print("Temp. => ");
  Serial.println(DHT.temperature,1);
  Serial.print("Hum. => ");
  Serial.println(DHT.humidity,1);
  Serial.print("Light => ");
  Serial.println(val,DEC);
  delay(5000);
}
