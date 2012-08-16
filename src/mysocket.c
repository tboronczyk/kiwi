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
#include "mysocket.h"

SocketAddress *init_socketaddress(long int addr, unsigned short port) {
    SocketAddress *sa;

    if ((sa = (SocketAddress *)calloc(1, sizeof(SocketAddress))) == NULL) {
        perror("calloc");
        exit(EXIT_FAILURE);
    }

    sa->sin_family = AF_INET;
    sa->sin_addr.s_addr = htonl(addr);
    sa->sin_port = htons(port);

    return sa;
}

void free_socketaddress(SocketAddress *sa) {
    free(sa);
}

Socket init_socket(SocketAddress *addr, int pending) {
    Socket sock;

    if ((sock = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP)) < 0) {
        perror("socket");
        exit(EXIT_FAILURE);
    }

    if (bind(sock, (struct sockaddr *)addr, sizeof(SocketAddress)) < 0) {
        perror("bind");
        exit(EXIT_FAILURE);
    }

    if (listen(sock, pending) < 0) {
        perror("listen");
        exit(EXIT_FAILURE);
    }

    return sock;
}

