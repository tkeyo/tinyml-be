## TinyML Backend

RestAPI for the TinyML projects. Implemented with Go using GinGonic. Uses AWS DynamoDB as backend.

**Roles:**
1. Saves data comming from an ESP32 IoT device that captures acceleration & gyro data, and identified gestures.
2. Retrieves data from the DB to show on the frontend. 

## Local Docker Run

```
docker run -it -p 8081:8081 -e API_AUTH_KEY=123 -e GIN_MODE=release -e AWS_ACCESS_KEY_ID=xyz AWS_SECRET_ACCESS_KEY=xyz tinymlbe:latest
```
