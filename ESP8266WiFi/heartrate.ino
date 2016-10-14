#define heartPin A0
#include <ESP8266WiFi.h>


void setup() {
  Serial.begin(115200);
}
void loop() {
  int heartValue = analogRead(heartPin);
  Serial.println(heartValue);
  delay(20);
}   
