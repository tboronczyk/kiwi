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
%lex-param { Scanner *s }
%parse-param { Scanner *s }
%parse-param { AST_Program **n }
%{
#include <stdio.h>
#include <stdlib.h>
#include <unicode/uchar.h>
#include <unicode/ustdio.h>
#include <unicode/ustring.h>
#include "ast.h"
#include "scanner.h"
#include "parser.h"

/* main.c */
extern UFILE *ustderr;

int yyerror(Scanner *, AST_Program **, const char *);
int yylex(YYSTYPE *, Scanner *);
%}

%union
{
    UChar *atom;
    Token token;

    AST_AssignStmt *assignstmt;
    AST_CompareExpr *compareexpr;
    AST_ComplexStmt *complexstmt;
    AST_CompoundBody *compoundbody;
    AST_CompoundBodyList *compoundbodylist;
    AST_CompoundStmt *compoundstmt;
    AST_ElseStmt *elsestmt;
    AST_Expr *expr;
    AST_ExprList *exprlist;
    AST_Factor *factor;
    AST_FuncCall *funccall;
    AST_FuncDef *funcdef;
    AST_FuncParamList *funcparamlist;
    AST_IfStmt *ifstmt;
    AST_MinorExpr *minorexpr;
    AST_NotExpr *notexpr;
    AST_Program *program;
    AST_ReturnStmt *returnstmt;
    AST_SimpleStmt *simplestmt;
    AST_Stmt *stmt;
    AST_StmtList *stmtlist;
    AST_Term *term;
    AST_VarStmt *varstmt;
    AST_VarStmtList *varstmtlist;
    AST_WhileStmt *whilestmt;
}

%type <atom> atom
%type <token> addop assignop compareop exprop mulop

%type <assignstmt> assignstmt
%type <compareexpr> compareexpr
%type <complexstmt> complexstmt
%type <compoundbody> compoundbody
%type <compoundbodylist> compoundbodylist
%type <compoundstmt> compoundstmt
%type <elsestmt> elsestmt
%type <expr> expr
%type <exprlist> exprlist
%type <factor> factor
%type <funcdef> funcdef
%type <funccall> funccall
%type <funcparamlist> funcparamlist
%type <ifstmt> ifstmt
%type <minorexpr> minorexpr
%type <notexpr> notexpr
%type <program> program
%type <returnstmt> returnstmt
%type <simplestmt> simplestmt
%type <stmt> stmt
%type <stmtlist> stmtlist
%type <term> term
%type <varstmt> varstmt
%type <varstmtlist> varstmtlist
%type <whilestmt> whilestmt

%token T_EOF 0
%token <token> T_ADD T_SUBTRACT 
%token <token> T_ASSIGN T_ADD_ASSIGN T_SUBTRACT_ASSIGN T_MULTIPLY_ASSIGN
%token <token> T_DIVIDE_ASSIGN T_MODULO_ASSIGN
%token <token> T_EQUAL T_NOT_EQUAL T_IS T_LESS T_LESS_EQUAL T_GREATER 
%token <token> T_GREATER_EQUAL 
%token <token> T_AND T_OR
%token <token> T_MULTIPLY T_DIVIDE T_MODULO

%token <atom> T_NUMBER T_IDENTIFIER T_STRING T_TRUE T_FALSE

%token T_IF T_NOT T_CONCAT
%token T_LBRACE T_RBRACE T_ELSE T_COMMA T_DOT T_RETURN T_VAR
%token T_LPAREN T_RPAREN T_WHILE T_FUNC
%token T_COMMENT T_COLON

%start program

%%
program
: T_EOF {
    *n = ast_program_init();
}
| stmtlist {
    *n = ast_program_init();
    (*n)->stmtlist = $1;
}
;

stmtlist
: stmt {
    $$ = ast_stmtlist_init();
    $$->stmt = $1;
}
| stmtlist stmt {
    $$ = ast_stmtlist_init();
    $$->stmtlist = $1;
    $$->stmt = $2;
}
;

stmt
: complexstmt {
    $$ = ast_stmt_init();
    $$->stmttype = AST_COMPLEXSTMT;
    $$->stmt.complexstmt = $1;
}
| simplestmt {
    $$ = ast_stmt_init();
    $$->stmttype = AST_SIMPLESTMT;
    $$->stmt.simplestmt = $1;
}
;

complexstmt
: compoundstmt {
    $$ = ast_complexstmt_init();
    $$->stmttype = AST_COMPOUNDSTMT;
    $$->stmt.compoundstmt = $1;
}
| funcdef {
    $$ = ast_complexstmt_init();
    $$->stmttype = AST_FUNCDEF;
    $$->stmt.funcdef = $1;
}
;

compoundstmt
: ifstmt {
    $$ = ast_compoundstmt_init();
    $$->stmttype = AST_IFSTMT;
    $$->stmt.ifstmt = $1;
}
| whilestmt {
    $$ = ast_compoundstmt_init();
    $$->stmttype = AST_WHILESTMT;
    $$->stmt.whilestmt = $1;
}
;

ifstmt
: T_IF expr compoundbody {
    $$ = ast_ifstmt_init();
    $$->expr = $2;
    $$->compoundbody = $3;
}
| T_IF expr compoundbody elsestmt {
    $$ = ast_ifstmt_init();
    $$->expr = $2;
    $$->compoundbody = $3;
    $$->elsestmt = $4;
}
;

expr
: notexpr {
    $$ = ast_expr_init();
    $$->notexpr = $1;
}
| expr exprop notexpr {
    $$ = ast_expr_init();
    $$->expr = $1;
    $$->exprop = $2;
    $$->notexpr = $3;
}
;

exprop
: T_AND { $$ = T_AND; }
| T_OR { $$ = T_OR; }
;

notexpr
: compareexpr {
    $$ = ast_notexpr_init();
    $$->compareexpr = $1;
}
| T_NOT compareexpr {
    $$ = ast_notexpr_init();
    $$->tnot = 1;
    $$->compareexpr = $2;
}
;

compareexpr
: minorexpr {
    $$ = ast_compareexpr_init();
    $$->minorexpr = $1;
}
| compareexpr compareop minorexpr {
    $$ = ast_compareexpr_init();
    $$->compareexpr = $1;
    $$->compareop = $2;
    $$->minorexpr = $3;
}
;

compareop
: T_EQUAL { $$ = T_EQUAL; }
| T_NOT_EQUAL { $$ = T_NOT_EQUAL; }
| T_LESS { $$ = T_LESS; }
| T_LESS_EQUAL { $$ = T_LESS_EQUAL; }
| T_GREATER { $$ = T_GREATER; }
| T_GREATER_EQUAL { $$ = T_GREATER_EQUAL; }
| T_IS { $$ = T_IS; }
;

minorexpr
: term {
    $$ = ast_minorexpr_init();
    $$->term = $1;
}
| minorexpr addop term {
    $$ = ast_minorexpr_init();
    $$->minorexpr = $1;
    $$->addop = $2;
    $$->term = $3;
}
;

addop
: T_ADD { $$ = T_ADD; }
| T_SUBTRACT { $$ = T_SUBTRACT; }
;

term
: factor {
    $$ = ast_term_init();
    $$->factor = $1;
}
| term mulop factor {
    $$ = ast_term_init();
    $$->term = $1;
    $$->mulop = $2;
    $$->factor = $3;
}
;

mulop
: T_MULTIPLY { $$ = T_MULTIPLY; }
| T_DIVIDE { $$ = T_DIVIDE; }
| T_MODULO { $$ = T_MODULO; }
;

factor
: atom {
    $$ = ast_factor_init();
    $$->factortype = AST_ATOM;
    $$->factor.atom = $1;
}
| funccall {
    $$ = ast_factor_init();
    $$->factortype = AST_FUNCCALL;
    $$->factor.funccall = $1;
}
| T_LPAREN expr T_RPAREN {
    $$ = ast_factor_init();
    $$->factortype = AST_EXPR;
    $$->factor.expr = $2;
}
| addop factor {
    $$ = ast_factor_init();
    $$->factortype = AST_FACTOR;
    $$->addop = $1;
    $$->factor.factor = $2;
}
;

atom
: T_IDENTIFIER { $$ = $1; }
| T_NUMBER { $$ = $1; }
| T_TRUE { $$ = $1; }
| T_FALSE { $$ = $1; }
| T_STRING { $$ = $1; }
;

funccall
: T_IDENTIFIER T_LPAREN T_RPAREN {
    $$ = ast_funccall_init();
    $$->identifier = $1;
}
| T_IDENTIFIER T_LPAREN exprlist T_RPAREN {
    $$ = ast_funccall_init();
    $$->identifier = $1;
    $$->exprlist = $3;
}
;

exprlist
: expr {
    $$ = ast_exprlist_init();
    $$->expr = $1;
}
| exprlist T_COMMA expr {
    $$ = ast_exprlist_init();
    $$->exprlist = $1;
    $$->expr = $3;
}
;

elsestmt
: T_ELSE compoundbody {
    $$ = ast_elsestmt_init();
    $$->stmttype = AST_COMPOUNDBODY;
    $$->stmt.compoundbody = $2;
}
| T_ELSE ifstmt {
    $$ = ast_elsestmt_init();
    $$->stmttype = AST_IFSTMT;
    $$->stmt.ifstmt = $2;
}
;

compoundbody
: T_LBRACE compoundbodylist T_RBRACE {
    $$ = ast_compoundbody_init();
    $$->compoundbodylist = $2;
}
;

compoundbodylist
: stmt {
    $$ = ast_compoundbodylist_init();
    $$->stmt = $1;
}
| compoundbodylist stmt {
    $$ = ast_compoundbodylist_init();
    $$->compoundbodylist = $1;
    $$->stmt = $2;
}
;

whilestmt
: T_WHILE expr compoundbody {
    $$ = ast_whilestmt_init();
    $$->expr = $2;
    $$->compoundbody = $3;
}
;

funcdef
: T_FUNC T_IDENTIFIER compoundbody {
    $$ = ast_funcdef_init();
    $$->identifier = $2;
    $$->compoundbody = $3;
}
| T_FUNC T_IDENTIFIER funcparamlist compoundbody {
    $$ = ast_funcdef_init();
    $$->identifier = $2;
    $$->funcparamlist = $3;
    $$->compoundbody = $4;
}
;

funcparamlist
: T_IDENTIFIER {
    $$ = ast_funcparamlist_init();
    $$->identifier = $1;
}
| funcparamlist T_COMMA T_IDENTIFIER {
    $$ = ast_funcparamlist_init();
    $$->funcparamlist = $1;
    $$->identifier = $3;
}
;

simplestmt
: assignstmt T_DOT {
    $$ = ast_simplestmt_init();
    $$->stmttype = AST_ASSIGNSTMT;
    $$->stmt.assignstmt = $1;
}
| returnstmt T_DOT {
    $$ = ast_simplestmt_init();
    $$->stmttype = AST_RETURNSTMT;
    $$->stmt.returnstmt = $1;
}
| varstmt T_DOT {
    $$ = ast_simplestmt_init();
    $$->stmttype = AST_VARSTMT;
    $$->stmt.varstmt = $1;
}
| expr T_DOT {
    $$ = ast_simplestmt_init();
    $$->stmttype = AST_EXPR;
    $$->stmt.expr = $1;
}
;

assignstmt
: T_IDENTIFIER assignop expr {
    $$ = ast_assignstmt_init();
    $$->identifier = $1;
    $$->assignop = $2;
    $$->expr = $3;
}
;

assignop
: T_ASSIGN  { $$ = T_ASSIGN; }
| T_ADD_ASSIGN { $$ = T_ADD_ASSIGN; }
| T_SUBTRACT_ASSIGN { $$ = T_SUBTRACT_ASSIGN; }
| T_MULTIPLY_ASSIGN { $$ = T_MULTIPLY_ASSIGN; }
| T_DIVIDE_ASSIGN { $$ = T_DIVIDE_ASSIGN; }
| T_MODULO_ASSIGN { $$ = T_MODULO_ASSIGN; }
;

returnstmt
: T_RETURN expr {
    $$ = ast_returnstmt_init();
    $$->expr = $2;
}
;

varstmt
: T_VAR varstmtlist {
    $$ = ast_varstmt_init();
    $$->varstmtlist = $2;
}
;

varstmtlist
: T_IDENTIFIER {
    $$ = ast_varstmtlist_init();
    $$->stmttype = AST_IDENTIFIER;
    $$->stmt.identifier = $1;
}
| assignstmt {
    $$ = ast_varstmtlist_init();
    $$->stmttype = AST_ASSIGNSTMT;
    $$->stmt.assignstmt = $1;
}
| varstmtlist T_COMMA T_IDENTIFIER {
    $$ = ast_varstmtlist_init();
    $$->stmttype = AST_IDENTIFIER;
    $$->varstmtlist = $1;
    $$->stmt.identifier = $3;
}
| varstmtlist T_COMMA assignstmt {
    $$ = ast_varstmtlist_init();
    $$->stmttype = AST_ASSIGNSTMT;
    $$->varstmtlist = $1;
    $$->stmt.assignstmt = $3;
}
;
%%

int yyerror(Scanner *s, AST_Program **node, const char *str)
{
    (void)*node;
    (void)u_fprintf(ustderr, "%s at %d:%d\n", str, s->linenum, s->linepos);
    return EXIT_FAILURE;
}

int yylex(YYSTYPE *yylval, Scanner *s)
{
    scanner_token(s);
    /* force re-read on comments */
    if (s->name == T_COMMENT) {
        return yylex(yylval, s);
    }
    else {
        if (s->name == T_NUMBER) {
            yylval->atom = calloc(s->ti + 1, sizeof(UChar));
            u_strcpy(yylval->atom, s->tbuf);
        }
        else if (s->name == T_IDENTIFIER || s->name == T_STRING) {
            yylval->atom = calloc(s->ti + 1, sizeof(UChar));
            u_strcpy(yylval->atom, s->tbuf);
        }
        else if (s->name == T_TRUE || s->name == T_FALSE) {
            yylval->atom = calloc(s->ti + 1, sizeof(UChar));
            u_strcpy(yylval->atom, s->tbuf);
        }
        return s->name;
    }
}
