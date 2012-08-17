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
#include "vm.h"

int *regs[NUM_REGS];
int stack[SIZE_STACK];
int sp;

void op_noop(void);
void op_move(void);
void op_xchg(void);
/*
void op_var(void);
void op_load(void);
void op_stor(void);
*/
void op_push(void);
void op_pop(void);

void op_add(void);
void op_sub(void);
void op_mul(void);
void op_div(void);
void op_neg(void);
/*
void op_ccat(void);
*/

void op_and(void);
void op_or(void);
void op_not(void);
/*
void op_cmp(void);
void op_jmp(void);
*/

int getmnem(void);
int getreg(void);
int getint(void);

int getmnem() {
    int m;
    if (scanf("%d", &m) != 1) {
        perror("scanf");
        exit(EXIT_FAILURE);
    }
    // return -1 on unknown
    return (m >= 0 && m < NUM_OPS) ? m : -1;
}

int getreg() {
    int r;
    if (scanf("%d", &r) != 1) {
        perror("scanf");
        exit(EXIT_FAILURE);
    }
    // return -1 on invalid register
    return (r >= 0 && r < NUM_REGS) ? r : -1;
}

int getint() {
    int i;
    if (scanf("%d", &i) != 1) {
        perror("scanf");
        exit(EXIT_FAILURE);
    }
    return i;
}


void op_noop() {
    return;
}

void op_move() {
    int dest = getreg();
    int src = getint();
    *regs[dest] = src;
}

void op_xchg() {
    int dest = getreg();
    int src = getreg();
    int tmp1 = *regs[dest];
    int tmp2 = *regs[src];
    *regs[dest] = tmp2;
    *regs[src] = tmp1;
}

/*
void op_var();
void op_load();
void op_stor();
*/
void op_push() {
    int src = getreg();

    if (sp != SIZE_STACK - 1) {
        sp++;
        stack[sp] = *regs[src];
    }
    else {
        fprintf(stderr, "PUSH: no room on stack");
        exit(EXIT_FAILURE);
    }
}

void op_pop() {
    int src = getreg();

    if (sp != -1) {
        *regs[src] = stack[sp];
        sp--;
    }
    else {
        fprintf(stderr, "POP: empty stack");
        exit(EXIT_FAILURE);
    }
}

void op_add() {
    int dest = getreg();
    int src = getreg();
    *regs[dest] += *regs[src];
}

void op_sub() {
    int dest = getreg();
    int src = getreg();
    *regs[dest] -= *regs[src];
}

void op_mul() {
    int dest = getreg();
    int src = getreg();
    *regs[dest] *= *regs[src];
}

void op_div() {
    int dest = getreg();
    int src = getreg();
    *regs[dest] /= *regs[src];
}

void op_neg() {
    int dest = getreg();
    *regs[dest] = -*regs[dest];
}

/*
void op_ccat();
*/
void op_and() {
    int dest = getreg();
    int src = getreg();
    *regs[dest] = *regs[dest] && *regs[src];
}

void op_or() {
    int dest = getreg();
    int src = getreg();
    *regs[dest] = *regs[dest] || *regs[src];
}

void op_not() {
    int dest = getreg();
    *regs[dest] = !*regs[dest];
}
/*
void op_cmp();
void op_jmp();
*/

int main() {

    int m, i;

    sp = -1;

    void (*ops[NUM_OPS])();
    for (i = 0; i < NUM_OPS; i++) {
        ops[i] = 0;
    }

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

    while ((m = getmnem()) > -1) {
        if (ops[m]) {
            ops[m]();
            printf(":%d %d %d\n", *regs[0], *regs[1], *regs[2]);
        }
        else {
            printf("not implemented");
        }
    }

    for (i = 0; i < NUM_REGS; i++) {
        free(regs[i]);
    }
}

