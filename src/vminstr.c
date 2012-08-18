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
#include <stdarg.h>
#include "opcodes.h"
#include "vminstr.h"

VM_Instr *vminstr_init(int op, ...) {
    VM_Instr *instr;
    va_list ap;
    if ((instr = (VM_Instr *)calloc(1, sizeof(VM_Instr))) == NULL) {
        perror("calloc");
        exit(EXIT_FAILURE);
    }

    instr->op = op;

    va_start(ap, op);
    if ((SINGLE_OPS & op) == op) {
        instr->dest = va_arg(ap, int);
    }
    else if ((DOUBLE_OPS & op) == op) {
        instr->dest = va_arg(ap, int);
        instr->src  = va_arg(ap, int);
    }
    va_end(ap);

    return instr;
}

void vminstr_free(VM_Instr *instr) {
    free(instr);
}

void vminstr_exec(VM_Machine *vm, VM_Instr *i) {
    int op = i->op;
    vm->ops[op](vm, i);
}
