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

#include <stdio.h>
#include "scanner.h"
#include "unicode/ustdio.h"
#include "y.tab.h"

UFILE *ustdout;

static void tokenize(Scanner *s) {
    while (scanner_token(s) != T_EOF) {
        u_fprintf(ustdout, "Found %d %S\n", s->name, s->tbuf);
    }
}

int main(int argc, char **argv)
{
    Scanner *s;
    int i;

    ustdout = u_finit(stdout, NULL, NULL);

    if (argc == 1) {
        s = scanner_init("stdin");
        tokenize(s);
        scanner_free(s);
    }
    else {
        for (i = 1; i < argc; i++) {
            s = scanner_init(argv[i]);
            tokenize(s);
            scanner_free(s);
        }
    }

    return 0;
}
