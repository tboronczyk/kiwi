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
#include "vmprogbuf.h"

static void vmprogbuf_grow(VM_ProgBuf *b);

VM_ProgBuf *vmprogbuf_init(void)
{
    VM_ProgBuf *b;

    /* allocate memory for program buffer */
    if ((b = (VM_ProgBuf *)calloc(1, sizeof(VM_ProgBuf))) == NULL) {
        perror("Allocate program buffer failed");
        exit(EXIT_FAILURE);
    }

    /* initialize buffer */
    b->tail = 0;
    b->len = VMPROGBUF_SIZE_INIT;
    if ((b->instr = (VM_Instr **)calloc(b->len, sizeof(VM_Instr *))) == NULL) {
        perror("Allocate program buffer instruction storage failed");
        exit(EXIT_FAILURE);
    }

    return b;
}

void vmprogbuf_free(VM_ProgBuf *b)
{
    int i;
    for (i = 0; i < b->tail; i++) {
        vminstr_free(b->instr[i]);
    }
    free(b->instr);
    free(b);
}

void vmprogbuf_push(VM_ProgBuf *b, VM_Instr *i)
{
    b->instr[b->tail] = i;
    b->tail++;
    /* increase buffer size if necessary */
    if (b->tail == b->len) {
        vmprogbuf_grow(b);
    }
}

void vmprogbuf_exec(VM_ProgBuf *b, VM_Machine *vm)
{
    for ( ; vm->ip < b->tail; vm->ip++) {
        vminstr_exec(vm, b->instr[vm->ip]);
        printf(":%d %d %d\n", *vm->regs[0], *vm->regs[1], *vm->regs[2]);
    }
}

static void vmprogbuf_grow(VM_ProgBuf *b)
{
    /* increase storage capacity of buffer */
    b->len += VMPROGBUF_SIZE_INCR;
    if ((b->instr = (VM_Instr **)realloc(b->instr, sizeof(VM_Instr *) * b->len)) == NULL) {
        perror("Reallocate program buffer instruction storage failed");
        exit(EXIT_FAILURE);
    }
}

