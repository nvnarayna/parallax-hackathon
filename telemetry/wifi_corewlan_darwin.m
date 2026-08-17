//go:build darwin

#import <CoreWLAN/CoreWLAN.h>
#import <Foundation/Foundation.h>

#include <stdlib.h>
#include <string.h>

#include "wifi_corewlan_darwin.h"

static char *copyNSString(NSString *value) {
    if (value == nil) {
        return NULL;
    }

    const char *utf8 = [value UTF8String];

    if (utf8 == NULL) {
        return NULL;
    }

    return strdup(utf8);
}

/*
 * Convert Wi-Fi channel number to center frequency in MHz.
 *
 * 2.4 GHz:
 *   channel 1 = 2412 MHz
 *   channel 6 = 2437 MHz
 *   channel 11 = 2462 MHz
 *
 * 5 GHz:
 *   channel 36 = 5180 MHz
 *   channel 40 = 5200 MHz
 *   ...
 *
 * 6 GHz:
 *   channel 1 = 5955 MHz
 *   channel 5 = 5975 MHz
 *   ...
 */
static int frequencyForChannel(CWChannel *channel) {
    if (channel == nil) {
        return 0;
    }

    NSInteger number = [channel channelNumber];
    CWChannelBand band = [channel channelBand];

    if (number <= 0) {
        return 0;
    }

    switch (band) {
        case kCWChannelBand2GHz:
            if (number == 14) {
                return 2484;
            }

            if (number >= 1 && number <= 13) {
                return 2407 + (number * 5);
            }

            break;

        case kCWChannelBand5GHz:
            return 5000 + (number * 5);

        case kCWChannelBand6GHz:
            return 5950 + (number * 5);

        default:
            break;
    }

    return 0;
}

int getWiFiInfo(WiFiInfo *info) {
    if (info == NULL) {
        return -1;
    }

    memset(info, 0, sizeof(WiFiInfo));

    @autoreleasepool {
        CWWiFiClient *client = [CWWiFiClient sharedWiFiClient];

        if (client == nil) {
            return -1;
        }

        CWInterface *interface = [client interface];

        if (interface == nil) {
            return -1;
        }

        if (![interface powerOn]) {
            return -1;
        }

        /*
         * Basic connection information.
         */
        info->ssid = copyNSString([interface ssid]);
        info->bssid = copyNSString([interface bssid]);

        /*
         * Signal strength.
         */
        NSInteger rssi = [interface rssiValue];

        if (rssi != 0) {
            info->rssi = (int)rssi;
            info->rssi_valid = 1;
        }

        /*
         * Noise floor.
         */
        NSInteger noise = [interface noiseMeasurement];

        if (noise != 0) {
            info->noise = (int)noise;
            info->noise_valid = 1;
        }

        /*
         * Channel / band / frequency.
         */
        CWChannel *channel = [interface wlanChannel];

        if (channel != nil) {
            NSInteger channelNumber = [channel channelNumber];

            if (channelNumber > 0) {
                info->channel = (int)channelNumber;
                info->channel_valid = 1;
            }

            switch ([channel channelBand]) {
                case kCWChannelBand2GHz:
                    info->band = strdup("2.4GHz");
                    break;

                case kCWChannelBand5GHz:
                    info->band = strdup("5GHz");
                    break;

                case kCWChannelBand6GHz:
                    info->band = strdup("6GHz");
                    break;

                default:
                    break;
            }

            int frequency = frequencyForChannel(channel);

            if (frequency > 0) {
                info->frequency = frequency;
                info->frequency_valid = 1;
            }
        }

        /*
         * Current PHY/link rate.
         */
        double transmitRate = [interface transmitRate];

        if (transmitRate > 0) {
            info->link_mbps = transmitRate;
            info->link_mbps_valid = 1;
        }
    }

    /*
     * RSSI is the minimum information required for a valid
     * Wi-Fi measurement.
     */
    if (!info->rssi_valid) {
        freeWiFiInfo(info);
        return -1;
    }

    return 0;
}

void freeWiFiInfo(WiFiInfo *info) {
    if (info == NULL) {
        return;
    }

    free(info->ssid);
    free(info->bssid);
    free(info->band);

    info->ssid = NULL;
    info->bssid = NULL;
    info->band = NULL;
}