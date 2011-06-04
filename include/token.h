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

#ifndef TOKEN_H
#define TOKEN_H

#include "unicode/umachine.h"

// operators
#define T_ASSIGN                 1
#define T_ADD                    2
#define T_ADD_ASSIGN             3
#define T_SUBTRACT               4
#define T_SUBTRACT_ASSIGN        5
#define T_MULTIPLY               6
#define T_MULTIPLY_ASSIGN        7
#define T_DIVIDE                 8
#define T_DIVIDE_ASSIGN          9
#define T_EQUAL                  10
#define T_NOT_EQUAL              11
#define T_LESS                   12
#define T_LESS_EQUAL             13
#define T_GREATER                14
#define T_GREATER_EQUAL          15
#define T_LOG_AND                16
#define T_LOG_OR                 17
#define T_LOG_XOR                18
#define T_LOG_NOT                19

// literals
#define T_WILDCARD               20
#define T_NUMBER                 21
#define T_NUMBER_INT_2           22
#define T_NUMBER_INT_8           23
#define T_NUMBER_INT_16          24
#define T_STRING                 25
#define T_TRUE                   26
#define T_FALSE                  27

// punctuators
#define T_BRACE_LEFT             28
#define T_BRACE_RIGHT            29
#define T_PAREN_LEFT             30
#define T_PAREN_RIGHT            31
#define T_COLON                  32
#define T_SEMICOLON              33
#define T_COMMA                  34

// comments
#define T_COMMENT                35
#define T_COMMENT_MULTI          36

// identifier
#define T_IDENTIFIER             37

// keywords
#define T_IF                     38
#define T_ELSE                   39
#define T_IS                     40
#define T_VAR                    41
#define T_WHILE                  42
#define T_NUMBER_TYPE            43
#define T_STRING_TYPE            44
#define T_BOOLEAN_TYPE           45

typedef struct {
    int name;
    UChar *lexeme;
} Token;

void token_free(Token *t);

#endif
