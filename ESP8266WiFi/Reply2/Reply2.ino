/**
    Get the Weather from intenet with esp8266
 
    get data from 心知天气：[url=http://www.thinkpage.cn/]http://www.thinkpage.cn/[/url]
    api 文档说明：[url=http://www.thinkpage.cn/doc]http://www.thinkpage.cn/doc[/url]
    city id list download ：[url=http://www.thinkpage.cn/data/thinkpage_cities.zip]http://www.thinkpage.cn/data/thinkpage_cities.zip[/url]
 
    Created by yfrobot, 2016.8.23
    This example is in public domain.
*/
 
#include <ESP8266WiFi.h>
#include <ArduinoJson.h>
 
WiFiClient client;
 
const char* ssid     = "Hackx";         
const char* password = "fudan2016";         
 
 
const char* host = "10.221.65.97";
const unsigned long BAUD_RATE = 115200;                   
const unsigned long HTTP_TIMEOUT = 2100;               
const size_t MAX_CONTENT_SIZE = 2048;                  
 
bool skipResponseHeaders() {
  // HTTP headers end with an empty line
  char endOfHeaders[] = "\r\n\r\n";
 
  client.setTimeout(HTTP_TIMEOUT);
  bool ok = client.find(endOfHeaders);
 
  if (!ok) {
    Serial.println("No response or invalid response!");
  }
 
  return ok;
}
 
// 发送请求指令
bool sendRequest(const char* host) {
  String GetUrl = "/json";
  client.print(String("GET ") + GetUrl + " HTTP/1.1\r\n" +
               "Host: " + host + "\r\n" +
               "Connection: close\r\n\r\n");
  return true;
}
 
// Read the body of the response from the HTTP server -- 从HTTP服务器响应中读取正文
void readReponseContent(char* content, size_t maxSize) {
  size_t length = client.peekBytes(content, maxSize);
  delay(100);
  Serial.println("Get the data from Internet!");
  content[length] = 0;
  Serial.println(content);
  Serial.println("Read Over!");
}

struct UserData {
  char city[16];
  char cnty[16];
  char weather[32];
  char temp[16];
  char feel[16];
  char hum[16];
  char visi[16];
  char udate[32];
};
 
// 解析数据
bool parseUserData(char* content, struct UserData* userData) {
  DynamicJsonBuffer jsonBuffer;
 
  JsonObject& root = jsonBuffer.parseObject(content);
 
  if (!root.success()) {
    Serial.println("JSON parsing failed!");
    return false;
  }
  //  const char* x = root["results"][0]["location"]["name"];//
  //  Serial.println(x);
  // Here were copy the strings we're interested in -- 复制我们感兴趣的字符串
  strcpy(userData->city, root["results"][0]["location"]["name"]);
  strcpy(userData->cnty, root["results"][0]["location"]["country"]);
  strcpy(userData->weather, root["results"][0]["now"]["text"]);
  strcpy(userData->temp, root["results"][0]["now"]["temperature"]);
  strcpy(userData->feel, root["results"][0]["now"]["feels_like"]);
  strcpy(userData->hum, root["results"][0]["now"]["humidity"]);
  strcpy(userData->visi, root["results"][0]["now"]["visibility"]);
  strcpy(userData->udate, root["results"][0]["last_update"]);
  // It's not mandatory to make a copy, you could just use the pointers
  // Since, they are pointing inside the "content" buffer, so you need to make
  // sure it's still in memory when you read the string
  //  -- 这不是强制复制，你可以使用指针，因为他们是指向“内容”缓冲区内，所以你需要确保
  //   当你读取字符串时它仍在内存中
 
  return true;
}
 
// Print the data extracted from the JSON -- 打印从JSON中提取的数据
void printUserData(const struct UserData* userData) {
  Serial.println("Print parsed data :");
  Serial.print("City : ");
  Serial.print(userData->city);
  Serial.print(", \t");
  Serial.print("Country : ");
  Serial.println(userData->cnty);
   
  Serial.print("Weather : ");
  Serial.print(userData->weather);
  Serial.print(",\t");
  Serial.print("Temp : ");
  Serial.print(userData->temp);
  Serial.print(" C");
  Serial.print(",\t");
  Serial.print("Feel : ");
  Serial.print(userData->feel);
  Serial.print(" C");
  Serial.print(",\t");
  Serial.print("Humidity : ");
  Serial.print(userData->hum);
  Serial.print(" %");
  Serial.print(",\t");
  Serial.print("visibility : ");
  Serial.print(userData->visi);
  Serial.println(" km");
   
  Serial.print("Last Updata : ");
  Serial.print(userData->udate);
  Serial.println("");
}
 
// Close the connection with the HTTP server -- 关闭与HTTP服务器连接
void stopConnect() {
  Serial.println("Disconnect");
  client.stop();
}
 
void setup(){
  WiFi.mode(WIFI_STA);     //设置esp8266 工作模式
 
  Serial.begin(BAUD_RATE );      //设置波特率
  Serial.println();
  Serial.print("connecting to ");
  Serial.println(ssid);
 
  WiFi.begin(ssid, password);   //连接wifi
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.println("WiFi connected");
  delay(500);
  // Check if a client has connected
  if (!client.connect(host, 80)) {
    Serial.println("connection failed");
    return;
  }
 
  if (sendRequest(host) && skipResponseHeaders()) {
    char response[MAX_CONTENT_SIZE];
    readReponseContent(response, sizeof(response));
    UserData userData;
    if (parseUserData(response, &userData)) {
      printUserData(&userData);
    }
  }
  stopConnect();
}
void loop(){
  delay(3000);
}
