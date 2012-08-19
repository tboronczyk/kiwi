#ifndef VM_H
#define VM_H

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

#define NUM_OPS 19

#define OP_NOOP 0
#define OP_MOVE 1
#define OP_XCHG 2
#define OP_VAR  3
#define OP_LOAD 4
#define OP_STOR 5
#define OP_PUSH 6
#define OP_POP  7
#define OP_ADD  8
#define OP_SUB  9
#define OP_MUL  10
#define OP_DIV  11
#define OP_NEG  12
#define OP_CCAT 13
#define OP_AND  14
#define OP_OR   15
#define OP_NOT  16
#define OP_CMP  17
#define OP_JMP  18

int op_numargs[NUM_OPS] = {
    /*OP_NOOP*/   0,
    /* OP_MOVE */ 2,
    /* OP_XCHG */ 2,
    /* OP_VAR  */ 1,
    /* OP_LOAD */ 2,
    /* OP_STOR */ 2,
    /* OP_PUSH */ 2,
    /* OP_POP  */ 1,
    /* OP_ADD  */ 2,
    /* OP_SUB  */ 2,
    /* OP_MUL  */ 2,
    /* OP_DIV  */ 2,
    /* OP_NEG  */ 1,
    /* OP_CCAT */ 2,
    /* OP_AND  */ 2,
    /* OP_OR   */ 2,
    /* OP_NOT  */ 1,
    /* OP_CMP  */ 2,
    /* OP_JMP  */ 1
};

#define VMMACH_NUM_REGS 3
#define VMMACH_SIZE_STACK 80

typedef struct _VM_Instr
{
    int op;
    int dest;
    int src;
}
VM_Instr;

typedef struct _VM_ProgBuf
{
    int len, tail;
    VM_Instr **instr;
}
VM_ProgBuf;

typedef struct _VM_Machine
{
    int sp,                             /* stack pointer */
        ip,                             /* instruction pointer */
        *regs[VMMACH_NUM_REGS],         /* registers */
        stack[VMMACH_SIZE_STACK];       /* stack */
    void (*ops[1 << (NUM_OPS - 1)])();  /* operation functions */
}
VM_Machine;

VM_Machine *vmmach_init(void);
/*
void vmmach_load(VM_Machine *);
void vmmach_exec(VM_Machine *);
*/
void vmmach_free(VM_Machine *);

#endif

