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

#include <stdlib.h>
#include <stdio.h>
#include "ast.h"

TokenNode *tnode_init(Token *token) {
    // allocate token node
    TokenNode *n = (TokenNode *)calloc(1, sizeof(TokenNode));

    // set token as data member of token node
    n->data = token;
    return n;
}

void tnode_free(TokenNode *node) {
    // free used memory
    token_free(node->data);
    free(node);
}

AST *ast_init(void) {
    // allocate abstract syntax tree
    AST *ast = calloc(1, sizeof(AST));
    ast->count = 0;
    return ast;
}

void ast_append(AST *ast, TokenNode *node) {
    if (ast->count) {
        // append token node
        ast->tail->next = node;
        ast->tail = ast->tail->next;
    }
    else {
        // set token node as head
        ast->head = node;
        ast->tail = ast->head;
    }
    ast->count++;
}

void ast_dump(AST *ast) {
    // @TESTING: this dumps contents of ast for testing
    TokenNode *n;
    n = ast->head;
    if (ast->count) {
        while (n) {
            printf("%d\n", n->data->name);
            n = n->next;
        }
    }
    else {
        printf("AST IS EMPTY\n");
    }

}

void ast_free(AST *ast) {
    TokenNode *n, *m;

    // traverse tree to free used memory
    n = ast->head;
    while (n) {
        m = n->next;
        tnode_free(n);
        n = m;
    }
    free(ast);
}
