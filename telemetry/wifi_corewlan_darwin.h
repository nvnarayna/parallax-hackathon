#ifndef WIFI_COREWLAN_DARWIN_H
#define WIFI_COREWLAN_DARWIN_H

typedef struct {
    char *ssid;
    char *bssid;
    char *band;

    int rssi;
    int rssi_valid;

    int noise;
    int noise_valid;

    int channel;
    int channel_valid;

    int frequency;
    int frequency_valid;

    double link_mbps;
    int link_mbps_valid;
} WiFiInfo;

int getWiFiInfo(WiFiInfo *info);
void freeWiFiInfo(WiFiInfo *info);

#endif