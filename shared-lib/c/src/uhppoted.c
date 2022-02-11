#include "../include/uhppoted.h"
#include "libuhppoted.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>

__thread char *err = NULL;
__thread UHPPOTE *u = NULL;

char *errmsg() {
    return err;
}

/* (optional) setup for UHPPOTE network configuration. Defaults to:
 * - bind:        0.0.0.0:0
 * - broadcast:   255.255.255.255:60000
 * - listen:      0.0.0.0:60001
 * - timeout:     5s
 * - controllers: (none)
 * - debug:       false
 *
 * NOTES: 1. https://wiki.sei.cmu.edu/confluence/display/cplusplus/EXP58-CPP.+Pass+an+object+of+the+correct+type+to+va_start
 *        2. https://www.linkedin.com/pulse/modern-c-variadic-functions-how-shoot-yourself-foot-avoid-zinin
 */ 
void setup(const char *bind, const char *broadcast, const char *listen, int timeout, int debug, ...) {
    if (u != NULL) {
        teardown();
    }

    if ((u = (UHPPOTE *) malloc(sizeof(UHPPOTE))) != NULL) {
        u->bind = bind;
        u->broadcast = broadcast;
        u->listen = listen;
        u->timeout = timeout;
        u->devices = NULL;
        u->debug = debug != 0;

        va_list args;
        va_start(args, debug);

        controller *p = va_arg(args, controller *);
        udevice    *q = NULL;
        udevice    *previous = NULL;

        // NTS: duplicates address because the controllers may go out of scope after the invocation of setup(..)
        while(p != NULL) {
            if ((q = (udevice *) malloc(sizeof(udevice))) != NULL) {
                q->id = p->id;
                q->address = strdup(p->address);
                q->next=previous;
                previous = q;
            }

            p = va_arg(args, controller *);
        }
        va_end(args);

        u->devices = q;
    }
}

void teardown() {
    if (u != NULL) {
        udevice *d = u->devices;

        while (d != NULL) {
            udevice *next = d->next;
            free((char *) d->address); // Known to be malloc'd in setup
            free(d);
            d = next;
        }        
    }

    free(u);
}

void set_error(const char *errmsg) {
    unsigned l = strlen(errmsg) + 1;

    if (err != NULL) {
        free(err);
    }

    if ((err = malloc(l)) != NULL) {
        snprintf(err, l, "%s", errmsg);            
    }
}

// All this finagling because you can't return a slice from Go
int get_devices(uint32_t **devices, int *N) {
    struct GetDevices_return rc;
    uint32_t *list = NULL;        
    int size = 0;

    do {
        size += 16;

        uint32_t *p;        
        if ((p = realloc(list, size * sizeof(uint32_t))) == NULL) {
            free(list);
            set_error("Error allocating storage for slice");
            return -1;            
        } else {
            list = p;
        }

        GoSlice slice = { list,0,size} ;

        rc = GetDevices(u, slice);
        if (rc.r1 != NULL) {
            free(list);
            set_error(rc.r1);
            return -1;
        }

    } while (rc.r0 > size);

    *N = rc.r0;
    *devices = malloc(rc.r0 * sizeof(uint32_t));

    memmove(*devices, list, rc.r0 * sizeof(uint32_t));
    free(list);

    return 0;
}

int get_device(unsigned id, struct device *d) {
    struct GetDevice_return rc = GetDevice(u,id);

    if (rc.r1 != NULL) {        
        set_error(rc.r1);
        return -1;
    }

    d->ID = rc.r0.ID;
    snprintf(d->address, sizeof(d->address), "%s", rc.r0.address);
    snprintf(d->subnet,  sizeof(d->subnet),  "%s", rc.r0.subnet);
    snprintf(d->gateway, sizeof(d->gateway), "%s", rc.r0.gateway);
    snprintf(d->MAC,     sizeof(d->MAC),     "%s", rc.r0.MAC);
    snprintf(d->version, sizeof(d->version), "%s", rc.r0.version);
    snprintf(d->date,    sizeof(d->date),    "%s", rc.r0.date);

    free(rc.r0.address);
    free(rc.r0.subnet);
    free(rc.r0.gateway);
    free(rc.r0.MAC);
    free(rc.r0.version);
    free(rc.r0.date);

    return 0;
}