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

ASTNode_AssignStmt *astnode_assignstmt_init(void);
void astnode_assignstmt_free(ASTNode_AssignStmt *);

ASTNode_CompareExpr *astnode_compareexpr_init(void);
void astnode_compareexpr_free(ASTNode_CompareExpr *);

ASTNode_ComplexStmt *astnode_complexstmt_init(void);
void astnode_complexstmt_free(ASTNode_ComplexStmt *);

ASTNode_CompoundBody *astnode_compoundbody_init(void);
void astnode_compoundbody_free(ASTNode_CompoundBody *);

ASTNode_CompoundBodyList *astnode_compoundbodylist_init(void);
void astnode_compoundbodylist_free(ASTNode_CompoundBodyList *);

ASTNode_CompoundStmt *astnode_compoundstmt_init(void);
void astnode_compoundstmt_free(ASTNode_CompoundStmt *);

ASTNode_ElseStmt *astnode_elsestmt_init(void);
void astnode_elsestmt_free(ASTNode_ElseStmt *);

ASTNode_Expr *astnode_expr_init(void);
void astnode_expr_free(ASTNode_Expr *);

ASTNode_ExprList *astnode_exprlist_init(void);
void astnode_exprlist_free(ASTNode_ExprList *);

ASTNode_Factor *astnode_factor_init(void);
void astnode_factor_free(ASTNode_Factor *);

ASTNode_FuncCall *astnode_funccall_init(void);
void astnode_funccall_free(ASTNode_FuncCall *);

ASTNode_FuncDef *astnode_funcdef_init(void);
void astnode_funcdef_free(ASTNode_FuncDef *);

ASTNode_FuncParamList *astnode_funcparamlist_init(void);
void astnode_funcparamlist_free(ASTNode_FuncParamList *);

ASTNode_IfStmt *astnode_ifstmt_init(void);
void astnode_ifstmt_free(ASTNode_IfStmt *);

ASTNode_MinorExpr *astnode_minorexpr_init(void);
void astnode_minorexpr_free(ASTNode_MinorExpr *);

ASTNode_NotExpr *astnode_notexpr_init(void);
void astnode_notexpr_free(ASTNode_NotExpr *);

ASTNode_Program *astnode_program_init(void);
void astnode_program_free(ASTNode_Program *);

ASTNode_ReturnStmt *astnode_returnstmt_init(void);
void astnode_returnstmt_free(ASTNode_ReturnStmt *);

ASTNode_SimpleStmt *astnode_simplestmt_init(void);
void astnode_simplestmt_free(ASTNode_SimpleStmt *);

ASTNode_Stmt *astnode_stmt_init(void);
void astnode_stmt_free(ASTNode_Stmt *);

ASTNode_StmtList *astnode_stmtlist_init(void);
void astnode_stmtlist_free(ASTNode_StmtList *);

ASTNode_Term *astnode_term_init(void);
void astnode_term_free(ASTNode_Term *);

ASTNode_VarStmt *astnode_varstmt_init(void);
void astnode_varstmt_free(ASTNode_VarStmt *);

ASTNode_VarStmtList *astnode_varstmtlist_init(void);
void astnode_varstmtlist_free(ASTNode_VarStmtList *);

ASTNode_WhileStmt *astnode_whilestmt_init(void);
void astnode_whilestmt_free(ASTNode_WhileStmt *);

struct s_ASTNode_AssignStmt 
{
    UChar *identifier;
    Token assignop;
    ASTNode_Expr *expr;
};

struct s_ASTNode_CompareExpr
{
    ASTNode_CompareExpr *compareexpr;
    Token compareop;
    ASTNode_MinorExpr *minorexpr;
};

struct s_ASTNode_ComplexStmt
{
    ASTNode_Type stmttype;
    union
    {
        ASTNode_CompoundStmt *compoundstmt;
        ASTNode_FuncDef *funcdef;
    } stmt;
};

struct s_ASTNode_CompoundBody
{
    ASTNode_CompoundBodyList *compoundbodylist;
};

struct s_ASTNode_CompoundBodyList
{
    ASTNode_CompoundBodyList *compoundbodylist;
    ASTNode_Stmt *stmt;
};

struct s_ASTNode_CompoundStmt
{
    ASTNode_Type stmttype;
    union
    {
        ASTNode_IfStmt *ifstmt;
        ASTNode_WhileStmt *whilestmt;
    } stmt;
};

struct s_ASTNode_ElseStmt
{
    ASTNode_Type stmttype;
    union
    {
        ASTNode_CompoundBody *compoundbody;
        ASTNode_IfStmt *ifstmt;
    } stmt;
};

struct s_ASTNode_Expr
{
    ASTNode_Expr *expr;
    Token exprop;
    ASTNode_NotExpr *notexpr;
};

struct s_ASTNode_ExprList
{
    ASTNode_ExprList *exprlist;
    ASTNode_Expr *expr;
};

struct s_ASTNode_Factor
{
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
    UChar* identifier;
    ASTNode_ExprList *exprlist;
};

struct s_ASTNode_FuncDef
{
    UChar* identifier;
    ASTNode_FuncParamList *funcparamlist;
    ASTNode_CompoundBody *compoundbody;
};

struct s_ASTNode_FuncParamList
{
    UChar *identifier;
    ASTNode_FuncParamList *funcparamlist;
};

struct s_ASTNode_IfStmt
{
    ASTNode_Expr *expr;
    ASTNode_CompoundBody *compoundbody;
    ASTNode_ElseStmt *elsestmt;
};

struct s_ASTNode_MinorExpr
{
    ASTNode_MinorExpr *minorexpr;
    Token addop;
    ASTNode_Term *term;
};

struct s_ASTNode_NotExpr
{
    int tnot;
    ASTNode_CompareExpr *compareexpr;
};

struct s_ASTNode_Program
{
    ASTNode_StmtList *stmtlist;
};

struct s_ASTNode_ReturnStmt
{
    ASTNode_Expr *expr;
};

struct s_ASTNode_SimpleStmt
{
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
    ASTNode_Type stmttype;
    union
    {
        ASTNode_ComplexStmt *complexstmt;
        ASTNode_SimpleStmt *simplestmt;
    } stmt;
};

struct s_ASTNode_StmtList
{
    ASTNode_StmtList *stmtlist;
    ASTNode_Stmt *stmt;
};

struct s_ASTNode_Term
{
    ASTNode_Term *term;
    Token mulop;
    ASTNode_Factor *factor;
};

struct s_ASTNode_VarStmt
{
    ASTNode_VarStmtList *varstmtlist;
};

struct s_ASTNode_VarStmtList
{
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
    ASTNode_Expr *expr;
    ASTNode_CompoundBody *compoundbody;
};

#endif
