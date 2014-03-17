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

ASTNode_AssignStmt *astnode_assignstmt_init(void)
{
    ASTNode_AssignStmt *n = calloc(1, sizeof(ASTNode_AssignStmt));
    return n;
}

void astnode_assignstmt_free(ASTNode_AssignStmt *n)
{
    free(n->identifier);
    astnode_expr_free(n->expr);
    free(n);
}

ASTNode_CompareExpr *astnode_compareexpr_init(void)
{
    ASTNode_CompareExpr *n = calloc(1, sizeof(ASTNode_CompareExpr));
    return n;
}

void astnode_compareexpr_free(ASTNode_CompareExpr *n)
{
    if (n->compareexpr != NULL) {
        astnode_compareexpr_free(n->compareexpr);
    }
    astnode_minorexpr_free(n->minorexpr);
    free(n); 
}

ASTNode_ComplexStmt *astnode_complexstmt_init(void)
{
    ASTNode_ComplexStmt *n = calloc(1, sizeof(ASTNode_ComplexStmt));
    return n;
}

void astnode_complexstmt_free(ASTNode_ComplexStmt *n)
{
    switch (n->stmttype) {
        case ASTNODE_COMPOUNDSTMT:
            astnode_compoundstmt_free(n->stmt.compoundstmt);
            break;
        case ASTNODE_FUNCDEF:
            astnode_funcdef_free(n->stmt.funcdef);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

ASTNode_CompoundBody *astnode_compoundbody_init(void)
{
    ASTNode_CompoundBody *n = calloc(1, sizeof(ASTNode_CompoundBody));
    return n;
}

void astnode_compoundbody_free(ASTNode_CompoundBody *n)
{
    astnode_compoundbodylist_free(n->compoundbodylist);
    free(n);
}

ASTNode_CompoundBodyList *astnode_compoundbodylist_init(void)
{
    ASTNode_CompoundBodyList *n = calloc(1, sizeof(ASTNode_CompoundBodyList));
    return n;
}

void astnode_compoundbodylist_free(ASTNode_CompoundBodyList *n)
{
    if (n->compoundbodylist != NULL) {
        astnode_compoundbodylist_free(n->compoundbodylist);
    }
    astnode_stmt_free(n->stmt);
    free(n);
}

ASTNode_CompoundStmt *astnode_compoundstmt_init(void)
{
    ASTNode_CompoundStmt *n = calloc(1, sizeof(ASTNode_CompoundStmt));
    return n;
}

void astnode_compoundstmt_free(ASTNode_CompoundStmt *n)
{
    switch(n->stmttype) {
        case ASTNODE_IFSTMT:
            astnode_ifstmt_free(n->stmt.ifstmt);
            break;
        case ASTNODE_WHILESTMT:
            astnode_whilestmt_free(n->stmt.whilestmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

ASTNode_ElseStmt *astnode_elsestmt_init(void)
{
    ASTNode_ElseStmt *n = calloc(1, sizeof(ASTNode_ElseStmt));
    return n;
}

void astnode_elsestmt_free(ASTNode_ElseStmt *n)
{
    switch(n->stmttype) {
        case ASTNODE_COMPOUNDBODY:
            astnode_compoundbody_free(n->stmt.compoundbody);
            break;
        case ASTNODE_IFSTMT:
            astnode_ifstmt_free(n->stmt.ifstmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

ASTNode_Expr *astnode_expr_init(void)
{
    ASTNode_Expr *n = calloc(1, sizeof(ASTNode_Expr));
    return n;
}

void astnode_expr_free(ASTNode_Expr *n)
{
    if (n->expr != NULL) {
        astnode_expr_free(n->expr);
    }
    astnode_notexpr_free(n->notexpr);
    free(n);
}

ASTNode_ExprList *astnode_exprlist_init(void)
{
    ASTNode_ExprList *n = calloc(1, sizeof(ASTNode_ExprList));
    return n;
}

void astnode_exprlist_free(ASTNode_ExprList *n)
{
    if (n->exprlist != NULL) {
        astnode_exprlist_free(n->exprlist);
    }
    astnode_expr_free(n->expr);
    free(n);
}

ASTNode_Factor *astnode_factor_init(void)
{
    ASTNode_Factor *n = calloc(1, sizeof(ASTNode_Factor));
    return n;
}

void astnode_factor_free(ASTNode_Factor *n)
{
    switch (n->factortype) {
        case ASTNODE_ATOM:
            free(n->factor.atom);
            break;
        case ASTNODE_FUNCCALL:
            astnode_funccall_free(n->factor.funccall);
            break;
        case ASTNODE_EXPR:
            astnode_expr_free(n->factor.expr);
            break;
        case ASTNODE_FACTOR:
            astnode_factor_free(n->factor.factor);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

ASTNode_FuncCall *astnode_funccall_init(void)
{
    ASTNode_FuncCall *n = calloc(1, sizeof(ASTNode_FuncCall));
    return n;
}

void astnode_funccall_free(ASTNode_FuncCall *n)
{
    free(n->identifier);
    astnode_exprlist_free(n->exprlist);
    free(n);
}

ASTNode_FuncDef *astnode_funcdef_init(void)
{
    ASTNode_FuncDef *n = calloc(1, sizeof(ASTNode_FuncDef));
    return n;
}

void astnode_funcdef_free(ASTNode_FuncDef *n)
{
    free(n->identifier);
    astnode_funcparamlist_free(n->funcparamlist);
    astnode_compoundbody_free(n->compoundbody);
    free(n);
}

ASTNode_FuncParamList *astnode_funcparamlist_init(void)
{
    ASTNode_FuncParamList *n = calloc(1, sizeof(ASTNode_FuncParamList));
    return n;
}

void astnode_funcparamlist_free(ASTNode_FuncParamList *n)
{
    free(n->identifier);
    if (n->funcparamlist != NULL) {
        astnode_funcparamlist_free(n->funcparamlist);
    }
    free(n);
}

ASTNode_IfStmt *astnode_ifstmt_init(void)
{
    ASTNode_IfStmt *n = calloc(1, sizeof(ASTNode_IfStmt));
    return n;
}

void astnode_ifstmt_free(ASTNode_IfStmt *n)
{
    astnode_expr_free(n->expr);
    astnode_compoundbody_free(n->compoundbody);
    if (n->elsestmt != NULL) {
        astnode_elsestmt_free(n->elsestmt);
    }
    free(n);
}

ASTNode_MinorExpr *astnode_minorexpr_init(void)
{
    ASTNode_MinorExpr *n = calloc(1, sizeof(ASTNode_MinorExpr));
    return n;
}

void astnode_minorexpr_free(ASTNode_MinorExpr *n)
{
    if (n->minorexpr != NULL) {
        astnode_minorexpr_free(n->minorexpr);
    }
    astnode_term_free(n->term);
    free(n);
}

ASTNode_NotExpr *astnode_notexpr_init(void)
{
    ASTNode_NotExpr *n = calloc(1, sizeof(ASTNode_NotExpr));
    return n;
}

void astnode_notexpr_free(ASTNode_NotExpr *n)
{
    astnode_compareexpr_free(n->compareexpr);
    free(n);
}

ASTNode_Program *astnode_program_init(void)
{
    ASTNode_Program *n = calloc(1, sizeof(ASTNode_Program));
    return n;
}

void astnode_program_free(ASTNode_Program *n)
{
    astnode_stmtlist_free(n->stmtlist);
    free(n);
}

ASTNode_ReturnStmt *astnode_returnstmt_init(void)
{
    ASTNode_ReturnStmt *n = calloc(1, sizeof(ASTNode_ReturnStmt));
    return n;
}

void astnode_returnstmt_free(ASTNode_ReturnStmt *n)
{
    astnode_expr_free(n->expr);
    free(n);
}

ASTNode_SimpleStmt *astnode_simplestmt_init(void)
{
    ASTNode_SimpleStmt *n = calloc(1, sizeof(ASTNode_SimpleStmt));
    return n;
}

void astnode_simplestmt_free(ASTNode_SimpleStmt *n)
{
    switch (n->stmttype) {
        case ASTNODE_ASSIGNSTMT:
            astnode_assignstmt_free(n->stmt.assignstmt);
            break;
        case ASTNODE_RETURNSTMT:
            astnode_returnstmt_free(n->stmt.returnstmt);
            break;
        case ASTNODE_VARSTMT:
            astnode_varstmt_free(n->stmt.varstmt);
            break;
        case ASTNODE_EXPR:
            astnode_expr_free(n->stmt.expr);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

ASTNode_Stmt *astnode_stmt_init(void)
{
    ASTNode_Stmt *n = calloc(1, sizeof(ASTNode_Stmt));
    return n;
}

void astnode_stmt_free(ASTNode_Stmt *n)
{
    switch (n->stmttype) {
        case ASTNODE_COMPLEXSTMT:
            astnode_complexstmt_free(n->stmt.complexstmt);
            break;
        case ASTNODE_SIMPLESTMT:
            astnode_simplestmt_free(n->stmt.simplestmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);
}

ASTNode_StmtList *astnode_stmtlist_init(void)
{
    ASTNode_StmtList *n = calloc(1, sizeof(ASTNode_StmtList));
    return n;
}

void astnode_stmtlist_free(ASTNode_StmtList *n)
{
    if (n->stmtlist != NULL) {
        astnode_stmtlist_free(n->stmtlist);
    }
    astnode_stmt_free(n->stmt);
    free(n);
}

ASTNode_Term *astnode_term_init(void)
{
    ASTNode_Term *n = calloc(1, sizeof(ASTNode_Term));
    return n;
}

void astnode_term_free(ASTNode_Term *n)
{
    if (n->term != NULL) {
        astnode_term_free(n->term);
    }
    astnode_factor_free(n->factor);
    free(n);
}

ASTNode_VarStmt *astnode_varstmt_init(void)
{
    ASTNode_VarStmt *n = calloc(1, sizeof(ASTNode_VarStmt));
    return n;
}

void astnode_varstmt_free(ASTNode_VarStmt *n)
{
    astnode_varstmtlist_free(n->varstmtlist);
    free(n);
}

ASTNode_VarStmtList *astnode_varstmtlist_init(void)
{
    ASTNode_VarStmtList *n = calloc(1, sizeof(ASTNode_VarStmtList));
    return n;
}

void astnode_varstmtlist_free(ASTNode_VarStmtList *n)
{
    if (n->varstmtlist != NULL) {
        astnode_varstmtlist_free(n->varstmtlist);
    }
    switch (n->stmttype) {
        case ASTNODE_IDENTIFIER:
            free(n->stmt.identifier);
            break;
        case ASTNODE_ASSIGNSTMT:
            astnode_assignstmt_free(n->stmt.assignstmt);
            break;
        default:
            exit(EXIT_FAILURE);
    }
    free(n);

}

ASTNode_WhileStmt *astnode_whilestmt_init(void)
{
    ASTNode_WhileStmt *n = calloc(1, sizeof(ASTNode_WhileStmt));
    return n;
}

void astnode_whilestmt_free(ASTNode_WhileStmt *n)
{
    astnode_expr_free(n->expr);
    astnode_compoundbody_free(n->compoundbody);
    free(n);
}
