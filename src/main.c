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
#include "scanner.h"
#include "y.tab.h"
#include "unicode/ustdio.h"

extern void scanner_parse(Scanner *s); /* y.tab.c */

int scanner_error(Scanner *s, const char *str)
{
    fprintf(stderr, "%s line %d\n", str, s->lineno);
    return 1;
}

UFILE *ustdout;

int scanner_lex(YYSTYPE *lvalp, Scanner *s)
{
    scanner_token(s);
    /* force re-read on comments */
    if (s->name == T_COMMENT) {
        return scanner_lex(lvalp, s);
    }
    else {
/*
        if (s->name == T_NUMBER) {
            yylval.number = ...
        }
        else if (s->name == T_IDENTIFIER || s->name == T_STRING) {
            yylval.string = ...
        }
*/
        return s->name;
    }
}

int main()
{
    Scanner *s;

    ustdout = u_finit(stdout, NULL, NULL);

    s = scanner_init();
    scanner_parse(s);
    scanner_free(s);

    return EXIT_SUCCESS;
}
