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

#include "parser.h"
#include "unicode/ustdio.h"

UFILE *ustdout;

static void parsefile(Parser *p) {
    AST *ast;
    while (1) {
        ast = parser_parse_stmt(p);
        if (ast->count) {
            ast_dump(ast);
            ast_free(ast);
        }
        else {
            break;
        }
    }
    ast_free(ast);
}

int main(int argc, char **argv)
{
    Parser *p;
    int i;

    ustdout = u_finit(stdout, NULL, NULL);

    if (argc == 1) {
        p = parser_init("stdin");
        parsefile(p);
        parser_free(p);
    }
    else {
        for (i = 1; i < argc; i++) {
            p = parser_init(argv[i]);
            parsefile(p);
            parser_free(p);
        }
    }

    u_fclose(ustdout);

    return 0;
}
