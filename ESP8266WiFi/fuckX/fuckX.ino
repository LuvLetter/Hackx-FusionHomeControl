#include <ArduinoJson.h>

#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <OneWire.h>
#define heartPin D0
OneWire  ds(D1);  // on pin D1

const char* ssid = "Hackx";
const char* password = "fudan2016";
float celsius, fahrenheit;
float rate;
int minTime;
int count;
boolean isHigh;
void flush() {
  rate = count*3;
  count = 0;
}
ESP8266WebServer server(80);

//root page can be accessed only if authentification is ok
void handleRoot(){
  Serial.println("Enter handleRoot");

  String content = "<html><body><H2>Success</H2><br>";

  char ptr[25];
  dtostrf(celsius,2,3,ptr);
  content += "<h3>Temperature = ";
  content += ptr;
  content += "</h3><br>";

  char ptr2[25];
  itoa(digitalRead(D2),ptr2,10);
  content += "<h3>hasPeopleMoved = ";
  content += ptr2;
  content += "</h3><br>";

  char ptr3[25];
  itoa(analogRead(A0),ptr3,10);
  content += "<h3>CO = ";
  content += ptr3;
  content += "</h3><br>";


  char ptr4[25];
  itoa(rate,ptr4,10);
  content += "<h3>heartRate = ";
  content += ptr4;
  content += "</h3><br>";

  server.send(200, "text/html", content);
}

void retJson() {
  Serial.println("sending json");
  String content;
  content += "{\"Temperature\":";
  content += celsius;
  content += ",";

  content += "\"hasPeopleMoved\":";
  if(digitalRead(D2) == 1){
    content += "true";
  }
  else{
    content += "false" ;
  }
  content += ",";

  content += "\"CO\":";
  content += analogRead(A0);
  content += ",";

  content += "\"heartRate\":";
  content += rate;
  content += "}";

  server.send(200,"application/json",content);
}

//no need authentification

void setup(void){
  isHigh = false;
  Serial.begin(115200);
  WiFi.begin(ssid, password);
  Serial.println("");

  // Wait for connection
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.print("Connected to ");
  Serial.println(ssid);
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());


  server.on("/", handleRoot);
  server.on("/json", retJson);
  server.begin();
  Serial.println("HTTP server started");
}

void loop(void){
  byte i;
  byte present = 0;
  byte type_s;
  byte data[12];
  byte addr[8];

  if ( !ds.search(addr)) {
    Serial.println("No more addresses.");
    Serial.println();
    ds.reset_search();
    delay(250);
    return;
  }


  if (OneWire::crc8(addr, 7) != addr[7]) {
    return;
  }
  type_s = 0;

  ds.reset();
  ds.select(addr);
  ds.write(0x44, 1);        // start conversion, with parasite power on at the end

  delay(1000);     // maybe 750ms is enough, maybe not
  // we might do a ds.depower() here, but the reset will take care of it.

  present = ds.reset();
  ds.select(addr);
  ds.write(0xBE);         // Read Scratchpad

  Serial.print(" ");
  for ( i = 0; i < 9; i++) {           // we need 9 bytes
    data[i] = ds.read();
  }
//  Serial.print(OneWire::crc8(data, 8), HEX);
//  Serial.println();

  int16_t raw = (data[1] << 8) | data[0];
  if (type_s) {
    raw = raw << 3; // 9 bit resolution default
    if (data[7] == 0x10) {
      // "count remain" gives full 12 bit resolution
      raw = (raw & 0xFFF0) + 12 - data[6];
    }
  } else {
    byte cfg = (data[4] & 0x60);
    // at lower res, the low bits are undefined, so let's zero them
    if (cfg == 0x00) raw = raw & ~7;  // 9 bit resolution, 93.75 ms
    else if (cfg == 0x20) raw = raw & ~3; // 10 bit res, 187.5 ms
    else if (cfg == 0x40) raw = raw & ~1; // 11 bit res, 375 ms
    //// default is 12 bit resolution, 750 ms conversion time
  }
  celsius = (float)raw / 16.0;
  fahrenheit = celsius * 1.8 + 32.0;
  Serial.println(celsius);
  Serial.println(digitalRead(D2));
  Serial.println(analogRead(A0));
  minTime=millis();
  while((millis()-minTime)<20000) {
    int heartValue = analogRead(heartPin);
    if(heartValue > 600&&isHigh==false) {
      isHigh=true;
      count++;
    }
    if(heartValue < 600&&isHigh==true) {
      isHigh=false;
    }
    delay(20);
  }
  flush();
  server.handleClient();
}
