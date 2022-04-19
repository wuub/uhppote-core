#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "tests.h"
#include "uhppoted.h"

void usage();
bool all();

typedef bool (*f)();

typedef struct test {
    const char *name;
    f fn;
} test;

// typedef struct result {
//     const char *field;
//     const int type;
//     union {
//         bool expected;
//         bool value;
//     } Boolean;
// } result;

const uint32_t DEVICE_ID = 405419896;
const uint32_t CARD_NUMBER = 8165538;
const uint32_t CARD_INDEX = 19;
const uint32_t EVENT_INDEX = 51;
const uint8_t DOOR = 4;

const test tests[] = {
    {.name = "get-devices", .fn = getDevices},
    {.name = "get-device", .fn = getDevice},
    {.name = "set-address", .fn = setAddress},
    {.name = "get-status", .fn = getStatus},
    {.name = "get-time", .fn = getTime},
    {.name = "set-time", .fn = setTime},
    {.name = "get-listener", .fn = getListener},
    {.name = "set-listener", .fn = setListener},
    {.name = "get-door-control", .fn = getDoorControl},
    {.name = "set-door-control", .fn = setDoorControl},
    {.name = "open-door", .fn = openDoor},
    {.name = "get-cards", .fn = getCards},
    {.name = "get-card", .fn = getCard},
    {.name = "get-card-by-index", .fn = getCardByIndex},
    {.name = "put-card", .fn = putCard},
    {.name = "delete-card", .fn = deleteCard},
    {.name = "delete-cards", .fn = deleteCards},
    {.name = "get-event-index", .fn = getEventIndex},
    {.name = "set-event-index", .fn = setEventIndex},
    {.name = "get-event", .fn = getEvent},
    {.name = "record-special-events", .fn = recordSpecialEvents},
};

controller alpha = {.id = 405419896, .address = "192.168.1.100"};
controller beta = {.id = 303986753, .address = "192.168.1.100"};

int main(int argc, char **argv) {
    bool ok = true;
    char *cmd;

    if (argc > 1) {
        cmd = argv[1];
    }

    setup("192.168.1.100:0", "192.168.1.255:60000", "192.168.1.100:60001", 2500, true, &alpha, &beta, NULL);

    if (cmd == NULL || strncmp(cmd, "all", 3) == 0) {
        ok = all();
    } else if (strcmp(cmd, "help") == 0) {
        printf("\n");
        usage();
    } else {
        int N = sizeof(tests) / sizeof(test);

        for (int i = 0; i < N; i++) {
            test t = tests[i];
            if (strncmp(cmd, t.name, strlen(t.name)) == 0) {
                ok = t.fn();
                goto done; // <evil cackle> always wanted to do this just to annoy somebody on the Internet
            }
        }

        printf("\n*** ERROR invalid command (%s)\n\n", cmd);
        usage();
        ok = false;
    }

done:
    teardown();

    return ok ? 0 : -1;
}

bool all() {
    bool ok = true;
    int N = sizeof(tests) / sizeof(test);

    for (int i = 0; i < N; i++) {
        test t = tests[i];
        ok = t.fn() ? ok : false;
    }

    return ok;
}

void usage() {
    int N = sizeof(tests) / sizeof(test);

    printf("   Usage: test <command>\n");
    printf("\n");
    printf("   Supported commands:\n");

    for (int i = 0; i < N; i++) {
        test t = tests[i];
        printf("      %s\n", t.name);
    }

    printf("\n");
}

// bool evaluate(const char *tag, result resultset[]) {
// }

bool passed(const char *tag) {
    printf("%-21s ok\n", tag);

    return true;
}

bool failed(const char *tag) {
    printf("%-21s failed\n", tag);

    return false;
}
