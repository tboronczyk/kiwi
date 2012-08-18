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

#ifndef OPCODES_H
#define OPCODES_H

#define NUM_OPS 19

#define OP_NOOP	1
#define OP_MOVE (1 << 1)	
#define OP_XCHG	(1 << 2)
#define OP_VAR  (1 << 3)	
#define OP_LOAD (1 << 4)
#define OP_STOR	(1 << 5)
#define OP_PUSH	(1 << 6)
#define OP_POP  (1 << 7)
#define OP_ADD  (1 << 8)
#define OP_SUB	(1 << 9)
#define OP_MUL	(1 << 10)
#define OP_DIV	(1 << 11)
#define OP_NEG	(1 << 12)
#define OP_CCAT	(1 << 13)
#define OP_AND	(1 << 14)
#define OP_OR	(1 << 15)
#define OP_NOT	(1 << 16)
#define OP_CMP	(1 << 17)
#define OP_JMP	(1 << 18)

#define NONE_OPS   (OP_NOOP)
#define SINGLE_OPS (OP_POP | OP_NEG | OP_NOT)
#define DOUBLE_OPS (OP_MOVE | OP_XCHG | OP_PUSH | OP_ADD | OP_SUB | OP_MUL \
                   | OP_MUL | OP_DIV | OP_AND | OP_OR)
#endif

