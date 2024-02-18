#ifndef __EVENT_H__
#define __EVENT_H__

typedef struct _event_t {
    int  pid;
    char comm[16];
    char username[80];
    char password[80];
} event_t;

typedef struct pam_handle
{
  char *authtok;
  unsigned caller_is;
  void *pam_conversation;
  char *oldauthtok;
  char *prompt;
  char *service_name;
  char *user;
  char *rhost;
  char *ruser;
  char *tty;
  char *xdisplay;
  char *authtok_type;
  void *data;
  void *env;
} pam_handle_t;

#endif
