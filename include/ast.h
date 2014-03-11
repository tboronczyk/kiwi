#ifndef AST_H
#define AST_H

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

#include <unicode/uchar.h>

typedef int Token;

typedef struct s_ASTNode_Node ASTNode_Node;

typedef struct s_ASTNode_AssignStmt ASTNode_AssignStmt;
typedef struct s_ASTNode_CompareExpr ASTNode_CompareExpr;
typedef struct s_ASTNode_ComplexStmt ASTNode_ComplexStmt;
typedef struct s_ASTNode_CompoundBody ASTNode_CompoundBody;
typedef struct s_ASTNode_CompoundBodyList ASTNode_CompoundBodyList;
typedef struct s_ASTNode_CompoundStmt ASTNode_CompoundStmt;
typedef struct s_ASTNode_ElseStmt ASTNode_ElseStmt;
typedef struct s_ASTNode_Expr ASTNode_Expr;
typedef struct s_ASTNode_ExprList ASTNode_ExprList;
typedef struct s_ASTNode_Factor ASTNode_Factor;
typedef struct s_ASTNode_FuncCall ASTNode_FuncCall;
typedef struct s_ASTNode_FuncDef ASTNode_FuncDef;
typedef struct s_ASTNode_FuncParamList ASTNode_FuncParamList;
typedef struct s_ASTNode_IfStmt ASTNode_IfStmt;
typedef struct s_ASTNode_MinorExpr ASTNode_MinorExpr;
typedef struct s_ASTNode_NotExpr ASTNode_NotExpr;
typedef struct s_ASTNode_Program ASTNode_Program;
typedef struct s_ASTNode_ReturnStmt ASTNode_ReturnStmt;
typedef struct s_ASTNode_SimpleStmt ASTNode_SimpleStmt;
typedef struct s_ASTNode_Stmt ASTNode_Stmt;
typedef struct s_ASTNode_StmtList ASTNode_StmtList;
typedef struct s_ASTNode_Term ASTNode_Term;
typedef struct s_ASTNode_VarStmt ASTNode_VarStmt;
typedef struct s_ASTNode_VarStmtList ASTNode_VarStmtList;
typedef struct s_ASTNode_WhileStmt ASTNode_WhileStmt;

typedef enum
{
    ASTNODE_ADDOP,
    ASTNODE_ASSIGNOP,
    ASTNODE_ASSIGNSTMT,
    ASTNODE_ATOM,
    ASTNODE_COMPAREEXPR,
    ASTNODE_COMPAREOP,
    ASTNODE_COMPLEXSTMT,
    ASTNODE_COMPOUNDBODY,
    ASTNODE_COMPOUNDBODYLIST,
    ASTNODE_COMPOUNDSTMT,
    ASTNODE_ELSESTMT,
    ASTNODE_EXPR,
    ASTNODE_EXPRLIST,
    ASTNODE_EXPROP,
    ASTNODE_FACTOR,
    ASTNODE_FUNCCALL,
    ASTNODE_FUNCDEF,
    ASTNODE_FUNCPARAMLIST,
    ASTNODE_IDENTIFIER,
    ASTNODE_IFSTMT,
    ASTNODE_MINOREXPR,
    ASTNODE_MULOP,
    ASTNODE_NOTEXPR,
    ASTNODE_PROGRAM,
    ASTNODE_RETURNSTMT,
    ASTNODE_SIMPLESTMT,
    ASTNODE_STMT,
    ASTNODE_STMTLIST,
    ASTNODE_TERM,
    ASTNODE_VARSTMT,
    ASTNODE_VARSTMTLIST,
    ASTNODE_WHILESTMT
}
ASTNode_Type;

void *astnode_init(ASTNode_Type);

struct s_ASTNode_Node
{
    ASTNode_Type nodetype;
};

struct s_ASTNode_AssignStmt 
{
    ASTNode_Type nodetype;
    UChar *identifier;
    Token assignop;
    ASTNode_Expr *expr;
};

struct s_ASTNode_CompareExpr
{
    ASTNode_Type nodetype;
    ASTNode_CompareExpr *compareexpr;
    Token compareop;
    ASTNode_MinorExpr *minorexpr;
};

struct s_ASTNode_ComplexStmt
{
    ASTNode_Type nodetype;
    ASTNode_Type stmttype;
    union
    {
        ASTNode_CompoundStmt *compoundstmt;
        ASTNode_FuncDef *funcdef;
    } stmt;
};

struct s_ASTNode_CompoundBody
{
    ASTNode_Type nodetype;
    ASTNode_CompoundBodyList *compoundbodylist;
};

struct s_ASTNode_CompoundBodyList
{
    ASTNode_Type nodetype;
    ASTNode_CompoundBodyList *compoundbodylist;
    ASTNode_Stmt *stmt;
};

struct s_ASTNode_CompoundStmt
{
    ASTNode_Type nodetype;
    ASTNode_Type stmttype;
    union
    {
        ASTNode_IfStmt *ifstmt;
        ASTNode_WhileStmt *whilestmt;
    } stmt;
};

struct s_ASTNode_ElseStmt
{
    ASTNode_Type nodetype;
    ASTNode_Type stmttype;
    union
    {
        ASTNode_CompoundBody *compoundbody;
        ASTNode_IfStmt *ifstmt;
    } stmt;
};

struct s_ASTNode_Expr
{
    ASTNode_Type nodetype;
    ASTNode_Expr *expr;
    Token exprop;
    ASTNode_NotExpr *notexpr;
};

struct s_ASTNode_ExprList
{
    ASTNode_Type nodetype;
    ASTNode_ExprList *exprlist;
    ASTNode_Expr *expr;
};

struct s_ASTNode_Factor
{
    ASTNode_Type nodetype;
    ASTNode_Type factortype;
    union
    {
        UChar *atom;
        ASTNode_FuncCall *funccall;
        ASTNode_Expr *expr;
        ASTNode_Factor *factor;
    } factor;
    Token addop;
};

struct s_ASTNode_FuncCall
{
    ASTNode_Type nodetype;
    UChar* identifier;
    ASTNode_ExprList *exprlist;
};

struct s_ASTNode_FuncDef
{
    ASTNode_Type nodetype;
    UChar* identifier;
    ASTNode_FuncParamList *funcparamlist;
    ASTNode_CompoundBody *compoundbody;
};

struct s_ASTNode_FuncParamList
{
    ASTNode_Type nodetype;
    UChar *identifier;
    ASTNode_FuncParamList *funcparamlist;
};

struct s_ASTNode_IfStmt
{
    ASTNode_Type nodetype;
    ASTNode_Expr *expr;
    ASTNode_CompoundBody *compoundbody;
    ASTNode_ElseStmt *elsestmt;
};

struct s_ASTNode_MinorExpr
{
    ASTNode_Type nodetype;
    ASTNode_MinorExpr *minorexpr;
    Token addop;
    ASTNode_Term *term;
};

struct s_ASTNode_NotExpr
{
    ASTNode_Type nodetype;
    int tnot;
    ASTNode_CompareExpr *compareexpr;
};

struct s_ASTNode_Program
{
    ASTNode_Type nodetype;
    ASTNode_StmtList *stmtlist;
};

struct s_ASTNode_ReturnStmt
{
    ASTNode_Type nodetype;
    ASTNode_Expr *expr;
};

struct s_ASTNode_SimpleStmt
{
    ASTNode_Type nodetype;
    ASTNode_Type stmttype;
    union
    {
        ASTNode_AssignStmt *assignstmt;
        ASTNode_ReturnStmt *returnstmt;
        ASTNode_VarStmt *varstmt;
        ASTNode_Expr * expr;
    } stmt;
};

struct s_ASTNode_Stmt
{
    ASTNode_Type nodetype;
    ASTNode_Type stmttype;
    union
    {
        ASTNode_ComplexStmt *complexstmt;
        ASTNode_SimpleStmt *simplestmt;
    } stmt;
};

struct s_ASTNode_StmtList
{
    ASTNode_Type nodetype;
    ASTNode_StmtList *stmtlist;
    ASTNode_Stmt *stmt;
};

struct s_ASTNode_Term
{
    ASTNode_Type nodetype;
    ASTNode_Term *term;
    Token mulop;
    ASTNode_Factor *factor;
};

struct s_ASTNode_VarStmt
{
    ASTNode_Type nodetype;
    ASTNode_VarStmtList *varstmtlist;
};

struct s_ASTNode_VarStmtList
{
    ASTNode_Type nodetype;
    ASTNode_VarStmtList *varstmtlist;
    ASTNode_Type stmttype;
    union
    {
        UChar *identifier;
        ASTNode_AssignStmt *assignstmt;
    } stmt;
};

struct s_ASTNode_WhileStmt
{
    ASTNode_Type nodetype;
    ASTNode_Expr *expr;
    ASTNode_CompoundBody *compoundbody;
};

#endif
