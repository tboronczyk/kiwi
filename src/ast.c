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
#include "ast.h"

AST_AssignStmt *ast_assignstmt_init(void)
{
    AST_AssignStmt *n = calloc(1, sizeof(AST_AssignStmt));
    return n;
}

void ast_assignstmt_free(AST_AssignStmt *n)
{
    free(n->identifier);
    ast_expr_free(n->expr);
    free(n);
}

AST_CompareExpr *ast_compareexpr_init(void)
{
    AST_CompareExpr *n = calloc(1, sizeof(AST_CompareExpr));
    return n;
}

void ast_compareexpr_free(AST_CompareExpr *n)
{
    if (n->compareexpr != NULL) {
        ast_compareexpr_free(n->compareexpr);
    }
    ast_minorexpr_free(n->minorexpr);
    free(n); 
}

AST_ComplexStmt *ast_complexstmt_init(void)
{
    AST_ComplexStmt *n = calloc(1, sizeof(AST_ComplexStmt));
    return n;
}

void ast_complexstmt_free(AST_ComplexStmt *n)
{
    switch (n->stmttype) {
        case AST_COMPOUNDSTMT:
            ast_compoundstmt_free(n->stmt.compoundstmt);
            break;
        case AST_FUNCDEF:
            ast_funcdef_free(n->stmt.funcdef);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

AST_CompoundBody *ast_compoundbody_init(void)
{
    AST_CompoundBody *n = calloc(1, sizeof(AST_CompoundBody));
    return n;
}

void ast_compoundbody_free(AST_CompoundBody *n)
{
    ast_compoundbodylist_free(n->compoundbodylist);
    free(n);
}

AST_CompoundBodyList *ast_compoundbodylist_init(void)
{
    AST_CompoundBodyList *n = calloc(1, sizeof(AST_CompoundBodyList));
    return n;
}

void ast_compoundbodylist_free(AST_CompoundBodyList *n)
{
    if (n->compoundbodylist != NULL) {
        ast_compoundbodylist_free(n->compoundbodylist);
    }
    ast_stmt_free(n->stmt);
    free(n);
}

AST_CompoundStmt *ast_compoundstmt_init(void)
{
    AST_CompoundStmt *n = calloc(1, sizeof(AST_CompoundStmt));
    return n;
}

void ast_compoundstmt_free(AST_CompoundStmt *n)
{
    switch(n->stmttype) {
        case AST_IFSTMT:
            ast_ifstmt_free(n->stmt.ifstmt);
            break;
        case AST_WHILESTMT:
            ast_whilestmt_free(n->stmt.whilestmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

AST_ElseStmt *ast_elsestmt_init(void)
{
    AST_ElseStmt *n = calloc(1, sizeof(AST_ElseStmt));
    return n;
}

void ast_elsestmt_free(AST_ElseStmt *n)
{
    switch(n->stmttype) {
        case AST_COMPOUNDBODY:
            ast_compoundbody_free(n->stmt.compoundbody);
            break;
        case AST_IFSTMT:
            ast_ifstmt_free(n->stmt.ifstmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

AST_Expr *ast_expr_init(void)
{
    AST_Expr *n = calloc(1, sizeof(AST_Expr));
    return n;
}

void ast_expr_free(AST_Expr *n)
{
    if (n->expr != NULL) {
        ast_expr_free(n->expr);
    }
    ast_notexpr_free(n->notexpr);
    free(n);
}

AST_ExprList *ast_exprlist_init(void)
{
    AST_ExprList *n = calloc(1, sizeof(AST_ExprList));
    return n;
}

void ast_exprlist_free(AST_ExprList *n)
{
    if (n->exprlist != NULL) {
        ast_exprlist_free(n->exprlist);
    }
    ast_expr_free(n->expr);
    free(n);
}

AST_Factor *ast_factor_init(void)
{
    AST_Factor *n = calloc(1, sizeof(AST_Factor));
    return n;
}

void ast_factor_free(AST_Factor *n)
{
    switch (n->factortype) {
        case AST_ATOM:
            free(n->factor.atom);
            break;
        case AST_FUNCCALL:
            ast_funccall_free(n->factor.funccall);
            break;
        case AST_EXPR:
            ast_expr_free(n->factor.expr);
            break;
        case AST_FACTOR:
            ast_factor_free(n->factor.factor);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

AST_FuncCall *ast_funccall_init(void)
{
    AST_FuncCall *n = calloc(1, sizeof(AST_FuncCall));
    return n;
}

void ast_funccall_free(AST_FuncCall *n)
{
    free(n->identifier);
    ast_exprlist_free(n->exprlist);
    free(n);
}

AST_FuncDef *ast_funcdef_init(void)
{
    AST_FuncDef *n = calloc(1, sizeof(AST_FuncDef));
    return n;
}

void ast_funcdef_free(AST_FuncDef *n)
{
    free(n->identifier);
    ast_funcparamlist_free(n->funcparamlist);
    ast_compoundbody_free(n->compoundbody);
    free(n);
}

AST_FuncParamList *ast_funcparamlist_init(void)
{
    AST_FuncParamList *n = calloc(1, sizeof(AST_FuncParamList));
    return n;
}

void ast_funcparamlist_free(AST_FuncParamList *n)
{
    free(n->identifier);
    if (n->funcparamlist != NULL) {
        ast_funcparamlist_free(n->funcparamlist);
    }
    free(n);
}

AST_IfStmt *ast_ifstmt_init(void)
{
    AST_IfStmt *n = calloc(1, sizeof(AST_IfStmt));
    return n;
}

void ast_ifstmt_free(AST_IfStmt *n)
{
    ast_expr_free(n->expr);
    ast_compoundbody_free(n->compoundbody);
    if (n->elsestmt != NULL) {
        ast_elsestmt_free(n->elsestmt);
    }
    free(n);
}

AST_MinorExpr *ast_minorexpr_init(void)
{
    AST_MinorExpr *n = calloc(1, sizeof(AST_MinorExpr));
    return n;
}

void ast_minorexpr_free(AST_MinorExpr *n)
{
    if (n->minorexpr != NULL) {
        ast_minorexpr_free(n->minorexpr);
    }
    ast_term_free(n->term);
    free(n);
}

AST_NotExpr *ast_notexpr_init(void)
{
    AST_NotExpr *n = calloc(1, sizeof(AST_NotExpr));
    return n;
}

void ast_notexpr_free(AST_NotExpr *n)
{
    ast_compareexpr_free(n->compareexpr);
    free(n);
}

AST_Program *ast_program_init(void)
{
    AST_Program *n = calloc(1, sizeof(AST_Program));
    return n;
}

void ast_program_free(AST_Program *n)
{
    ast_stmtlist_free(n->stmtlist);
    free(n);
}

AST_ReturnStmt *ast_returnstmt_init(void)
{
    AST_ReturnStmt *n = calloc(1, sizeof(AST_ReturnStmt));
    return n;
}

void ast_returnstmt_free(AST_ReturnStmt *n)
{
    ast_expr_free(n->expr);
    free(n);
}

AST_SimpleStmt *ast_simplestmt_init(void)
{
    AST_SimpleStmt *n = calloc(1, sizeof(AST_SimpleStmt));
    return n;
}

void ast_simplestmt_free(AST_SimpleStmt *n)
{
    switch (n->stmttype) {
        case AST_ASSIGNSTMT:
            ast_assignstmt_free(n->stmt.assignstmt);
            break;
        case AST_RETURNSTMT:
            ast_returnstmt_free(n->stmt.returnstmt);
            break;
        case AST_VARSTMT:
            ast_varstmt_free(n->stmt.varstmt);
            break;
        case AST_EXPR:
            ast_expr_free(n->stmt.expr);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

AST_Stmt *ast_stmt_init(void)
{
    AST_Stmt *n = calloc(1, sizeof(AST_Stmt));
    return n;
}

void ast_stmt_free(AST_Stmt *n)
{
    switch (n->stmttype) {
        case AST_COMPLEXSTMT:
            ast_complexstmt_free(n->stmt.complexstmt);
            break;
        case AST_SIMPLESTMT:
            ast_simplestmt_free(n->stmt.simplestmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

AST_StmtList *ast_stmtlist_init(void)
{
    AST_StmtList *n = calloc(1, sizeof(AST_StmtList));
    return n;
}

void ast_stmtlist_free(AST_StmtList *n)
{
    if (n->stmtlist != NULL) {
        ast_stmtlist_free(n->stmtlist);
    }
    ast_stmt_free(n->stmt);
    free(n);
}

AST_Term *ast_term_init(void)
{
    AST_Term *n = calloc(1, sizeof(AST_Term));
    return n;
}

void ast_term_free(AST_Term *n)
{
    if (n->term != NULL) {
        ast_term_free(n->term);
    }
    ast_factor_free(n->factor);
    free(n);
}

AST_VarStmt *ast_varstmt_init(void)
{
    AST_VarStmt *n = calloc(1, sizeof(AST_VarStmt));
    return n;
}

void ast_varstmt_free(AST_VarStmt *n)
{
    ast_varstmtlist_free(n->varstmtlist);
    free(n);
}

AST_VarStmtList *ast_varstmtlist_init(void)
{
    AST_VarStmtList *n = calloc(1, sizeof(AST_VarStmtList));
    return n;
}

void ast_varstmtlist_free(AST_VarStmtList *n)
{
    if (n->varstmtlist != NULL) {
        ast_varstmtlist_free(n->varstmtlist);
    }
    switch (n->stmttype) {
        case AST_IDENTIFIER:
            free(n->stmt.identifier);
            break;
        case AST_ASSIGNSTMT:
            ast_assignstmt_free(n->stmt.assignstmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);

}

AST_WhileStmt *ast_whilestmt_init(void)
{
    AST_WhileStmt *n = calloc(1, sizeof(AST_WhileStmt));
    return n;
}

void ast_whilestmt_free(AST_WhileStmt *n)
{
    ast_expr_free(n->expr);
    ast_compoundbody_free(n->compoundbody);
    free(n);
}
