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

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "scanner.h"
#include "unicode/uchar.h"
#include "unicode/ustdio.h"
#include "unicode/ustring.h"
#include "y.tab.h"

/* main.c */
extern UFILE *ustdin;
extern UFILE *ustdout;
extern UFILE *ustderr;

#define SCANBUF_SIZE_INIT 10 
#define SCANBUF_SIZE_INCR 5

#ifdef DEBUG
#define scan_error_exit(s) \
    u_fprintf(ustderr, "scanner:%s:%d: Unexpected lexeme (%s)\n", s->fname, s->lineno, __func__), \
    exit(EXIT_FAILURE)
#else
#define scan_error_exit(s) \
    u_fprintf(ustderr, "scanner:%s:%d: Unexpected lexeme\n", s->fname, s->lineno), \
    exit(EXIT_FAILURE);
#endif

#define set_single(s,t) \
    s->name = t; \
    append_advance(s)

#define set_double(s,x,t) \
    append_advance(s); \
    if (s->c == x) { set_single(s, t); } else scan_error_exit(s)

#define set_maybe_double(s,x,t1,t2) \
    append_advance(s); \
    if (s->c == x) { set_single(s, t2); } else s->name = t1

static void append_advance(Scanner *);
static void buffer_append(Scanner *);
static void buffer_grow(Scanner *);
static void buffer_init(Scanner *);
static void buffer_reset(Scanner *);
static void read_comment_multi_inner(Scanner *);
static void read_identifier(Scanner *);
static void read_number(Scanner *);
static void read_slash(Scanner *);
static void read_string(Scanner *);
static void stream_advance(Scanner *);
static void stream_init(Scanner *);
static void stream_read_token(Scanner *);
static void stream_skip_whitespace(Scanner *);

Scanner* scanner_init(void) 
{
    char *fname = "stdin";
    Scanner *s;

    /* allocate scanner */
    if ((s = (Scanner *)calloc(1, sizeof(Scanner))) == NULL) {
        perror("Allocate scanner failed");
        exit(EXIT_FAILURE);
    }
    /* set filename */
    assert(s->fname == NULL);
    if ((s->fname = (char *)calloc(strlen(fname) + 1, sizeof(char))) == NULL) {
        perror("Allocate scanner filename failed");
        exit(EXIT_FAILURE);
    }
    memcpy(s->fname, fname, sizeof(char) * strlen(fname));
    /* open stream to file */
    assert(s->fp == NULL);
    if (strcmp("stdin", fname) == 0) {
        s->fp = u_finit(stdin, NULL, NULL);
    }
    else if ((s->fp = u_fopen(fname, "r", NULL, NULL)) == NULL) {
        perror("Open file failed");
        exit(EXIT_FAILURE);
    }
    /* initialize scanner */
    buffer_init(s);
    stream_init(s);
    return s;
}

void scanner_token(Scanner *s) 
{
    buffer_reset(s);

    /* advance stream past whitespace */
    stream_skip_whitespace(s);

    /* obtain token */
    stream_read_token(s);
}

void scanner_free(Scanner *s) 
{
    u_fclose(s->fp);
    free(s->fname);
    free(s->tbuf);
    free(s);
}

static void append_advance(Scanner *s)
{
    buffer_append(s);
    stream_advance(s);
}

static void buffer_append(Scanner *s)
{
    s->tbuf[s->ti] = s->c;
    s->ti++;
    /* increase token buffer size if necessary */
    if (s->ti == (unsigned int)s->tlen) {
        buffer_grow(s);
    }
}

static void buffer_grow(Scanner *s)
{
    /* increase storage of token buffer */
    s->tlen += SCANBUF_SIZE_INCR;
    if ((s->tbuf = (UChar *)realloc(s->tbuf, sizeof(UChar) * s->tlen)) == NULL) {
        perror("Reallocate token buffer failed");
        exit(EXIT_FAILURE);
    }
    /* ensure new buffer space is clear */
    memset(&s->tbuf[s->ti], 0, sizeof(UChar) * (s->tlen - s->ti));
}

static void buffer_init(Scanner *s)
{
    /* initialize token buffer */
    s->ti = 0;
    s->tlen = SCANBUF_SIZE_INIT;
    assert(s->tbuf == NULL);
    if ((s->tbuf = (UChar *)calloc((size_t)s->tlen, sizeof(UChar))) == NULL) {
        perror("Allocate token buffer failed");
        exit(EXIT_FAILURE);
    }
}

static void buffer_reset(Scanner *s)
{
    s->ti = 0;
    memset(s->tbuf, 0, sizeof(UChar) * s->tlen);
}

static void read_comment_multi_inner(Scanner *s)
{
    /* need to keep track of previous character */
    UChar prev;
    prev = s->c;

    /* read characters until end of comment is seen */
    append_advance(s);
    while (!(prev == (UChar)'*' && s->c == (UChar)'/')) {
        /* support nested comments */
        if (prev == (UChar)'/' && s->c == (UChar)'*') {
            read_comment_multi_inner(s);
        }
        prev = s->c;
        append_advance(s);
    }
    append_advance(s);
}

static void read_identifier(Scanner *s)
{
    static int init = 1;

    /* string literals for comparison */
    U_STRING_DECL(ustr_else, "else", 4);
    U_STRING_DECL(ustr_if, "if", 2);
    U_STRING_DECL(ustr_is, "is", 2);
    U_STRING_DECL(ustr_var, "var", 3);
    U_STRING_DECL(ustr_while, "while", 5);
    U_STRING_DECL(ustr_true, "true", 4);
    U_STRING_DECL(ustr_false, "false", 5);
    U_STRING_DECL(ustr_func, "func", 4);
    U_STRING_DECL(ustr_return, "return", 6);
    U_STRING_DECL(ustr_backtick, "`", 1);
    if (init == 1) {
        U_STRING_INIT(ustr_else, "else", 4);
        U_STRING_INIT(ustr_if, "if", 2);
        U_STRING_INIT(ustr_is, "is", 2);
        U_STRING_INIT(ustr_var, "var", 3);
        U_STRING_INIT(ustr_while, "while", 5);
        U_STRING_INIT(ustr_true, "true", 4);
        U_STRING_INIT(ustr_false, "false", 5);
        U_STRING_INIT(ustr_func, "func", 4);
        U_STRING_INIT(ustr_return, "return", 6);
        U_STRING_INIT(ustr_backtick, "`", 1);
        init = 0;
    }

    append_advance(s);
    while (u_isIDPart(s->c) == TRUE || s->c == (UChar)'_') {
        append_advance(s);
    }
    /* match keywords */
    if (u_strcmp(s->tbuf, ustr_else) == 0) { s->name = T_ELSE; }
    else if (u_strcmp(s->tbuf, ustr_if) == 0) { s->name = T_IF; }
    else if (u_strcmp(s->tbuf, ustr_is) == 0) { s->name = T_IS; }
    else if (u_strcmp(s->tbuf, ustr_var) == 0) { s->name = T_VAR; }
    else if (u_strcmp(s->tbuf, ustr_while) == 0) { s->name = T_WHILE; }
    else if (u_strcmp(s->tbuf, ustr_true) == 0) { s->name = T_TRUE; }
    else if (u_strcmp(s->tbuf, ustr_false) == 0) { s->name = T_FALSE; }
    else if (u_strcmp(s->tbuf, ustr_func) == 0) { s->name = T_FUNC; }
    else if (u_strcmp(s->tbuf, ustr_return) == 0) { s->name = T_RETURN; }
    /* ... */

    /* assign token as identifier */
    else {
        s->name = T_IDENTIFIER;
    }

    /* backtick is a convenience to allow a programmer to use a reserved-keyword
       as an identifier, a single backtick itself is not considered valid */
    if (u_strcmp(s->tbuf, ustr_backtick) == 0) {
        scan_error_exit(s);
    }
}

static void read_number(Scanner *s)
{
    static int init = 1;
    int j;

    /* string literals for comparison */
    U_STRING_DECL(ustr_base2, "2", 1);
    U_STRING_DECL(ustr_base8, "8", 1);
    U_STRING_DECL(ustr_base16, "16", 2);
    U_STRING_DECL(ustr_octvals, "01234567", 8);
    if (init == 1) {
        U_STRING_INIT(ustr_base2, "2", 1);
        U_STRING_INIT(ustr_base8, "8", 1);
        U_STRING_INIT(ustr_base16, "16", 2);
        U_STRING_INIT(ustr_octvals, "01234567", 8);
        init = 0;
    }

    j = 0;
    /* append digits, initially assuming a base-10 number */
    s->name = T_NUMBER;
    while (u_isdigit(s->c) == TRUE) {
        append_advance(s);
    }

    /* read value for non base-10 numbers */
    if (s->c == (UChar)'#') {
        /* read hexadecimal number, assume hex if there is no leading radix */
        if (s->ti == 0 || u_strcmp(s->tbuf, ustr_base16) == 0) {
            append_advance(s);
            while (u_isxdigit(s->c) == TRUE) {
                j++;
                append_advance(s);
            }
        }
        /* read binary number */
        else if (u_strcmp(s->tbuf, ustr_base2) == 0) {
            append_advance(s);
            while (s->c == (UChar)'0' || s->c == (UChar)'1') {
                j++;
                append_advance(s);
            }
        }
        /* read octal number */
        else if (u_strcmp(s->tbuf, ustr_base8) == 0) {
            append_advance(s);
            while (u_memchr(ustr_octvals, s->c, u_strlen(ustr_octvals)) != NULL) {
                j++;
                append_advance(s);
            }
        }
        /* only bases 2, 8, and 16 are supported */
        else {
            scan_error_exit(s);
        }

        /* radix with no number part is invalid */
        if (j == 0) {
            scan_error_exit(s);
        }
    }
}

static void read_slash(Scanner *s)
{
    append_advance(s);
    /* match single-line comment */
    if (s->c == (UChar)'/') {
        s->name = T_COMMENT;
        while (s->c != (UChar)'\n') {
            append_advance(s);
        }
    }
    /* match multi-line comment */
    else if (s->c == (UChar)'*') {
        s->name = T_COMMENT;
        append_advance(s);
        read_comment_multi_inner(s);
    }
    /* match shorthand divide assign operator */
    else if (s->c == (UChar)'=') {
        set_single(s, T_DIVIDE_ASSIGN);
    }
    /* assumed match division operator */
    else {
        s->name = T_DIVIDE;
    }
}

static void read_string(Scanner *s)
{
    UChar tmp;
    s->name = T_STRING;

    /* do not include initial quote in string value */
    stream_advance(s);

    while (s->c != (UChar)'"') {
        /* handle escaped literals */
        if (s->c == (UChar)'\\') {
            stream_advance(s);
            if (s->c == (UChar)'"') { }
            else if (s->c == (UChar)'r') { s->c = (UChar)'\r'; }
            else if (s->c == (UChar)'n') { s->c = (UChar)'\n'; }
            else if (s->c == (UChar)'t') { s->c = (UChar)'\t'; }
            else if (s->c == (UChar)'\\') { s->c =(UChar) '\\'; }
            else {
                tmp = s->c;
                s->c = (UChar)'\\';
                buffer_append(s);
                s->c = tmp;
            }
        }
        append_advance(s);
    }
    /* do not include final quote in string value */
    stream_advance(s);
}

static void stream_advance(Scanner *s)
{
    s->c = u_fgetc(s->fp);

    /* update file position */
    if (s->c == (UChar)'\n') {
        s->lineno++;
    }
}

static void stream_init(Scanner *s)
{
    /* set file position and advance through stream to first non-whitespace
       character */
    s->lineno = 1;
    stream_advance(s);
    stream_skip_whitespace(s);
}

static void stream_read_token(Scanner *s)
{
/*
    static int init = 1;

    // string literals for comparison
    U_STRING_DECL(ustr_wildcard, "_", 1);
    if (init == 1) {
        U_STRING_INIT(ustr_wildcard, "_", 1);
        init = 0;
    }
*/
    /* the first character will determine the parsing logic for various tokens */
    if (s->c == U_EOF) { set_single(s, T_EOF); }
    else if (s->c == (UChar)':') { set_maybe_double(s, (UChar)'=', T_COLON, T_ASSIGN); }
    else if (s->c == (UChar)'.') { set_maybe_double(s, (UChar)'.', T_DOT, T_CONCAT); }
    else if (s->c == (UChar)'+') { set_maybe_double(s, (UChar)':', T_ADD, T_ADD_ASSIGN); }
    else if (s->c == (UChar)'-') { set_maybe_double(s, (UChar)':', T_SUBTRACT, T_SUBTRACT_ASSIGN); }
    else if (s->c == (UChar)'*') { set_maybe_double(s, (UChar)':', T_MULTIPLY, T_MULTIPLY_ASSIGN); }
    else if (s->c == (UChar)'%') { set_maybe_double(s, (UChar)':', T_MODULO, T_MODULO_ASSIGN); }
    else if (s->c == (UChar)'/') { read_slash(s); }
    else if (s->c == (UChar)'=') { set_single(s, T_EQUAL); }
    else if (s->c == (UChar)'~') { set_maybe_double(s, (UChar)'=', T_NOT, T_NOT_EQUAL); }
    else if (s->c == (UChar)'<') { set_maybe_double(s, (UChar)'=', T_LESS, T_LESS_EQUAL); }
    else if (s->c == (UChar)'>') { set_maybe_double(s, (UChar)'=', T_GREATER, T_GREATER_EQUAL); }
    else if (s->c == (UChar)'&') { set_double(s, (UChar)'&', T_AND); }
    else if (s->c == (UChar)'|') { set_double(s, (UChar)'|', T_NOT); }
    else if (s->c == (UChar)'{') { set_single(s, T_LBRACE); }
    else if (s->c == (UChar)'}') { set_single(s, T_RBRACE); }
    else if (s->c == (UChar)'(') { set_single(s, T_LPAREN); }
    else if (s->c == (UChar)')') { set_single(s, T_RPAREN); }
    else if (s->c == (UChar)',') { set_single(s, T_COMMA); }
    else if (s->c == (UChar)'"') { read_string(s); }
    else if (u_isdigit(s->c) == TRUE || s->c == (UChar)'#') { read_number(s); }
    else if (u_isIDStart(s->c) == TRUE || s->c == (UChar)'_' || s->c == (UChar)'`') {
        read_identifier(s);
/*
        // single _ is wildcard token
        if (u_strcmp(s->tbuf, ustr_wildcard) == 0) {
            s->name = T_WILDCARD;
        }
*/
    }
    /* invalid */
    else {
        scan_error_exit(s);
    }
}

static void stream_skip_whitespace(Scanner *s)
{
    while (u_isWhitespace(s->c) == TRUE) {
        stream_advance(s);
    }
}
