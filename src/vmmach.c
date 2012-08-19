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
#include "opcodes.h"
#include "vmprogbuf.h"
#include "vmop.h"

VM_Machine *vmmach_init(void)
{
    VM_Machine *vm;
    int i;

    /* allocate memory for machine */
    if ((vm = (VM_Machine *)calloc(1, sizeof(VM_Machine))) == NULL) {
        perror("Allocate machine failured");
        exit(EXIT_FAILURE);
    }

    /* initialize machine */
    vm->sp = -1;
    vm->ip = 0;
    for (i = 0; i < VMMACH_NUM_REGS; i++) {
        vm->regs[i] = (int *)calloc(1, sizeof(int));
    }

    /* bind machine functions to operators */
    vm->ops[OP_NOOP] = vmop_noop;
    vm->ops[OP_MOVE] = vmop_move;
    vm->ops[OP_XCHG] = vmop_xchg;
/*
    vm->ops[OP_VAR] = vmop_var;
    vm->ops[OP_LOAD] = vmop_load;
    vm->ops[OP_STOR] = vmop_stor;
*/
    vm->ops[OP_PUSH] = vmop_push;
    vm->ops[OP_POP] = vmop_pop;
    vm->ops[OP_ADD] = vmop_add;
    vm->ops[OP_SUB] = vmop_sub;
    vm->ops[OP_MUL] = vmop_mul;
    vm->ops[OP_DIV] = vmop_div;
    vm->ops[OP_NEG] = vmop_neg;
/*
    vm->ops[OP_CCAT] = vmop_ccat;
*/
    vm->ops[OP_AND] = vmop_and;
    vm->ops[OP_OR] = vmop_or;
    vm->ops[OP_NOT] = vmop_not;
/*
    vm->ops[OP_CMP] = vmop_cmp;
    vm->ops[OP_JMP] = vmop_jmp;
*/
    return vm;
}

void vmmach_free(VM_Machine *vm)
{
    int i;
    for (i = 0; i < VMMACH_NUM_REGS; i++) {
        free(vm->regs[i]);
    }
    free(vm);
}
