#ifndef VMOP_H
#define VMOP_H

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

#include "vmmach.h"
#include "vminstr.h"

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

#endif
