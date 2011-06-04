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
#include "parser.h"

Parser *parser_init(const char *fname) {
    // allocate parser
    Parser *p = (Parser *)calloc(1, sizeof(Parser));

    // bind scanner to parser
    p->s = scanner_init(fname);

    return p;
}

AST *parser_stmt(Parser *p) {
    AST *ast;
    Token *t;
    TokenNode *node;

    // initialize abstract syntax tree
    ast = ast_init();

    // @TESTING: for now just read in stream of tokens
    do {
        if ((t = scanner_token(p->s))) {
            node = tnode_init(t);
            ast_append(ast, node);
        }
    }
    while (t != NULL && t->name != T_SEMICOLON);

    // @TESTING: for now statements are expect to end with a semicolon
    if (t == NULL && ast && ast->tail && ast->tail->data->name != T_SEMICOLON) {
        fprintf(stderr, "parser:%s:%d: Incomplete statement\n", p->s->fname, p->s->lineno);
        exit(EXIT_FAILURE);
    }

    return ast; 
}

void parser_free(Parser *p) {
    // free used memory
    scanner_free(p->s);
    free(p);
}
