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
#include "vmop.h"

#define UNUSED(x) (void)(x)

void vmop_noop(VM_Machine *vm, VM_Instr *instr)
{
    UNUSED(vm);
    UNUSED(instr);
    return;
}

void vmop_move(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    int src = instr->src;
    *vm->regs[dest] = src;
}

void vmop_xchg(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    int src = instr->src;
    int tmp1 = *vm->regs[dest];
    int tmp2 = *vm->regs[src];
    *vm->regs[dest] = tmp2;
    *vm->regs[src] = tmp1;
}
/*
void vmop_var(VM_Machine *vm, VM_Instr *instr)
{
}

void vmop_load(VM_Machine *vm, VM_Instr *instr)
{
}

void vmop_stor(VM_Machine *vm, VM_Instr *instr)
{
}
*/
void vmop_push(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;

    if (vm->sp != VMMACH_SIZE_STACK - 1) {
        vm->sp++;
        vm->stack[vm->sp] = *vm->regs[dest];
    }
    else {
        perror("Stack push failed");
        exit(EXIT_FAILURE);
    }
}

void vmop_pop(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;

    if (vm->sp != -1) {
        *vm->regs[dest] = vm->stack[vm->sp];
        vm->sp--;
    }
    else {
        perror("Stack pop failed");
        exit(EXIT_FAILURE);
    }
}

void vmop_add(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    int src = instr->src;
    *vm->regs[dest] += *vm->regs[src];
} 

void vmop_sub(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    int src = instr->src;
    *vm->regs[dest] -= *vm->regs[src];
}

void vmop_mul(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    int src = instr->src;
    *vm->regs[dest] *= *vm->regs[src];
} 

void vmop_div(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    int src = instr->src;
    *vm->regs[dest] /= *vm->regs[src];
}

void vmop_neg(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    *vm->regs[dest] = -*vm->regs[dest];
}
/*
void vmop_ccat(VM_Machine *vm, VM_Instr *instr)
{
}
*/
void vmop_and(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    int src = instr->src;
    *vm->regs[dest] = *vm->regs[dest] && *vm->regs[src];
}

void vmop_or(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    int src = instr->src;
    *vm->regs[dest] = *vm->regs[dest] || *vm->regs[src];
}

void vmop_not(VM_Machine *vm, VM_Instr *instr)
{
    int dest = instr->dest;
    *vm->regs[dest] = !*vm->regs[dest];
}
/*
void vmop_cmp(VM_Machine *vm, VM_Instr *instr)
{
}

void vmop_jmp(VM_Machine *vm, VM_Instr *instr)
{
}
*/

