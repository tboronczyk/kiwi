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

%pure-parser
%lex-param { scanner_t *s}
%parse-param { scanner_t *s }

%{
#include <stdio.h>
#include "unicode/uchar.h"
#include "unicode/ustdio.h"
#include "unicode/ustring.h"
#include "scanner.h"
#include "y.tab.h"
%}

%union
{
    int number;
    UChar *string;
}

%token <number> T_NUMBER
%token <string> T_IDENTIFIER T_STRING

%token T_IF T_AND T_OR T_NOT T_EQUAL T_ADD T_SUBTRACT T_CONCAT T_MULTIPLY
%token T_DIVIDE T_NOT_EQUAL T_LESS T_LESS_EQUAL T_GREATER T_GREATER_EQUAL
%token T_LBRACE T_RBRACE T_ELSE T_COMMA T_DOT T_ASSIGN T_RETURN T_VAR T_TRUE
%token T_FALSE T_LPAREN T_RPAREN T_WHILE T_FUNC T_IS T_MODULO T_ADD_ASSIGN
%token T_SUBTRACT_ASSIGN T_MULTIPLY_ASSIGN T_DIVIDE_ASSIGN T_MODULO_ASSIGN 
%token T_COMMENT T_COLON

%start program

%%
program
	: /* empty */
	| stmtlist
	;


stmtlist
	: stmt
	| stmtlist stmt
	;

stmt
	: complexstmt
	| simplestmt
	;

complexstmt
	: compoundstmt
	| funcdef
	;

compoundstmt
	: ifstmt
	| whilestmt
	;

ifstmt
	: T_IF expr compoundbody
	| T_IF expr compoundbody elsestmt
	;

expr
	: notexpr
	| expr exprop notexpr
	;

exprop
	: T_AND
	| T_OR
	;
notexpr
	: compareexpr
	| T_NOT compareexpr
	;

compareexpr
	: minorexpr 
	| compareexpr compop minorexpr
	;

compop
	: T_EQUAL
	| T_NOT_EQUAL
	| T_LESS
	| T_LESS_EQUAL
	| T_GREATER
	| T_GREATER_EQUAL
        | T_IS
	;

minorexpr
	: term
	| minorexpr addop term
	;

addop
	: T_ADD
	| T_SUBTRACT;

term
	: factor
	| term mulop factor
	;

mulop
	: T_MULTIPLY
	| T_DIVIDE
	;

factor
	: atom
	| addop factor
	;

atom
	: T_IDENTIFIER
	| funccall
	| T_NUMBER 
	| T_TRUE
	| T_FALSE
	| T_STRING
	| T_LPAREN expr T_RPAREN 
	;

funccall
	: T_IDENTIFIER T_LPAREN T_RPAREN
	| T_IDENTIFIER T_LPAREN exprlist T_RPAREN
	;

exprlist
	: expr
	| exprlist T_COMMA expr
	;

elsestmt
	: T_ELSE compoundbody
	| T_ELSE ifstmt
	;

compoundbody
	: T_LBRACE compoundbodylist T_RBRACE
	;

compoundbodylist
	: stmt
	| compoundbodylist stmt
	;

whilestmt
	: T_WHILE expr compoundbody
	;

funcdef
	: T_FUNC T_IDENTIFIER compoundbody
	| T_FUNC T_IDENTIFIER funcparamlist compoundbody
	;

funcparamlist
	: T_IDENTIFIER
	| funcparamlist T_COMMA T_IDENTIFIER
	;

simplestmt
	: assignstmt T_DOT
	| returnstmt T_DOT
	| varstmt T_DOT
	| expr T_DOT
	;
	
assignstmt
	: T_IDENTIFIER assignop expr
	;

assignop
	: T_ASSIGN
        | T_ADD_ASSIGN
        | T_SUBTRACT_ASSIGN
        | T_MULTIPLY_ASSIGN
        | T_DIVIDE_ASSIGN
        | T_MODULO_ASSIGN
	;


returnstmt
	: T_RETURN expr
	;

varstmt
	: T_VAR varstmtlist
	;

varstmtlist
	: T_IDENTIFIER
	| assignstmt
	| varstmtlist T_COMMA T_IDENTIFIER
	| varstmtlist T_COMMA assignstmt
	;
%%

