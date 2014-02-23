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
#include "symtable.h"

static SymTable_Entry *symtable_entry_init(void);
static unsigned int symtable_hash(char *);

SymTable *symtable_init(void)
{
    SymTable *t;

    /* initialize symbol table */
    if ((t = (SymTable *)calloc(1, sizeof(SymTable))) == NULL) {
        perror("Allocate symbol table failed");
        exit(EXIT_FAILURE);
    }
    t->hash = symtable_hash;

    assert(t->stack == NULL);
    assert(t->entries == NULL);

    if ((t->entries = (SymTable_Entry **)calloc(SYMTAB_SIZE, sizeof(SymTable_Entry))) == NULL) {
        perror("Allocate symbol table entry storage failed");
        exit(EXIT_FAILURE);
    }

    return t;
}

void symtable_enter_scope(SymTable *t)
{
    SymTable_Stack *s;

    /* enter current scope to stack */
    if ((s = (SymTable_Stack *)calloc(1, sizeof(SymTable_Stack))) == NULL) {
        perror("Allocate symbol table stack failed");
        exit(EXIT_FAILURE);
    }

    assert(s->entries == NULL);
    assert(s->next == NULL);

    s->entries = t->entries;
    s->next = t->stack;
    t->stack = s;

    /* init new scope */
    if ((t->entries = (SymTable_Entry **)calloc(SYMTAB_SIZE, sizeof(SymTable_Entry))) == NULL) {
        perror("Allocate symbol table entry failed");
        exit(EXIT_FAILURE);
    }
}

void symtable_leave_scope(SymTable *t)
{
    SymTable_Stack *s;

    /* restore scope from stack */
    s = t->stack;
    t->entries = s->entries;

    /* discard top of stack */
    t->stack = s->next;
}

void symtable_insert(SymTable *t, char *key, SymTable_EntryType type, void *value)
{
    unsigned int i;
    SymTable_Entry *e;

    /* create new symbol entry and enter to front of entry list */
    i = t->hash(key) % SYMTAB_SIZE;
    e = symtable_entry_init();

    assert(e->key == NULL);
    assert(e->value == NULL);
    assert(e->next == NULL);

    e->key = key;
    e->type = type;
    e->value = value;
    e->next = t->entries[i];

    t->entries[i] = e;
}

void *symtable_lookup(SymTable *t, char *key)
{
    unsigned int i;
    SymTable_Entry *e;

    i = t->hash(key) % SYMTAB_SIZE;
    for (e = t->entries[i]; e; e = e->next) {
        if (strcmp(e->key, key) == 0) {
            return e->value;
        }
    }
    return NULL;
}

void symtable_delete(SymTable *t, char *key)
{
    unsigned int i;
    SymTable_Entry *e, *tmp;

    i = t->hash(key) % SYMTAB_SIZE;
    tmp = NULL;
    for (e = t->entries[i]; e; tmp = e, e = e->next) {
        if (strcmp(e->key, key) == 0) {
            if (tmp == NULL) {
                t->entries[i] = e->next;
            }
            else {
                tmp->next = e->next;
            }
            free(e->key);
            free(e->value);
            free(e);
            return;
        }
    }
}

static SymTable_Entry *symtable_entry_init(void)
{
    SymTable_Entry *e;

    /* initialize new symbol entry */
    if ((e = (SymTable_Entry *)calloc(1, sizeof(SymTable_Entry))) == NULL) {
        perror("Allocate symbol table entry failed");
        exit(EXIT_FAILURE);
    }

    return e;
}

static unsigned int symtable_hash(char *s)
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
