#include <WiFi.h>
#include <PubSubClient.h>
#include <Wire.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_MPU6050.h>

// Wi-Fi and MQTT Credentials
const char* ssid = "Your_WiFi_SSID";  // Replace with your WiFi SSID
const char* password = "Your_WiFi_Password";  // Replace with your WiFi password
const char* mqttServer = "broker.hivemq.com";  // Replace with your broker
const int mqttPort = 1883;
const char* mqttTopic = "robotic_arm/commands";

// MQTT Client
WiFiClient espClient;
PubSubClient mqttClient(espClient);

// Pin Definitions for Sensors
#define FLEX_SENSOR_1_PIN 36
#define FLEX_SENSOR_2_PIN 39
#define TOUCH_SENSOR_1_PIN 4
#define TOUCH_SENSOR_2_PIN 14
#define MPU_SDA 21
#define MPU_SCL 22

// Gesture Codes
#define GESTURE_FIST 0x01
#define GESTURE_OPEN_HAND 0x02
#define GESTURE_LEFT 0x03
#define GESTURE_RIGHT 0x04

// Function Prototypes
void connectToWiFi();
void connectToMQTT();
void ensureWiFiConnection();
void ensureMQTTConnection();
int readFlexSensor(int pin);
void readMPU6050(float &x, float &y, float &z);
int debounceTouchSensor(int pin);
int detectGesture();
void publishGesture(int gestureCode);
void debugSensors(float flex1, float flex2, float accelX, float accelY, float accelZ, int touch1, int touch2);

// MPU6050 Configuration
Adafruit_MPU6050 mpu;

// Timing for non-blocking delay
unsigned long lastReadTime = 0;
const unsigned long readInterval = 500;  // 500 ms interval

// Setup Function
void setup() {
    Serial.begin(115200);

    // Initialize WiFi
    connectToWiFi();

    // Initialize MQTT
    mqttClient.setServer(mqttServer, mqttPort);

    connectToMQTT();

    // Initialize MPU6050
    if (!mpu.begin()) {
        Serial.println("Failed to find MPU6050 chip.");
        while (1);
    }
    mpu.setAccelerometerRange(MPU6050_RANGE_8_G);
    mpu.setGyroRange(MPU6050_RANGE_500_DEG);
    mpu.setFilterBandwidth(MPU6050_BAND_21_HZ);
}

// Main Loop Function
void loop() {
    ensureWiFiConnection();
    ensureMQTTConnection();

    unsigned long currentTime = millis();
    if (currentTime - lastReadTime >= readInterval) {
        lastReadTime = currentTime;
        int gestureCode = detectGesture();
        if (gestureCode != 0) {
            publishGesture(gestureCode);
        }
    }
}

// Function Definitions
void connectToWiFi() {
    Serial.print("Connecting to Wi-Fi: ");
    Serial.println(ssid);
    WiFi.begin(ssid, password);
    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }
    Serial.println("\nWi-Fi Connected.");
}

void ensureWiFiConnection() {
    if (WiFi.status() != WL_CONNECTED) {
        connectToWiFi();
    }
}

void connectToMQTT() {
    while (!mqttClient.connected()) {
        Serial.print("Connecting to MQTT...");
        if (mqttClient.connect("ESP32")) {
            Serial.println("Connected to MQTT Broker.");
        } else {
            Serial.print("Failed. State=");
            Serial.println(mqttClient.state());
            delay(2000);
        }
    }
}

void ensureMQTTConnection() {
    if (!mqttClient.connected()) {
        connectToMQTT();
    }
}

int readFlexSensor(int pin) {
    int value = analogRead(pin);
    return value;
}

void readMPU6050(float &x, float &y, float &z) {
    sensors_event_t a, g, temp;
    mpu.getEvent(&a, &g, &temp);
    x = a.acceleration.x;
    y = a.acceleration.y;
    z = a.acceleration.z;
}

int debounceTouchSensor(int pin) {
    int stableValue = digitalRead(pin);
    delay(50);  // Debounce delay
    if (digitalRead(pin) == stableValue) {
        return stableValue;
    }
    return LOW;  // Default to no touch
}

int detectGesture() {
    // Read sensors
    int flex1 = readFlexSensor(FLEX_SENSOR_1_PIN);
    int flex2 = readFlexSensor(FLEX_SENSOR_2_PIN);
    float accelX, accelY, accelZ;
    readMPU6050(accelX, accelY, accelZ);
    int touch1 = debounceTouchSensor(TOUCH_SENSOR_1_PIN);
    int touch2 = debounceTouchSensor(TOUCH_SENSOR_2_PIN);

    // Debugging sensors
    debugSensors(flex1, flex2, accelX, accelY, accelZ, touch1, touch2);

    // Detect gestures
    if (flex1 < 500 && flex2 < 500 && touch1 == HIGH) {
        return GESTURE_FIST;
    } else if (flex1 > 2000 && flex2 > 2000 && touch2 == HIGH) {
        return GESTURE_OPEN_HAND;
    } else if (accelX < -5.0) {
        return GESTURE_LEFT;
    } else if (accelX > 5.0) {
        return GESTURE_RIGHT;
    }
    return 0;  // No gesture detected
}

void debugSensors(float flex1, float flex2, float accelX, float accelY, float accelZ, int touch1, int touch2) {
    Serial.print("Flex1: "); Serial.print(flex1);
    Serial.print(", Flex2: "); Serial.print(flex2);
    Serial.print(", AccelX: "); Serial.print(accelX);
    Serial.print(", AccelY: "); Serial.print(accelY);
    Serial.print(", AccelZ: "); Serial.print(accelZ);
    Serial.print(", Touch1: "); Serial.print(touch1);
    Serial.print(", Touch2: "); Serial.println(touch2);
}

void publishGesture(int gestureCode) {
    char message[10];
    snprintf(message, sizeof(message), "0x%02X", gestureCode);
    if (!mqttClient.publish(mqttTopic, message)) {
        Serial.println("Failed to publish gesture code!");
    } else {
        Serial.print("Published gesture code: ");
        Serial.println(message);
    }
}
