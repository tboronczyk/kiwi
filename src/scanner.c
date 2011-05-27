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
#include <stdlib.h>
#include <string.h>
#include "scanner.h"
#include "token.h"
#include "uhelpers.h"

#define BUFFER_SIZE_INIT 7
#define BUFFER_SIZE_INCR 1.5

#define perror_exit(f) \
    perror(f); \
    exit(EXIT_FAILURE)

#ifdef DEBUG
#define scan_error_exit(s) \
    fprintf(stderr, "%s:%d: Unexpected lexeme (%s)\n", s->fname, s->lineno, __func__), \
    exit(EXIT_FAILURE)
#else
#define scan_error_exit(s) \
    fprintf(stderr, "%s:%d: Unexpected lexeme\n", s->fname, s->lineno), \
    exit(EXIT_FAILURE);
#endif

#define set_single(s,t) \
    s->name = t; \
    append_advance(s)

#define set_double(s,x,t) \
    append_advance(s); \
    if (strcmp(s->ubuf, x) == 0) { set_single(s,t); } else scan_error_exit(s)

#define set_maybe_double(s,x,t1,t2) \
    append_advance(s); \
    if (strcmp(s->ubuf, x) == 0) { set_single(s,t2); } else s->name = t1

static void buffers_init(Scanner *s) {
    // initialize token buffer
    s->ti = 0;
    s->tlen = BUFFER_SIZE_INIT;
    if ((s->tbuf = (char *)calloc(s->tlen, 1)) == NULL) {
        perror_exit("calloc");
    }
    // initialize UTF-8 byte buffer
    s->ui = 0;
    if ((s->ubuf = (char *)calloc(U_MAX_BYTES, sizeof(char))) == NULL) {
        perror_exit("calloc");
    }
}

static void token_buffer_reset(Scanner *s) {
    // clear token buffer
    s->ti = 0;
    memset(s->tbuf, 0, s->tlen);
}

static void token_buffer_grow(Scanner *s) {
    // increase storage of token buffer
    s->tlen = (int)((double)s->tlen * BUFFER_SIZE_INCR);
    if ((s->tbuf = (char *)realloc(s->tbuf, s->tlen)) == NULL) {
        perror_exit("realloc");
    }
    // ensure new buffer space is clear
    memset(s->tbuf + s->ti, 0, s->tlen - s->ti);
}

static void token_buffer_append(Scanner *s) {
    // enough space to append UTF-8 bytes to token buffer?
    while (s->ti + s->ui >= s->tlen) {
        token_buffer_grow(s);
    }
    // append buffer
    memcpy(&s->tbuf[s->ti], s->ubuf, s->ui);
    s->ti += s->ui;
}

static void buffers_free(Scanner *s) {
    // free buffer memory
    free(s->tbuf);
    free(s->ubuf);
}

static void stream_advance(Scanner *s) {
    // obtain next UTF-8 byte sequence
    s->ui = u_getc(s->fp, s->ubuf);
    if (ferror(s->fp)) {
        perror_exit("u_getc");
    }
    // update file position
    if (u_isnewline(s->ubuf)) {
        s->lineno++;
    }
}

static void append_advance(Scanner *s) {
    // append the current UTF-8 byte buffer to the token buffer and
    // advance the stream to the next byte sequence
    token_buffer_append(s);
    stream_advance(s);
}

static void stream_skip_whitespace(Scanner *s) {
    // advance stream to first non-whitespace byte sequence
    while (u_isspace(s->ubuf)) {
        stream_advance(s);
    }
}

static void stream_init(Scanner *s) {
    // set file position and advance through stream to first non-whitespace
    // byte sequence
    s->lineno = 1;
    stream_advance(s);
    stream_skip_whitespace(s);
}

static void read_comment_multi_inner(Scanner *s) {
    // need to keep track of previous byte sequence
    char prev[U_MAX_BYTES];
    memcpy(prev, s->ubuf, U_MAX_BYTES);

    // read bytes until end of comment is seen
    append_advance(s);
    while (!(strcmp(prev, "*") == 0 && strcmp(s->ubuf, "/") == 0)) {
        // support nested comments
        if (strcmp(prev, "/") == 0 && strcmp(s->ubuf, "*") == 0) {
            read_comment_multi_inner(s);
        }
        memcpy(prev, s->ubuf, U_MAX_BYTES);
        append_advance(s);
    }
    append_advance(s);
}

static void read_slash(Scanner *s) {
    append_advance(s);
    // match single-line comment
    if (strcmp(s->ubuf, "/") == 0) {
        s->name = T_COMMENT;
        while (!u_isnewline(s->ubuf)) {
            append_advance(s);
        }
    }
    // match multi-line comment
    else if (strcmp(s->ubuf, "*") == 0) {
        s->name = T_COMMENT_MULTI;
        append_advance(s);
        read_comment_multi_inner(s);
    }
    // match shorthand divide assign operator
    else if (strcmp(s->ubuf, "=") == 0) {
        set_single(s, T_DIVIDE_ASSIGN);
    }
    // assumed match division operator
    else {
        s->name = T_DIVIDE;
    }
}

static void read_number(Scanner *s) {
    long base = 10;
    int j = 0;
    // append digits, initially assuming a base-10 number
    s->name = T_NUMBER;
    while (u_isdigit(s->ubuf)) {
        append_advance(s);
    }

    // read value for non base-10 numbers
    if (strcmp(s->ubuf, "#") == 0) {
        // assume hexadecimal if there is no leading radix
        if (s->ti == 0) {
            base = 16;
        }
        else {
            base = strtol(s->tbuf, NULL, 10);
        }
        // only base 2, 8 and 16 are supported
        if (base != 2 && base != 8 && base != 16) {
            scan_error_exit(s);
        }
        append_advance(s);

        switch (base) {
            // read binary number
            case 2:
                s->name = T_NUMBER_INT_2;
                while (u_is2digit(s->ubuf)) {
                    j++;
                    append_advance(s);
                }
                break;

                // read octal number
            case 8:
                s->name = T_NUMBER_INT_8;
                while (u_is8digit(s->ubuf)) {
                    j++;
                    append_advance(s);
                }
                break;

                // read hexadecimal number
            case 16:
                s->name = T_NUMBER_INT_16;
                while (u_is16digit(s->ubuf)) {
                    j++;
                    append_advance(s);
                }
                break;
        }

        // radix with no number part is invalid
        if (j == 0) {
            scan_error_exit(s);
        }
    }
}

static void read_alpha(Scanner *s) {
    while (u_isalnum(s->ubuf) || strcmp(s->ubuf, "_") == 0) {
        append_advance(s);
    }
    // match keywords
    if (strcmp(s->tbuf, "else") == 0) { s->name = T_ELSE; }
    else if (strcmp(s->tbuf, "if") == 0) { s->name = T_IF; }
    else if (strcmp(s->tbuf, "is") == 0) { s->name = T_IS; }
    // ...
    // assign token as identifier
    else {
        s->name = T_IDENTIFIER;
    }
}

static void stream_read_token(Scanner *s) {
    // the first byte sequence will determine the parsing logic for various
    // tokens
    if (strcmp(s->ubuf, ":") == 0) { set_double(s, "=", T_ASSIGN); }
    else if (strcmp(s->ubuf, "+") == 0) { set_maybe_double(s, "=", T_ADD, T_ADD_ASSIGN); }
    else if (strcmp(s->ubuf, "-") == 0) { set_maybe_double(s, "=", T_SUBTRACT, T_SUBTRACT_ASSIGN); }
    else if (strcmp(s->ubuf, "*") == 0) { set_maybe_double(s, "=", T_MULTIPLY, T_MULTIPLY_ASSIGN); }
    else if (strcmp(s->ubuf, "/") == 0) { read_slash(s); }
    else if (strcmp(s->ubuf, "=") == 0) { set_double(s, "=", T_EQUAL); }
    else if (strcmp(s->ubuf, "~") == 0) { set_maybe_double(s, "=", T_LOG_NOT, T_NOT_EQUAL); }
    else if (strcmp(s->ubuf, "<") == 0) { set_maybe_double(s, "=", T_LESS, T_LESS_EQUAL); }
    else if (strcmp(s->ubuf, ">") == 0) { set_maybe_double(s, "=", T_GREATER, T_GREATER_EQUAL); }
    else if (strcmp(s->ubuf, "&") == 0) { set_double(s, "&", T_LOG_AND); }
    else if (strcmp(s->ubuf, "|") == 0) { set_double(s, "&", T_LOG_NOT); }
    else if (strcmp(s->ubuf, "^") == 0) { set_double(s, "^", T_LOG_XOR); }
    else if (strcmp(s->ubuf, "{") == 0) { set_single(s, T_BRACE_LEFT); }
    else if (strcmp(s->ubuf, "}") == 0) { set_single(s, T_BRACE_RIGHT); }
    else if (strcmp(s->ubuf, "(") == 0) { set_single(s, T_PAREN_LEFT); }
    else if (strcmp(s->ubuf, ")") == 0) { set_single(s, T_PAREN_RIGHT); }
    else if (u_isdigit(s->ubuf) || strcmp(s->ubuf, "#") == 0) { read_number(s); }
    else if (u_isalpha(s->ubuf) || strcmp(s->ubuf, "_") == 0) {
        read_alpha(s);
        // single _ is wildcard token
        if (strcmp(s->tbuf, "_") == 0) {
            s->name = T_WILDCARD;
        }
    }
    // invalid
    else {
        scan_error_exit(s);
    }
}

Scanner *scanner_init(const char *fname) {
    Scanner *s;
    // allocate scanner
    if ((s = (Scanner *)calloc(1, sizeof(Scanner))) == NULL) {
        perror_exit("calloc");
    }
    // set filename
    if ((s->fname = (char *)calloc(strlen(fname) + 1, sizeof(char))) == NULL) {
        perror_exit("calloc");
    }
    memcpy(s->fname, fname, strlen(fname));
    // open stream to file
    if (strcmp("stdin", fname) == 0) {
        s->fp = stdin;
    }
    else if ((s->fp = fopen(fname, "rb")) == NULL) {
        perror_exit("fopen");
    }
    // initialize scanner
    buffers_init(s);
    stream_init(s);

    return s;
}

Token *scanner_token(Scanner *s) {
    // allocate token
    Token *t = (Token *)calloc(1, sizeof(Token));
    token_buffer_reset(s);

    // advance stream past whitespace
    stream_skip_whitespace(s);


    // end of file was reached or an error was encountered
    if (s->ui == 0 && s->ubuf[0]== EOF) {
        free(t);
        return 0;
    }
    // obtain and return token
    else {
        stream_read_token(s);
        t->name = s->name;
        if ((t->lexeme = (char *)calloc(s->ti + 1, 1)) == NULL) {
            perror_exit("calloc");
        }
        memmove(t->lexeme, s->tbuf, s->ti);

        return t;
    }
}

void scanner_free(Scanner *s) {
    // close file stream and free memory
    if (s->fp != stdin) {
        fclose(s->fp);
    }
    free(s->fname);
    buffers_free(s);
    free(s);
}

