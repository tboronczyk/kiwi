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

int getmnem() {
    int m;
    scanf("%d", &m);
    return (m >= 0 && m < NUM_OPS) ? m : -1;
}

int getreg() {
    int r;
    scanf("%d", &r);
    return (r >= 0 && r < NUM_REGS) ? r : -1;
}

int getint() {
    int i;
    scanf("%d", &i);
    return i;
}

int *regs[NUM_REGS];

void op_noop() {
    return;
}

void op_move() {
    int dest = getreg();
    int orig = getint();
    *regs[dest] = orig;
}

void op_xchg() {
    int dest = getreg();
    int orig = getreg();
    *regs[dest] = *regs[orig];
}

/*
void op_var();
void op_load();
void op_stor();
void op_push();
void op_pop();
*/

void op_add() {
    int dest = getreg();
    int orig = getreg();
    *regs[dest] += *regs[orig];
}


void op_sub() {
    int dest = getreg();
    int orig = getreg();
    *regs[dest] -= *regs[orig];
}

void op_mul() {
    int dest = getreg();
    int orig = getreg();
    *regs[dest] *= *regs[orig];
}

void op_div() {
    int dest = getreg();
    int orig = getreg();
    *regs[dest] /= *regs[orig];
}

void op_neg() {
    int dest = getreg();
    *regs[dest] = -*regs[dest];
}

/*
void op_ccat();
void op_and();
void op_or();
void op_cmp();
void op_jmp();
*/

int main() {

    int m, i;

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
    ops[OP_PUSH] = op_push;
    ops[OP_POP] = op_pop;
*/
    ops[OP_ADD] = op_add;
    ops[OP_SUB] = op_sub;
    ops[OP_MUL] = op_mul;
    ops[OP_DIV] = op_div;
    ops[OP_NEG] = op_neg;
/*
    ops[OP_CCAT] = op_ccat;
    ops[OP_AND] = op_and;
    ops[OP_OR] = op_or;
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

