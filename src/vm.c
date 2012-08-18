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
#include "vm.h"

#define UNUSED(x) (void)(x)

int *regs[NUM_REGS];
int stack[SIZE_STACK];
int sp;
void (*ops[1 << (NUM_OPS - 1)])();

Instruction *instr_init(int op, ...) {
    Instruction *instr;
    va_list ap;
    if ((instr = (Instruction *)calloc(1, sizeof(Instruction))) == NULL) {
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

void instr_free(Instruction *instr) {
    free(instr);
}

void instr_exec(ProgBuf *b, int i) {
    int op = b->instr[i]->op;
    ops[op](b->instr[i]);
}
    
ProgBuf *progbuf_init(void) {
    ProgBuf *b;
    if ((b = (ProgBuf *)calloc(1, sizeof(ProgBuf))) == NULL) {
        perror("calloc");
        exit(EXIT_FAILURE);
    }

    b->tail = 0;
    b->len = PROGBUF_SIZE_INIT;
    if ((b->instr = (Instruction **)calloc(b->len, sizeof(Instruction *))) == NULL) {
        perror("calloc");
        exit(EXIT_FAILURE);
    }

    return b;
}

void progbuf_free(ProgBuf *b) {
    int i;
    for (i = 0; i < b->tail; i++) {
        instr_free(b->instr[i]);
    }
    free(b->instr);
}

static void progbuf_grow(ProgBuf *b) {
    // increase storage capacity of buffer
    b->len += PROGBUF_SIZE_INCR;
    if ((b->instr = (Instruction **)realloc(b->instr, sizeof(Instruction *) * b->len)) == NULL) {
        perror("realloc");
        exit(EXIT_FAILURE);
    }
}

static void progbuf_push(ProgBuf *b, Instruction *i) {
    b->instr[b->tail] = i;
    b->tail++;
    // increase buffer size if necessary
    if (b->tail == b->len) {
        progbuf_grow(b);
    }
}

void op_noop(Instruction *instr) {
    UNUSED(instr);
    return;
}

void op_move(Instruction *instr) {
    int dest = instr->dest;
    int src = instr->src;
    *regs[dest] = src;
}

void op_xchg(Instruction *instr) {
    int dest = instr->dest;
    int src = instr->src;
    int tmp1 = *regs[dest];
    int tmp2 = *regs[src];
    *regs[dest] = tmp2;
    *regs[src] = tmp1;
}
/*
void op_var(Instruction *instr) {
}

void op_load(Instruction *instr) {
}

void op_stor(Instruction *instr) {
}
*/
void op_push(Instruction *instr) {
    int dest = instr->dest;

    if (sp != SIZE_STACK - 1) {
        sp++;
        stack[sp] = *regs[dest];
    }
    else {
        fprintf(stderr, "PUSH: no room on stack");
        exit(EXIT_FAILURE);
    }
}

void op_pop(Instruction *instr) {
    int dest = instr->dest;

    if (sp != -1) {
        *regs[dest] = stack[sp];
        sp--;
    }
    else {
        fprintf(stderr, "POP: empty stack");
        exit(EXIT_FAILURE);
    }
}

void op_add(Instruction *instr) {
    int dest = instr->dest;
    int src = instr->src;
    *regs[dest] += *regs[src];
}
void op_sub(Instruction *instr) {
    int dest = instr->dest;
    int src = instr->src;
    *regs[dest] -= *regs[src];
}
void op_mul(Instruction *instr) {
    int dest = instr->dest;
    int src = instr->src;
    *regs[dest] *= *regs[src];
}
void op_div(Instruction *instr) {
    int dest = instr->dest;
    int src = instr->src;
    *regs[dest] /= *regs[src];
}

void op_neg(Instruction *instr) {
    int dest = instr->dest;
    *regs[dest] = -*regs[dest];
}
/*
void op_ccat(Instruction *instr) {
}
*/
void op_and(Instruction *instr) {
    int dest = instr->dest;
    int src = instr->src;
    *regs[dest] = *regs[dest] && *regs[src];
}

void op_or(Instruction *instr) {
    int dest = instr->dest;
    int src = instr->src;
    *regs[dest] = *regs[dest] || *regs[src];
}

void op_not(Instruction *instr) {
    int dest = instr->dest;
    *regs[dest] = !*regs[dest];
}
/*
void op_cmp(Instruction *instr) {
}

void op_jmp(Instruction *instr) {
}
*/

int main() {

    int i;
    sp = -1;
    ProgBuf *b = progbuf_init();

    ops[OP_NOOP] = op_noop;
    ops[OP_MOVE] = op_move;
    ops[OP_XCHG] = op_xchg;
/*
    ops[OP_VAR] = op_var;
    ops[OP_LOAD] = op_load;
    ops[OP_STOR] = op_stor;
*/
    ops[OP_PUSH] = op_push;
    ops[OP_POP] = op_pop;
    ops[OP_ADD] = op_add;
    ops[OP_SUB] = op_sub;
    ops[OP_MUL] = op_mul;
    ops[OP_DIV] = op_div;
    ops[OP_NEG] = op_neg;
/*
    ops[OP_CCAT] = op_ccat;
*/
    ops[OP_AND] = op_and;
    ops[OP_OR] = op_or;
    ops[OP_NOT] = op_not;
/*
    ops[OP_CMP] = op_cmp;
    ops[OP_JMP] = op_jmp;
*/

    for (i = 0; i < NUM_REGS; i++) {
        regs[i] = (int *)calloc(1, sizeof(int));
    }

    progbuf_push(b, instr_init(OP_NOOP));
    progbuf_push(b, instr_init(OP_MOVE, 0, 1));
    progbuf_push(b, instr_init(OP_MOVE, 1, 1));
    progbuf_push(b, instr_init(OP_ADD, 0, 1));

    for (i = 0; i < b->tail; i++) {
        instr_exec(b, i);
        printf(":%d %d %d\n", *regs[0], *regs[1], *regs[2]);
    }

    progbuf_free(b);

    for (i = 0; i < NUM_REGS; i++) {
        free(regs[i]);
    }
}

