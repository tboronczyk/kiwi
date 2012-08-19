/*
 * Copyright (c) 2012, Timothy Boronczyk
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  1. Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *
 *  2. Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *
 *  3. The names of the authors may not be used to endorse or promote products
 *     derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED "AS IS" AND WITHOUT ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, WITHOUT LIMITATION, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.
 */

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "mysocket.h"

#define DEFAULT_PORT 11040
#define BUFSIZE 32

void handle_client(Socket sock);
void run_server(Socket sock);

void handle_client(Socket sock)
{
    char buf[BUFSIZE];
    int len;

    do {
        if ((len = recv(sock, buf, BUFSIZE, 0)) < 0) {
            perror("recv");
            exit(EXIT_FAILURE);
        }

        /* Just send it back for now... */
        if (send(sock, buf, len, 0) != len) {
            perror("send");
            exit(EXIT_FAILURE);
        }
    }
    while (len);

    close(sock);
}

void run_server(Socket sock)
{
    Socket clientsock;

    while (1) {
        if ((clientsock = accept(sock, NULL, NULL)) < 0) {
             perror("accept");
             exit(EXIT_FAILURE);
        }
        handle_client(clientsock);
    }
}

int main(int argc, char *argv[])
{
    unsigned short port = DEFAULT_PORT;

    if (argc > 2) {
       fprintf(stderr, "%s [port]\n", argv[0]);
       exit(EXIT_FAILURE);
    }
    else if (argc == 2) {
        port = atoi(argv[1]);
    }

    Socket sock;
    SocketAddress *addr;
    int limit = 0;

    addr = init_socketaddress(INADDR_ANY, port);
    sock = init_socket(addr, limit);
    run_server(sock);
    /* free_socketaddress(addr); */

    return EXIT_SUCCESS;
}
