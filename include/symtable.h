#ifndef SYMTABLE_H
#define SYMTABLE_H

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

typedef struct s_SymTable SymTable;
typedef struct s_SymTable_Entry SymTable_Entry;
typedef struct s_SymTable_Stack SymTable_Stack;

/* possible data types */
typedef enum
{
    SYMTAB_ENTRY_NUMBER,
    SYMTAB_ENTRY_STRING,
    SYMTAB_ENTRY_FUNC
}
SymTable_EntryType;

SymTable *symtable_init(void);
void symtable_enter_scope(SymTable *);
void symtable_leave_scope(SymTable *);
void symtable_insert(SymTable *, char *, SymTable_EntryType, void *);
void *symtable_lookup(SymTable *, char *);
void symtable_delete(SymTable *, char *);

/* symbol table hash table */
struct s_SymTable
{
    unsigned int (*hash)(char *);  /* function used for hashing keys */
    SymTable_Entry **entries;      /* array of entry lists (hash table) */
    SymTable_Stack *stack;         /* array of entry lists (hash table) */
};

/* linked list bucket for symtable_t hash table */
struct s_SymTable_Entry
{
    char *key;
    void *value;
    SymTable_EntryType type;
    SymTable_Entry *next;
};

/* linked list bucket for symtable_t hash table */
struct s_SymTable_Stack
{
    SymTable_Entry **entries;     /* array of entry lists (hash table) */
    SymTable_Stack *next;
};

#endif
