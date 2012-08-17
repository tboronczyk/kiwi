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

#define OP_NOOP	0

#define OP_MOVE	1
#define OP_XCHG	2
/*
#define OP_VAR	3
#define OP_LOAD	4
#define OP_STOR	5
*/
#define OP_PUSH	6
#define OP_POP	7

#define OP_ADD	8
#define OP_SUB	9
#define OP_MUL	10
#define OP_DIV	11
#define OP_NEG	12
/*
#define OP_CCAT	13
*/

#define OP_AND	14
#define OP_OR	15
#define OP_NOT	16

/*
#define OP_CMP	17
#define OP_JMP	18
*/
#endif

