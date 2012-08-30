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

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "symtab.h"

int main()
{
    symtab_t *t;
    t = symtab_init();

    int a =40;
    int b =41;

    symtab_insert(t, "a", SYMTAB_ENTRY_NUMBER, &a);
    symtab_enterscope(t);
    symtab_insert(t, "b", SYMTAB_ENTRY_NUMBER, &b);
    symtab_enterscope(t);

    if (symtab_lookup(t, "a")) printf("a=%d\n", *(int *)symtab_lookup(t, "a")); else printf("a not found\n");
    if (symtab_lookup(t, "b")) printf("b=%d\n", *(int *)symtab_lookup(t, "b")); else printf("b not found\n");
    symtab_leavescope(t);
    if (symtab_lookup(t, "a")) printf("a=%d\n", *(int *)symtab_lookup(t, "a")); else printf("a not found\n");
    if (symtab_lookup(t, "b")) printf("b=%d\n", *(int *)symtab_lookup(t, "b")); else printf("b not found\n");
    symtab_leavescope(t);
    if (symtab_lookup(t, "a")) printf("a=%d\n", *(int *)symtab_lookup(t, "a")); else printf("a not found\n");
    if (symtab_lookup(t, "b")) printf("b=%d\n", *(int *)symtab_lookup(t, "b")); else printf("b not found\n");

    return EXIT_SUCCESS;
}

