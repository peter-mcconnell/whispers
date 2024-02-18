#ifndef __EVENT_H__
#define __EVENT_H__

typedef struct _event_t {
    int  pid;
    char comm[16];
    char username[80];
    char password[80];
} event_t;

#endif
