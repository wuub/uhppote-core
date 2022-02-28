#pragma once

#include <stdbool.h>
#include <stdint.h>

typedef struct controller {
    uint32_t id;
    const char *address;
} controller;

typedef struct device {
    uint32_t ID;
    char address[16];
    char subnet[16];
    char gateway[16];
    char MAC[18];
    char version[6];
    char date[11];
} device;

void setup(const char *bind, const char *broadcast, const char *listen, int timeout, int debug, ...);
void teardown();

const char *errmsg();

int get_devices(uint32_t **devices, int *N);
int get_device(uint32_t id, struct device *);
