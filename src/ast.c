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

static void astnode_assignstmt_free(ASTNode_AssignStmt *);
static void astnode_compareexpr_free(ASTNode_CompareExpr *);
static void astnode_complexstmt_free(ASTNode_ComplexStmt *);
static void astnode_compoundbody_free(ASTNode_CompoundBody *);
static void astnode_compoundbodylist_free(ASTNode_CompoundBodyList *);
static void astnode_compoundstmt_free(ASTNode_CompoundStmt *);
static void astnode_elsestmt_free(ASTNode_ElseStmt *);
static void astnode_expr_free(ASTNode_Expr *);
static void astnode_exprlist_free(ASTNode_ExprList *);
static void astnode_factor_free(ASTNode_Factor *);
static void astnode_funccall_free(ASTNode_FuncCall *);
static void astnode_funcdef_free(ASTNode_FuncDef *);
static void astnode_funcparamlist_free(ASTNode_FuncParamList *);
static void astnode_ifstmt_free(ASTNode_IfStmt *);
static void astnode_minorexpr_free(ASTNode_MinorExpr *);
static void astnode_notexpr_free(ASTNode_NotExpr *);
static void astnode_program_free(ASTNode_Program *);
static void astnode_returnstmt_free(ASTNode_ReturnStmt *);
static void astnode_simplestmt_free(ASTNode_SimpleStmt *);
static void astnode_stmt_free(ASTNode_Stmt *);
static void astnode_stmtlist_free(ASTNode_StmtList *);
static void astnode_term_free(ASTNode_Term *);
static void astnode_varstmt_free(ASTNode_VarStmt *);
static void astnode_varstmtlist_free(ASTNode_VarStmtList *);
static void astnode_whilestmt_free(ASTNode_WhileStmt *);

void *astnode_init(ASTNode_Type type)
{
    /* determine size to allocate for node */
    size_t size;
    switch (type) {
        case ASTNODE_ASSIGNSTMT:
            size = sizeof(ASTNode_AssignStmt);
            break;
        case ASTNODE_COMPAREEXPR:
            size = sizeof(ASTNode_CompareExpr);
            break;
        case ASTNODE_COMPLEXSTMT:
            size = sizeof(ASTNode_ComplexStmt);
            break;
        case ASTNODE_COMPOUNDBODYLIST:
            size = sizeof(ASTNode_CompoundBodyList);
            break;
        case ASTNODE_COMPOUNDBODY:
            size = sizeof(ASTNode_CompoundBody);
            break;
        case ASTNODE_COMPOUNDSTMT:
            size = sizeof(ASTNode_CompoundStmt);
            break;
        case ASTNODE_ELSESTMT:
            size = sizeof(ASTNode_ElseStmt);
            break;
        case ASTNODE_EXPRLIST:
            size = sizeof(ASTNode_ExprList);
            break;
        case ASTNODE_EXPR:
            size = sizeof(ASTNode_Expr);
            break;
        case ASTNODE_FACTOR:
            size = sizeof(ASTNode_Factor);
            break;
        case ASTNODE_FUNCCALL:
            size = sizeof(ASTNode_FuncCall);
            break;
        case ASTNODE_FUNCDEF:
            size = sizeof(ASTNode_FuncDef);
            break;
        case ASTNODE_FUNCPARAMLIST:
            size = sizeof(ASTNode_FuncParamList);
            break;
        case ASTNODE_IFSTMT:
            size = sizeof(ASTNode_IfStmt);
            break;
        case ASTNODE_MINOREXPR:
            size = sizeof(ASTNode_MinorExpr);
            break;
        case ASTNODE_NOTEXPR:
            size = sizeof(ASTNode_NotExpr);
            break;
        case ASTNODE_PROGRAM:
            size = sizeof(ASTNode_Program);
            break;
        case ASTNODE_RETURNSTMT:
            size = sizeof(ASTNode_ReturnStmt);
            break;
        case ASTNODE_SIMPLESTMT:
            size = sizeof(ASTNode_SimpleStmt);
            break;
        case ASTNODE_STMTLIST:
            size = sizeof(ASTNode_StmtList);
            break;
        case ASTNODE_STMT:
            size = sizeof(ASTNode_Stmt);
            break;
        case ASTNODE_TERM:
            size = sizeof(ASTNode_Term);
            break;
        case ASTNODE_VARSTMTLIST:
            size = sizeof(ASTNode_VarStmtList);
            break;
        case ASTNODE_VARSTMT:
            size = sizeof(ASTNode_VarStmt);
            break;
        case ASTNODE_WHILESTMT:
            size = sizeof(ASTNode_WhileStmt);
            break;
        /* should never reach this */
        default:
            perror("Invalid ASTNode_Type passed to astnode_init");
            exit(EXIT_FAILURE);
    }

    void *n;
    if ((n = calloc(1, size)) == NULL) {
        perror("Allocate astnode failed");
        exit(EXIT_FAILURE);
    }
    ((ASTNode_Node *)n)->nodetype = type;

    return n;
}

void astnode_free(ASTNode_Node *n)
{
    /* free node with appropriate function */
    switch (n->nodetype) {
        case ASTNODE_ASSIGNSTMT:
            astnode_assignstmt_free((ASTNode_AssignStmt *)n);
            break;
        case ASTNODE_COMPAREEXPR:
            astnode_compareexpr_free((ASTNode_CompareExpr *)n);
            break;
        case ASTNODE_COMPLEXSTMT:
            astnode_complexstmt_free((ASTNode_ComplexStmt *)n);
            break;
        case ASTNODE_COMPOUNDBODY:
            astnode_compoundbody_free((ASTNode_CompoundBody *)n);
            break;
        case ASTNODE_COMPOUNDBODYLIST:
            astnode_compoundbodylist_free((ASTNode_CompoundBodyList *)n);
            break;
        case ASTNODE_COMPOUNDSTMT:
            astnode_compoundstmt_free((ASTNode_CompoundStmt *)n);
            break;
        case ASTNODE_ELSESTMT:
            astnode_elsestmt_free((ASTNode_ElseStmt *)n);
            break;
        case ASTNODE_EXPR:
            astnode_expr_free((ASTNode_Expr *)n);
            break;
        case ASTNODE_EXPRLIST:
            astnode_exprlist_free((ASTNode_ExprList *)n);
            break;
        case ASTNODE_FACTOR:
            astnode_factor_free((ASTNode_Factor *)n);
            break;
        case ASTNODE_FUNCCALL:
            astnode_funccall_free((ASTNode_FuncCall *)n);
            break;
        case ASTNODE_FUNCDEF:
            astnode_funcdef_free((ASTNode_FuncDef *)n);
            break;
        case ASTNODE_FUNCPARAMLIST:
            astnode_funcparamlist_free((ASTNode_FuncParamList *)n);
            break;
        case ASTNODE_IFSTMT:
            astnode_ifstmt_free((ASTNode_IfStmt *)n);
            break;
        case ASTNODE_MINOREXPR:
            astnode_minorexpr_free((ASTNode_MinorExpr *)n);
            break;
        case ASTNODE_NOTEXPR:
            astnode_notexpr_free((ASTNode_NotExpr *)n);
            break;
        case ASTNODE_PROGRAM:
            astnode_program_free((ASTNode_Program *)n);
            break;
        case ASTNODE_RETURNSTMT:
            astnode_returnstmt_free((ASTNode_ReturnStmt *)n);
            break;
        case ASTNODE_SIMPLESTMT:
            astnode_simplestmt_free((ASTNode_SimpleStmt *)n);
            break;
        case ASTNODE_STMT:
            astnode_stmt_free((ASTNode_Stmt *)n);
            break;
        case ASTNODE_STMTLIST:
            astnode_stmtlist_free((ASTNode_StmtList *)n);
            break;
        case ASTNODE_TERM:
            astnode_term_free((ASTNode_Term *)n);
            break;
        case ASTNODE_VARSTMT:
            astnode_varstmt_free((ASTNode_VarStmt *)n);
            break;
        case ASTNODE_VARSTMTLIST:
            astnode_varstmtlist_free((ASTNode_VarStmtList *)n);
            break;
        case ASTNODE_WHILESTMT:
            astnode_whilestmt_free((ASTNode_WhileStmt *)n);
            break;
        /* should never reach this */
        default:
            perror("Invalid astnode type passed on free");
            exit(EXIT_FAILURE);
    }
}

static void astnode_assignstmt_free(ASTNode_AssignStmt *n)
{
}

static void astnode_compareexpr_free(ASTNode_CompareExpr *n)
{
}

static void astnode_complexstmt_free(ASTNode_ComplexStmt *n)
{
}

static void astnode_compoundbody_free(ASTNode_CompoundBody *n)
{
}

static void astnode_compoundbodylist_free(ASTNode_CompoundBodyList *n)
{
}

static void astnode_compoundstmt_free(ASTNode_CompoundStmt *n)
{
}

static void astnode_elsestmt_free(ASTNode_ElseStmt *n)
{
}

static void astnode_expr_free(ASTNode_Expr *n)
{
}

static void astnode_exprlist_free(ASTNode_ExprList *n)
{
}

static void astnode_factor_free(ASTNode_Factor *n)
{
}

static void astnode_funccall_free(ASTNode_FuncCall *n)
{
}

static void astnode_funcdef_free(ASTNode_FuncDef *n)
{
}

static void astnode_funcparamlist_free(ASTNode_FuncParamList *n)
{
}

static void astnode_ifstmt_free(ASTNode_IfStmt *n)
{
}

static void astnode_minorexpr_free(ASTNode_MinorExpr *n)
{
}

static void astnode_notexpr_free(ASTNode_NotExpr *n)
{
}

static void astnode_program_free(ASTNode_Program *n)
{
}

static void astnode_returnstmt_free(ASTNode_ReturnStmt *n)
{
}

static void astnode_simplestmt_free(ASTNode_SimpleStmt *n)
{
}

static void astnode_stmt_free(ASTNode_Stmt *n)
{
}

static void astnode_stmtlist_free(ASTNode_StmtList *n)
{
}

static void astnode_term_free(ASTNode_Term *n)
{
}

static void astnode_varstmt_free(ASTNode_VarStmt *n)
{
}

static void astnode_varstmtlist_free(ASTNode_VarStmtList *n)
{
}

static void astnode_whilestmt_free(ASTNode_WhileStmt *n)
{
}
