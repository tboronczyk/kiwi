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

#include <assert.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <assert.h>
#include "symtab.h"

static unsigned int symtab_hash(char *s)
{
    /* TODO: need a good hash function */
    char *c;
    unsigned int hash;

    hash = 0;
    for (c = s; *c != '\0'; c++) {
        hash += (unsigned int)*c;
    }
    return hash;
}

static symtab_entry_t *symtab_entry_init(void)
{
    symtab_entry_t *e;

    /* initialize new symbol entry */
    if ((e = (symtab_entry_t *)calloc(1, sizeof(symtab_entry_t))) == NULL) {
        perror("Allocate symbol table entry failed");
        exit(EXIT_FAILURE);
    }

    return e;
}

symtab_t *symtab_init()
{
    symtab_t *t;

    /* initialize symbol table */
    if ((t = (symtab_t *)calloc(1, sizeof(symtab_t))) == NULL) {
        perror("Allocate symbol table failed");
        exit(EXIT_FAILURE);
    }
    t->hash = symtab_hash;

    assert(t->stack == NULL);
    assert(t->entries == NULL);

    if ((t->entries = (symtab_entry_t **)calloc(SYMTAB_SIZE, sizeof(symtab_entry_t))) == NULL) {
        perror("Allocate symbol table entry storage failed");
        exit(EXIT_FAILURE);
    }

    return t;
}

void symtab_enterscope(symtab_t *t)
{
    symtab_stack_t *s;

    /* enter current scope to stack */
    if ((s = (symtab_stack_t *)calloc(1, sizeof(symtab_stack_t))) == NULL) {
        perror("Allocate symbol table stack failed");
        exit(EXIT_FAILURE);
    }

    assert(s->entries == NULL);
    assert(s->next == NULL);

    s->entries = t->entries;
    s->next = t->stack;
    t->stack = s;

    /* init new scope */
    if ((t->entries = (symtab_entry_t **)calloc(SYMTAB_SIZE, sizeof(symtab_entry_t))) == NULL) {
        perror("Allocate symbol table entry failed");
        exit(EXIT_FAILURE);
    }
}

void symtab_leavescope(symtab_t *t)
{
    symtab_stack_t *s;

    /* restore scope from stack */
    s = t->stack;
    t->entries = s->entries;

    /* discard top of stack */
    t->stack = s->next;
}

void symtab_insert(symtab_t *t, char *key, symtab_entrytype_t type, void *value)
{
    unsigned int i;
    symtab_entry_t *e;

    /* create new symbol entry and enter to front of entry list */
    i = t->hash(key) % SYMTAB_SIZE;
    e = symtab_entry_init();

    assert(e->key == NULL);
    assert(e->value == NULL);
    assert(e->next == NULL);

    e->key = key;
    e->type = type;
    e->value = value;
    e->next = t->entries[i];

    t->entries[i] = e;
}

void *symtab_lookup(symtab_t *t, char *key)
{
    unsigned int i;
    symtab_entry_t *e;

    i = t->hash(key) % SYMTAB_SIZE;
    for (e = t->entries[i]; e; e = e->next) {
        if (strcmp(e->key, key) == 0) {
            return e->value;
        }
    }
    return NULL;
}

void symtab_delete(symtab_t *t, char *key)
{
    unsigned int i;
    symtab_entry_t *e;

    /* remove first occurence from entry list */
    i = t->hash(key) % SYMTAB_SIZE;
    e = t->entries[i];
    t->entries[i] = (t->entries[i])->next;

    /* free(e); */
}
