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
#include "vm.h"

#define UNUSED(x) (void)(x)

#define VMPROGBUF_SIZE_INIT 5 
#define VMPROGBUF_SIZE_INCR 5

void vmop_noop(VM_Machine *, VM_Instr *);
void vmop_move(VM_Machine *, VM_Instr *);
void vmop_xchg(VM_Machine *, VM_Instr *);
/*
void vmop_var(VM_Machine *, VM_Instr *);
void vmop_load(VM_Machine *, VM_Instr *);
void vmop_stor(VM_Machine *, VM_Instr *);
*/
void vmop_push(VM_Machine *, VM_Instr *);
void vmop_pop(VM_Machine *, VM_Instr *);

void vmop_add(VM_Machine *, VM_Instr *);
void vmop_sub(VM_Machine *, VM_Instr *);
void vmop_mul(VM_Machine *, VM_Instr *);
void vmop_div(VM_Machine *, VM_Instr *);
void vmop_neg(VM_Machine *, VM_Instr *);
/*
void vmop_ccat(VM_Machine *, VM_Instr *);
*/

void vmop_and(VM_Machine *, VM_Instr *);
void vmop_or(VM_Machine *, VM_Instr *);
void vmop_not(VM_Machine *, VM_Instr *);
/*
void vmop_cmp(VM_Machine *, VM_Instr *);
void vmop_jmp(VM_Machine *, VM_Instr *);
*/

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

VM_Instr *vminstr_init(int, ...);
void vminstr_free(VM_Instr *);

void vminstr_exec(VM_Machine *, VM_Instr *);

VM_Instr *vminstr_init(int op, ...)
{
    VM_Instr *instr;
    va_list ap;

    /* allocate memory for instruction */
    if ((instr = (VM_Instr *)calloc(1, sizeof(VM_Instr))) == NULL) {
        perror("Allocate instruction failed");
        exit(EXIT_FAILURE);
    }

    /* initalize instruction */
    instr->op = op;
    va_start(ap, op);
    if (op_numargs[op] > 0) { instr->dest = va_arg(ap, int); }
    if (op_numargs[op] > 1) { instr->src  = va_arg(ap, int); }
    va_end(ap);

    return instr;
}

void vminstr_free(VM_Instr *instr)
{
    free(instr);
}

void vminstr_exec(VM_Machine *vm, VM_Instr *i)
{
    /* execute instruction by calling appropriate machine function */
    int op = i->op;
    vm->ops[op](vm, i);
}


VM_ProgBuf *vmprogbuf_init(void);
void vmprogbuf_free(VM_ProgBuf *);
void vmprogbuf_grow(VM_ProgBuf *);
void vmprogbuf_push(VM_ProgBuf *, VM_Instr *);
void vmprogbuf_exec(VM_ProgBuf *, VM_Machine *);

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

void vmprogbuf_grow(VM_ProgBuf *b)
{
    /* increase storage capacity of buffer */
    b->len += VMPROGBUF_SIZE_INCR;
    if ((b->instr = (VM_Instr **)realloc(b->instr, sizeof(VM_Instr *) * b->len)) == NULL) {
        perror("Reallocate program buffer instruction storage failed");
        exit(EXIT_FAILURE);
    }
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

int main()
{
    VM_ProgBuf *b = vmprogbuf_init();
    VM_Machine *vm = vmmach_init();

    /* load program */
    vmprogbuf_push(b, vminstr_init(OP_NOOP));
    vmprogbuf_push(b, vminstr_init(OP_MOVE, 0, 10));
    vmprogbuf_push(b, vminstr_init(OP_MOVE, 1, 1));
    vmprogbuf_push(b, vminstr_init(OP_SUB, 0, 1));

    /* execute program */
    vmprogbuf_exec(b, vm);

    vmprogbuf_free(b);
    vmmach_free(vm);

    return EXIT_SUCCESS;
}

