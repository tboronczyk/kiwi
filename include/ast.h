/*
 * Copyright (c) 2011, Timothy Boronczyk
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

#ifndef AST_H
#define AST_H

#include "token.h"

// @TODO: TokenNode and AST are currently set up as linked list... this will
// change as language grammar evolves and a proper tree is built

typedef struct _TokenNode {
    Token *data;              // Token structure
    struct _TokenNode *next;  // pointer to next TokenNode in list
} TokenNode;

typedef struct _AST {
    int count;        // Number of TokenNodes contained in list
    TokenNode *head;  // TokenNode at head of list
    TokenNode *tail;  // TokenNode at tail
} AST;

TokenNode *tnode_init(Token *);
void tnode_free(TokenNode *n);

AST *ast_init(void);
void ast_append(AST *, TokenNode *);
void ast_free(AST *);

void ast_dump(AST *);

#endif
