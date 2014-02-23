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

#define RETURN_ON_ERROR(stmt,err) do { \
    if (((err) = (stmt)) != SCANERR_OK) { return (err); } \
} while (0)

#define SET_SINGLE(s,t,err) do { \
    (s)->name = (t); \
    RETURN_ON_ERROR(append_advance((s)), (err)); \
} while (0)

#define SET_DOUBLE(s,x,t,err) do { \
    RETURN_ON_ERROR(append_advance((s)), (err)); \
    if ((s)->c == x) { SET_SINGLE((s), (t), (err)); } \
    else { return SCANERR_UNEXPECTED_LEX; } \
} while (0)

#define SET_MAYBE_DOUBLE(s,x,t1,t2,err) do { \
    RETURN_ON_ERROR(append_advance((s)), (err)); \
    if ((s)->c == x) { SET_SINGLE((s), (t2), (err)); } \
    else { (s)->name = (t1); } \
} while (0)
    
static Scanner_ErrCode append_advance(Scanner *);
static Scanner_ErrCode buffer_append(Scanner *);
static Scanner_ErrCode buffer_grow(Scanner *);
static Scanner_ErrCode buffer_init(Scanner *);
static void buffer_reset(Scanner *);
static Scanner_ErrCode read_comment_multi_inner(Scanner *);
static Scanner_ErrCode read_identifier(Scanner *);
static Scanner_ErrCode read_number(Scanner *);
static Scanner_ErrCode read_slash(Scanner *);
static Scanner_ErrCode read_string(Scanner *);
static Scanner_ErrCode stream_advance(Scanner *);
static Scanner_ErrCode stream_init(Scanner *);
static Scanner_ErrCode stream_read_token(Scanner *);
static Scanner_ErrCode stream_skip_whitespace(Scanner *);

Scanner_ErrCode scanner_init(Scanner **s)
{
    char *fname = "stdin";

    /* allocate scanner */
    *s = calloc(1, sizeof(Scanner));
    if (*s == NULL) {
        return SCANERR_ALLOC_SCANNER;
    }

    /* set filename */
    assert((*s)->fname == NULL);
    (*s)->fname = calloc(strlen(fname) + 1, sizeof(char));
    if ((*s)->fname == NULL) {
        free(*s);
        return SCANERR_ALLOC_FILENAME;
    }
    memcpy((*s)->fname, fname, sizeof(char) * strlen(fname));

    /* open stream to file */
    assert((*s)->fp == NULL);
    if (strcmp("stdin", fname) == 0) {
        (*s)->fp = u_finit(stdin, NULL, NULL);
    }
    else {
        (*s)->fp = u_fopen(fname, "r", NULL, NULL);
    }
    if ((*s)->fp == NULL) {
        free((*s)->fname);
        free(*s);
        return SCANERR_FILEOPEN;
    }

    /* initialize scanner */
    Scanner_ErrCode err = buffer_init(*s);
    if (err != SCANERR_OK) {
        scanner_free(*s);
        return err;
    }
    stream_init(*s);
    return SCANERR_OK;
}

Scanner_ErrCode scanner_token(Scanner *s) 
{
    Scanner_ErrCode err;

    buffer_reset(s);

    /* advance stream past whitespace */
    RETURN_ON_ERROR(stream_skip_whitespace(s), err);

    /* obtain token */
    RETURN_ON_ERROR(stream_read_token(s), err);

    return SCANERR_OK;
}

void scanner_free(Scanner *s) 
{
    if (s != NULL) {
        u_fclose(s->fp);
        free(s->fname);
        free(s->tbuf);
        free(s);
    }
}

static Scanner_ErrCode append_advance(Scanner *s)
{
    Scanner_ErrCode err;
    RETURN_ON_ERROR(buffer_append(s), err);
    RETURN_ON_ERROR(stream_advance(s), err);
    return SCANERR_OK;
}

static Scanner_ErrCode buffer_append(Scanner *s)
{
    s->tbuf[s->ti] = s->c;
    s->ti++;
    /* increase token buffer size if necessary */
    if (s->ti == (unsigned int)s->tlen) {
        Scanner_ErrCode err;
        RETURN_ON_ERROR(buffer_grow(s), err);
    }
    return SCANERR_OK;
}

static Scanner_ErrCode buffer_grow(Scanner *s)
{
    /* increase storage of token buffer */
    s->tlen += SCANBUF_SIZE_INCR;
    if ((s->tbuf = realloc(s->tbuf, sizeof(UChar) * s->tlen)) == NULL) {
        return SCANERR_REALLOC_BUFFER;
    }
    /* ensure new buffer space is clear */
    memset(&s->tbuf[s->ti], 0, sizeof(UChar) * (s->tlen - s->ti));
    return SCANERR_OK;
}

static Scanner_ErrCode buffer_init(Scanner *s)
{
    /* initialize token buffer */
    s->ti = 0;
    s->tlen = SCANBUF_SIZE_INIT;
    assert(s->tbuf == NULL);
    if ((s->tbuf = calloc((size_t)s->tlen, sizeof(UChar))) == NULL) {
        return SCANERR_ALLOC_BUFFER;
    }
    return SCANERR_OK;
}

static void buffer_reset(Scanner *s)
{
    s->ti = 0;
    memset(s->tbuf, 0, sizeof(UChar) * s->tlen);
}

static Scanner_ErrCode read_comment_multi_inner(Scanner *s)
{
    Scanner_ErrCode err;
    /* need to keep track of previous character */
    UChar prev;
    prev = s->c;

    /* read characters until end of comment is seen */
    RETURN_ON_ERROR(append_advance(s), err);
    while (!(prev == (UChar)'*' && s->c == (UChar)'/')) {
        /* support nested comments */
        if (prev == (UChar)'/' && s->c == (UChar)'*') {
            RETURN_ON_ERROR(read_comment_multi_inner(s), err);
        }
        prev = s->c;
        RETURN_ON_ERROR(append_advance(s), err);
    }
    RETURN_ON_ERROR(append_advance(s), err);
    return SCANERR_OK;
}

static Scanner_ErrCode read_identifier(Scanner *s)
{
    Scanner_ErrCode err;

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

    RETURN_ON_ERROR(append_advance(s), err);
    while (u_isIDPart(s->c) == TRUE || s->c == (UChar)'_') {
        RETURN_ON_ERROR(append_advance(s), err);
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
        return SCANERR_UNEXPECTED_LEX;
    }

    return SCANERR_OK;
}

static Scanner_ErrCode read_number(Scanner *s)
{
    Scanner_ErrCode err;
    int i;

    static int init = 1;
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

    i = 0;
    /* append digits, initially assuming a base-10 number */
    s->name = T_NUMBER;
    while (u_isdigit(s->c) == TRUE) {
        RETURN_ON_ERROR(append_advance(s), err);
    }

    /* read value for non base-10 numbers */
    if (s->c == (UChar)'#') {
        /* read hexadecimal number, assume hex if there is no leading radix */
        if (s->ti == 0 || u_strcmp(s->tbuf, ustr_base16) == 0) {
            RETURN_ON_ERROR(append_advance(s), err);
            while (u_isxdigit(s->c) == TRUE) {
                i++;
                RETURN_ON_ERROR(append_advance(s), err);
            }
        }
        /* read binary number */
        else if (u_strcmp(s->tbuf, ustr_base2) == 0) {
            RETURN_ON_ERROR(append_advance(s), err);
            while (s->c == (UChar)'0' || s->c == (UChar)'1') {
                i++;
                RETURN_ON_ERROR(append_advance(s), err);
            }
        }
        /* read octal number */
        else if (u_strcmp(s->tbuf, ustr_base8) == 0) {
            RETURN_ON_ERROR(append_advance(s), err);
            while (u_memchr(ustr_octvals, s->c, u_strlen(ustr_octvals)) != NULL) {
                i++;
                RETURN_ON_ERROR(append_advance(s), err);
            }
        }
        /* only bases 2, 8, and 16 are supported */
        else {
            return SCANERR_UNEXPECTED_LEX;
        }

        /* radix with no number part is invalid */
        if (i == 0) {
            return SCANERR_UNEXPECTED_LEX;
        }
    }
    return SCANERR_OK;
}

static Scanner_ErrCode read_slash(Scanner *s)
{
    Scanner_ErrCode err;

    RETURN_ON_ERROR(append_advance(s), err);
    /* match single-line comment */
    if (s->c == (UChar)'/') {
        s->name = T_COMMENT;
        while (s->c != (UChar)'\n') {
            RETURN_ON_ERROR(append_advance(s), err);
        }
    }
    /* match multi-line comment */
    else if (s->c == (UChar)'*') {
        s->name = T_COMMENT;
        RETURN_ON_ERROR(append_advance(s), err);
        RETURN_ON_ERROR(read_comment_multi_inner(s), err);
    }
    /* match shorthand divide assign operator */
    else if (s->c == (UChar)'=') {
        SET_SINGLE(s, T_DIVIDE_ASSIGN, err);
    }
    /* assumed match division operator */
    else {
        s->name = T_DIVIDE;
    }
    return SCANERR_OK;
}

static Scanner_ErrCode read_string(Scanner *s)
{
    Scanner_ErrCode err;
    UChar tmp;
    s->name = T_STRING;

    /* do not include initial quote in string value */
    RETURN_ON_ERROR(stream_advance(s), err);

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
                RETURN_ON_ERROR(buffer_append(s), err);
                s->c = tmp;
            }
        }
        RETURN_ON_ERROR(append_advance(s), err);
    }
    /* do not include final quote in string value */
    RETURN_ON_ERROR(stream_advance(s), err);
    return SCANERR_OK;
}

static Scanner_ErrCode stream_advance(Scanner *s)
{
    s->c = u_fgetc(s->fp);
    s->linepos++;

    /* update file position */
    if (s->c == (UChar)'\n') {
        s->linenum++;
        s->linepos = 1;
    }
    return SCANERR_OK;
}

static Scanner_ErrCode stream_init(Scanner *s)
{
    Scanner_ErrCode err;
    /* set file position and advance through stream to first non-whitespace character */
    s->linenum = 1;
    s->linepos = 1;
    RETURN_ON_ERROR(stream_advance(s), err);
    RETURN_ON_ERROR(stream_skip_whitespace(s), err);
    return SCANERR_OK;
}

static Scanner_ErrCode stream_read_token(Scanner *s)
{
    Scanner_ErrCode err;
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
    if (s->c == U_EOF) { SET_SINGLE(s, T_EOF, err); }
    else if (s->c == (UChar)':') { SET_MAYBE_DOUBLE(s, (UChar)'=', T_COLON, T_ASSIGN, err); }
    else if (s->c == (UChar)'.') { SET_MAYBE_DOUBLE(s, (UChar)'.', T_DOT, T_CONCAT, err); }
    else if (s->c == (UChar)'+') { SET_MAYBE_DOUBLE(s, (UChar)':', T_ADD, T_ADD_ASSIGN, err); }
    else if (s->c == (UChar)'-') { SET_MAYBE_DOUBLE(s, (UChar)':', T_SUBTRACT, T_SUBTRACT_ASSIGN, err); }
    else if (s->c == (UChar)'*') { SET_MAYBE_DOUBLE(s, (UChar)':', T_MULTIPLY, T_MULTIPLY_ASSIGN, err); }
    else if (s->c == (UChar)'%') { SET_MAYBE_DOUBLE(s, (UChar)':', T_MODULO, T_MODULO_ASSIGN, err); }
    else if (s->c == (UChar)'/') { RETURN_ON_ERROR(read_slash(s), err); }
    else if (s->c == (UChar)'=') { SET_SINGLE(s, T_EQUAL, err); }
    else if (s->c == (UChar)'~') { SET_MAYBE_DOUBLE(s, (UChar)'=', T_NOT, T_NOT_EQUAL, err); }
    else if (s->c == (UChar)'<') { SET_MAYBE_DOUBLE(s, (UChar)'=', T_LESS, T_LESS_EQUAL, err); }
    else if (s->c == (UChar)'>') { SET_MAYBE_DOUBLE(s, (UChar)'=', T_GREATER, T_GREATER_EQUAL, err); }
    else if (s->c == (UChar)'&') { SET_DOUBLE(s, (UChar)'&', T_AND, err); }
    else if (s->c == (UChar)'|') { SET_DOUBLE(s, (UChar)'|', T_NOT, err); }
    else if (s->c == (UChar)'{') { SET_SINGLE(s, T_LBRACE, err); }
    else if (s->c == (UChar)'}') { SET_SINGLE(s, T_RBRACE, err); }
    else if (s->c == (UChar)'(') { SET_SINGLE(s, T_LPAREN, err); }
    else if (s->c == (UChar)')') { SET_SINGLE(s, T_RPAREN, err); }
    else if (s->c == (UChar)',') { SET_SINGLE(s, T_COMMA, err); }
    else if (s->c == (UChar)'"') { RETURN_ON_ERROR(read_string(s), err); }
    else if (u_isdigit(s->c) == TRUE || s->c == (UChar)'#') { RETURN_ON_ERROR(read_number(s), err); }
    else if (u_isIDStart(s->c) == TRUE || s->c == (UChar)'_' || s->c == (UChar)'`') {
        RETURN_ON_ERROR(read_identifier(s), err);
/*
        // single _ is wildcard token
        if (u_strcmp(s->tbuf, ustr_wildcard) == 0) {
            s->name = T_WILDCARD;
        }
*/
    }
    /* invalid */
    else {
        return SCANERR_UNEXPECTED_LEX;
    }
    return SCANERR_OK;
}

static Scanner_ErrCode stream_skip_whitespace(Scanner *s)
{
    Scanner_ErrCode err;
    while (u_isWhitespace(s->c) == TRUE) {
        RETURN_ON_ERROR(stream_advance(s), err);
    }
    return SCANERR_OK;
}
