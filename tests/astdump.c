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

#include <stdlib.h>
#include <unicode/ustdio.h>
#include "ast.h"
#include "astdump.h"
#include "scanner.h"
#include "parser.h"

extern UFILE *ustdout;

void astdump_assignstmt(AST_AssignStmt *n)
{
    u_fprintf(ustdout, "%S ", n->identifier);
    switch (n->assignop) {
        case T_ASSIGN:
            u_fprintf(ustdout, ":= ");
            break;
        case T_ADD_ASSIGN:
            u_fprintf(ustdout, "+: ");
            break;
        case T_SUBTRACT_ASSIGN:
            u_fprintf(ustdout, "-: ");
            break;
        case T_MULTIPLY_ASSIGN:
            u_fprintf(ustdout, "*: ");
            break;
        case T_DIVIDE_ASSIGN:
            u_fprintf(ustdout, "/: ");
            break;
        case T_MODULO_ASSIGN:
            u_fprintf(ustdout, "%%: ");
            break;
        default:
            exit(EXIT_FAILURE);
    }
    astdump_expr(n->expr);
}

void astdump_compareexpr(AST_CompareExpr *n)
{
    if (n->compareexpr != NULL) {
        astdump_compareexpr(n->compareexpr);
        switch (n->compareop) {
            case T_EQUAL:
                u_fprintf(ustdout, "= ");
                break;
            case T_NOT_EQUAL:
                u_fprintf(ustdout, "~= ");
                break;
            case T_LESS:
                u_fprintf(ustdout, "< ");
                break;
            case T_LESS_EQUAL:
                u_fprintf(ustdout, "<= ");
                break;
            case T_GREATER:
                u_fprintf(ustdout, "> ");
                break;
            case T_GREATER_EQUAL:
                u_fprintf(ustdout, ">= ");
                break;
            case T_IS:
                u_fprintf(ustdout, "IS ");
                break;
            default:
                exit(EXIT_FAILURE);
        }
    }
    astdump_minorexpr(n->minorexpr);
}

void astdump_complexstmt(AST_ComplexStmt *n)
{
    switch (n->stmttype) {
        case AST_COMPOUNDSTMT:
            astdump_compoundstmt(n->stmt.compoundstmt);
            break;
        case AST_FUNCDEF: 
            astdump_funcdef(n->stmt.funcdef);
            break;
        default:
            exit(EXIT_FAILURE);
    }
}

void astdump_compoundbody(AST_CompoundBody *n)
{
    u_fprintf(ustdout, "{ ");
    astdump_compoundbodylist(n->compoundbodylist);
    u_fprintf(ustdout, "} ");
}

void astdump_compoundbodylist(AST_CompoundBodyList *n)
{
    if (n->compoundbodylist != NULL) {
        astdump_compoundbodylist(n->compoundbodylist);
    }
    astdump_stmt(n->stmt);
}

void astdump_compoundstmt(AST_CompoundStmt *n)
{
    switch (n->stmttype) {
        case AST_IFSTMT:
            astdump_ifstmt(n->stmt.ifstmt);
            break;
        case AST_WHILESTMT:
            astdump_whilestmt(n->stmt.whilestmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
}

void astdump_elsestmt(AST_ElseStmt *n)
{
    u_fprintf(ustdout, "ELSE ");
    switch (n->stmttype) {
        case AST_COMPOUNDBODY:
            astdump_compoundbody(n->stmt.compoundbody);
            break;
        case AST_IFSTMT:
            astdump_ifstmt(n->stmt.ifstmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
}

void astdump_expr(AST_Expr *n)
{
    if (n->expr != NULL) {
        astdump_expr(n->expr);
        switch (n->exprop) {
            case T_AND:
                u_fprintf(ustdout, "AND ");
                break;
            case T_OR:
                u_fprintf(ustdout, "OR ");
                break;
            default:
                exit(EXIT_FAILURE);
        }
    }
    astdump_notexpr(n->notexpr);
}

void astdump_exprlist(AST_ExprList *n)
{
    if (n->exprlist != NULL) {
        astdump_exprlist(n->exprlist);
        u_fprintf(ustdout, ", ");
    }
    astdump_expr(n->expr);
}

void astdump_factor(AST_Factor *n)
{
    switch (n->factortype) {
        case AST_ATOM:
            u_fprintf(ustdout, "%S ", n->factor.atom);
            break;
        case AST_FUNCCALL:
            astdump_funccall(n->factor.funccall);
            break;
        case AST_EXPR:
            u_fprintf(ustdout, "( ");
            astdump_expr(n->factor.expr);
            u_fprintf(ustdout, ") ");
            break;
        case AST_FACTOR:
            switch (n->addop) {
                case T_ADD:
                    u_fprintf(ustdout, "+ ");
                    break;
                case T_SUBTRACT:
                    u_fprintf(ustdout, "- ");
                    break;
                default:
                    exit(EXIT_FAILURE);
            }
            astdump_factor(n->factor.factor);
            break;
        default:
            exit(EXIT_FAILURE);
    }
}

void astdump_funccall(AST_FuncCall *n)
{
     u_fprintf(ustdout, "%S ( ", n->identifier);
     if (n->exprlist != NULL) {
         astdump_exprlist(n->exprlist);
     }
     u_fprintf(ustdout, ") ");
}

void astdump_funcdef(AST_FuncDef *n)
{
     u_fprintf(ustdout, "FUNC %S ", n->identifier);
     if (n->funcparamlist != NULL) {
         astdump_funcparamlist(n->funcparamlist);
     }
     astdump_compoundbody(n->compoundbody);
}

void astdump_funcparamlist(AST_FuncParamList *n)
{
    if (n->funcparamlist != NULL) {
        astdump_funcparamlist(n->funcparamlist);
        u_fprintf(ustdout, ", ");
    }
    u_fprintf(ustdout, "%S ", n->identifier);
}

void astdump_ifstmt(AST_IfStmt *n)
{
    u_fprintf(ustdout, "IF ");
    astdump_expr(n->expr);
    astdump_compoundbody(n->compoundbody);
    if (n->elsestmt != NULL) {
        astdump_elsestmt(n->elsestmt);
    }
}

void astdump_minorexpr(AST_MinorExpr *n)
{
    if (n->minorexpr != NULL) {
        astdump_minorexpr(n->minorexpr);
        switch (n->addop) {
            case T_ADD:
                u_fprintf(ustdout, "+ ");
                break;
            case T_SUBTRACT:
                u_fprintf(ustdout, "- ");
                break;
            default:
                exit(EXIT_FAILURE);
        }
    }
    astdump_term(n->term);
}

void astdump_notexpr(AST_NotExpr *n)
{
    if (n->tnot) {
        u_fprintf(ustdout, "~ ");
    }
    astdump_compareexpr(n->compareexpr);
}

void astdump_program(AST_Program *n)
{
    if (n->stmtlist != NULL) {
        astdump_stmtlist(n->stmtlist);
    }
}

void astdump_returnstmt(AST_ReturnStmt *n)
{
   u_fprintf(ustdout, "RETURN ");
   astdump_expr(n->expr);
}

void astdump_simplestmt(AST_SimpleStmt *n)
{
    switch (n->stmttype) {
        case AST_ASSIGNSTMT:
            astdump_assignstmt(n->stmt.assignstmt);
            break;
        case AST_RETURNSTMT:
            astdump_returnstmt(n->stmt.returnstmt);
            break;
        case AST_VARSTMT:
            astdump_varstmt(n->stmt.varstmt);
            break;
        case AST_EXPR:
            astdump_expr(n->stmt.expr);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    u_fprintf(ustdout, ". ");
}

void astdump_stmt(AST_Stmt *n)
{
    switch (n->stmttype) {
        case AST_COMPLEXSTMT:
            astdump_complexstmt(n->stmt.complexstmt);
            break;
        case AST_SIMPLESTMT:
            astdump_simplestmt(n->stmt.simplestmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
}

void astdump_stmtlist(AST_StmtList *n)
{
    if (n->stmtlist != NULL) {
        astdump_stmtlist(n->stmtlist);
    }
    astdump_stmt(n->stmt);
}

void astdump_term(AST_Term *n)
{
    if (n->term != NULL) {
        astdump_term(n->term);
        switch (n->mulop) {
            case T_MULTIPLY:
                u_fprintf(ustdout, "* ");
                break;
            case T_DIVIDE:
                u_fprintf(ustdout, "/ ");
                break;
            case T_MODULO:
                u_fprintf(ustdout, "%% ");
                break;
            default:
                exit(EXIT_FAILURE);
        }
    }
    astdump_factor(n->factor);
}

void astdump_varstmt(AST_VarStmt *n)
{
    u_fprintf(ustdout, "VAR ");
    astdump_varstmtlist(n->varstmtlist);
}

void astdump_varstmtlist(AST_VarStmtList *n)
{
    if (n->varstmtlist != NULL) {
        astdump_varstmtlist(n->varstmtlist);
        u_fprintf(ustdout, ", ");
    }
    switch (n->stmttype) {
        case AST_IDENTIFIER:
           u_fprintf(ustdout, "%S ", n->stmt.identifier);
           break;
        case AST_ASSIGNSTMT:
           astdump_assignstmt(n->stmt.assignstmt);
           break;
        default:
            exit(EXIT_FAILURE);
    }
}

void astdump_whilestmt(AST_WhileStmt *n)
{
    u_fprintf(ustdout, "WHILE ");
    astdump_expr(n->expr);
    astdump_compoundbody(n->compoundbody);
}

