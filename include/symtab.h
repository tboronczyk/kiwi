#ifndef SYMTAB_H
#define SYMTAB_H

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

#define SYMTAB_SIZE 109

typedef struct _symtab symtab_t;
typedef struct _symtab_entry symtab_entry_t;
typedef struct _symtab_stack symtab_stack_t;

/* possible data types */
typedef enum
{
    SYMTAB_ENTRY_NUMBER,
    SYMTAB_ENTRY_STRING,
    SYMTAB_ENTRY_FUNC
}
symtab_entrytype_t;

symtab_t *symtab_init(void);
void symtab_enter_scope(symtab_t *);
void symtab_leave_scope(symtab_t *);
void symtab_insert(symtab_t *, char *, symtab_entrytype_t, void *);
void *symtab_lookup(symtab_t *, char *);
void symtab_delete(symtab_t *, char *);

/* symbol table hash table */
struct _symtab
{
    unsigned int (*hash)(char *);  /* function used for hashing keys */
    symtab_entry_t **entries;      /* array of entry lists (hash table) */
    symtab_stack_t *stack;         /* array of entry lists (hash table) */
};

/* linked list bucket for symtab_t hash table */
struct _symtab_entry
{
    char *key;
    void *value;
    symtab_entrytype_t type;
    symtab_entry_t *next;
};

/* linked list bucket for symtab_t hash table */
struct _symtab_stack
{
    symtab_entry_t **entries;     /* array of entry lists (hash table) */
    symtab_stack_t *next;
};

#endif
