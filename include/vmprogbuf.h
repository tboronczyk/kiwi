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

#ifndef VMPROGBUF_H
#define VMPROGBUF_H

#include "vminstr.h"
#include "vmmach.h"

#define VMPROGBUF_SIZE_INIT 5 
#define VMPROGBUF_SIZE_INCR 5

typedef struct _VM_ProgBuf{
    int len, tail;
    VM_Instr **instr;
}
VM_ProgBuf;

VM_ProgBuf *vmprogbuf_init(void);
void vmprogbuf_free(VM_ProgBuf *);

void vmprogbuf_push(VM_ProgBuf *, VM_Instr *);
void vmprogbuf_exec(VM_ProgBuf *, VM_Machine *);

#endif
