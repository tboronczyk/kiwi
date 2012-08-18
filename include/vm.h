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

#ifndef VM_H
#define VM_H

#define NUM_REGS 3
#define SIZE_STACK 80

#define PROGBUF_SIZE_INIT 5 
#define PROGBUF_SIZE_INCR 5

typedef struct _Instruction {
    int op;
    int dest;
    int src;
}
Instruction;

typedef struct {
    int len, tail;
    Instruction **instr;
}
ProgBuf;

Instruction *instr_init(int, ...);
void instr_exec(ProgBuf *, int);
void instr_free(Instruction *);

ProgBuf *progbuf_init(void);
void progbuf_free(ProgBuf *);

void op_noop(Instruction *);
void op_move(Instruction *);
void op_xchg(Instruction *);
/*
void op_var(Instruction *);
void op_load(Instruction *);
void op_stor(Instruction *);
*/
void op_push(Instruction *);
void op_pop(Instruction *);

void op_add(Instruction *);
void op_sub(Instruction *);
void op_mul(Instruction *);
void op_div(Instruction *);
void op_neg(Instruction *);
/*
void op_ccat(Instruction *);
*/

void op_and(Instruction *);
void op_or(Instruction *);
void op_not(Instruction *);
/*
void op_cmp(Instruction *);
void op_jmp(Instruction *);
*/

#endif

