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

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unicode/ustdio.h>
#include "ast.h"
#include "scanner.h"
#include "y.tab.h"

UFILE *ustdin,
      *ustdout,
      *ustderr;

static void tokenize(Scanner *s) {
    do {
        Scanner_ErrCode err = scanner_token(s);
        switch (err) {
            case SCANERR_OK:
                if (s->name == T_EOF) {
                    (void)u_fprintf(ustdout, "Found %d <EOF>\n", s->name);
                }
                else {
                    (void)u_fprintf(ustdout, "Found %d %S\n", s->name, s->tbuf);
                }
                break;

            case SCANERR_UNEXPECTED_LEX:
                (void)u_fprintf(ustderr, "Scanner: Unexpected lexeme in %s %d:%d\n", s->fname, s->linenum, s->linepos);
                return;

            default:
                (void)u_fprintf(ustderr, "Scanner: Error encountered (%d)\n", err);
        }
    }
    while (s->name != T_EOF);
}

int main(void)
{
    Scanner *s;
    Scanner_ErrCode err;

    /* prepare unicode file descriptors */
    ustdin  = u_finit(stdin,  NULL, NULL);
    ustdout = u_finit(stdout, NULL, NULL);
    ustderr = u_finit(stderr, NULL, NULL);

    err = scanner_init(&s);
    if (err != SCANERR_OK) {
        fprintf(stderr, "Allocate scanner failed");
        exit(EXIT_FAILURE);
    }

    tokenize(s);
    scanner_free(s);

    return EXIT_SUCCESS;
}
