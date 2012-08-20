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

#define VMMACH_NUM_REGS   3
#define VMMACH_SIZE_STACK 80

typedef enum
{
    OP_NOOP,
    OP_MOVE,
    OP_XCHG,
/*
    OP_VAR,
    OP_LOAD,
    OP_STOR,
*/
    OP_PUSH,
    OP_POP,
    OP_ADD,
    OP_SUB,
    OP_MUL,
    OP_DIV,
    OP_NEG,
/*
    OP_CCAT,
*/
    OP_AND,
    OP_OR,
    OP_NOT,
/*
    OP_CMP,
    OP_JUMP
*/
}
OpCode;

typedef struct _VM_Instr
{
    OpCode op;
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
    int sp,                       /* stack pointer */
        ip,                       /* instruction pointer */
        *regs[VMMACH_NUM_REGS],   /* registers */
        stack[VMMACH_SIZE_STACK]; /* stack */
}
VM_Machine;

VM_Machine *vmmach_init(void);
/*
void vmmach_load(VM_Machine *);
*/
void vmmach_exec(VM_Machine *, VM_ProgBuf *);
void vmmach_free(VM_Machine *);

#endif

