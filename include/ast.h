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

typedef struct s_AST_AssignStmt AST_AssignStmt;
typedef struct s_AST_CompareExpr AST_CompareExpr;
typedef struct s_AST_ComplexStmt AST_ComplexStmt;
typedef struct s_AST_CompoundBody AST_CompoundBody;
typedef struct s_AST_CompoundBodyList AST_CompoundBodyList;
typedef struct s_AST_CompoundStmt AST_CompoundStmt;
typedef struct s_AST_ElseStmt AST_ElseStmt;
typedef struct s_AST_Expr AST_Expr;
typedef struct s_AST_ExprList AST_ExprList;
typedef struct s_AST_Factor AST_Factor;
typedef struct s_AST_FuncCall AST_FuncCall;
typedef struct s_AST_FuncDef AST_FuncDef;
typedef struct s_AST_FuncParamList AST_FuncParamList;
typedef struct s_AST_IfStmt AST_IfStmt;
typedef struct s_AST_MinorExpr AST_MinorExpr;
typedef struct s_AST_NotExpr AST_NotExpr;
typedef struct s_AST_Program AST_Program;
typedef struct s_AST_ReturnStmt AST_ReturnStmt;
typedef struct s_AST_SimpleStmt AST_SimpleStmt;
typedef struct s_AST_Stmt AST_Stmt;
typedef struct s_AST_StmtList AST_StmtList;
typedef struct s_AST_Term AST_Term;
typedef struct s_AST_VarStmt AST_VarStmt;
typedef struct s_AST_VarStmtList AST_VarStmtList;
typedef struct s_AST_WhileStmt AST_WhileStmt;

typedef enum
{
    AST_ADDOP,
    AST_ASSIGNOP,
    AST_ASSIGNSTMT,
    AST_ATOM,
    AST_COMPAREEXPR,
    AST_COMPAREOP,
    AST_COMPLEXSTMT,
    AST_COMPOUNDBODY,
    AST_COMPOUNDBODYLIST,
    AST_COMPOUNDSTMT,
    AST_ELSESTMT,
    AST_EXPR,
    AST_EXPRLIST,
    AST_EXPROP,
    AST_FACTOR,
    AST_FUNCCALL,
    AST_FUNCDEF,
    AST_FUNCPARAMLIST,
    AST_IDENTIFIER,
    AST_IFSTMT,
    AST_MINOREXPR,
    AST_MULOP,
    AST_NOTEXPR,
    AST_PROGRAM,
    AST_RETURNSTMT,
    AST_SIMPLESTMT,
    AST_STMT,
    AST_STMTLIST,
    AST_TERM,
    AST_VARSTMT,
    AST_VARSTMTLIST,
    AST_WHILESTMT
}
AST_Type;

AST_AssignStmt *ast_assignstmt_init(void);
void ast_assignstmt_free(AST_AssignStmt *);

AST_CompareExpr *ast_compareexpr_init(void);
void ast_compareexpr_free(AST_CompareExpr *);

AST_ComplexStmt *ast_complexstmt_init(void);
void ast_complexstmt_free(AST_ComplexStmt *);

AST_CompoundBody *ast_compoundbody_init(void);
void ast_compoundbody_free(AST_CompoundBody *);

AST_CompoundBodyList *ast_compoundbodylist_init(void);
void ast_compoundbodylist_free(AST_CompoundBodyList *);

AST_CompoundStmt *ast_compoundstmt_init(void);
void ast_compoundstmt_free(AST_CompoundStmt *);

AST_ElseStmt *ast_elsestmt_init(void);
void ast_elsestmt_free(AST_ElseStmt *);

AST_Expr *ast_expr_init(void);
void ast_expr_free(AST_Expr *);

AST_ExprList *ast_exprlist_init(void);
void ast_exprlist_free(AST_ExprList *);

AST_Factor *ast_factor_init(void);
void ast_factor_free(AST_Factor *);

AST_FuncCall *ast_funccall_init(void);
void ast_funccall_free(AST_FuncCall *);

AST_FuncDef *ast_funcdef_init(void);
void ast_funcdef_free(AST_FuncDef *);

AST_FuncParamList *ast_funcparamlist_init(void);
void ast_funcparamlist_free(AST_FuncParamList *);

AST_IfStmt *ast_ifstmt_init(void);
void ast_ifstmt_free(AST_IfStmt *);

AST_MinorExpr *ast_minorexpr_init(void);
void ast_minorexpr_free(AST_MinorExpr *);

AST_NotExpr *ast_notexpr_init(void);
void ast_notexpr_free(AST_NotExpr *);

AST_Program *ast_program_init(void);
void ast_program_free(AST_Program *);

AST_ReturnStmt *ast_returnstmt_init(void);
void ast_returnstmt_free(AST_ReturnStmt *);

AST_SimpleStmt *ast_simplestmt_init(void);
void ast_simplestmt_free(AST_SimpleStmt *);

AST_Stmt *ast_stmt_init(void);
void ast_stmt_free(AST_Stmt *);

AST_StmtList *ast_stmtlist_init(void);
void ast_stmtlist_free(AST_StmtList *);

AST_Term *ast_term_init(void);
void ast_term_free(AST_Term *);

AST_VarStmt *ast_varstmt_init(void);
void ast_varstmt_free(AST_VarStmt *);

AST_VarStmtList *ast_varstmtlist_init(void);
void ast_varstmtlist_free(AST_VarStmtList *);

AST_WhileStmt *ast_whilestmt_init(void);
void ast_whilestmt_free(AST_WhileStmt *);

struct s_AST_AssignStmt 
{
    UChar *identifier;
    Token assignop;
    AST_Expr *expr;
};

struct s_AST_CompareExpr
{
    AST_CompareExpr *compareexpr;
    Token compareop;
    AST_MinorExpr *minorexpr;
};

struct s_AST_ComplexStmt
{
    AST_Type stmttype;
    union
    {
        AST_CompoundStmt *compoundstmt;
        AST_FuncDef *funcdef;
    } stmt;
};

struct s_AST_CompoundBody
{
    AST_CompoundBodyList *compoundbodylist;
};

struct s_AST_CompoundBodyList
{
    AST_CompoundBodyList *compoundbodylist;
    AST_Stmt *stmt;
};

struct s_AST_CompoundStmt
{
    AST_Type stmttype;
    union
    {
        AST_IfStmt *ifstmt;
        AST_WhileStmt *whilestmt;
    } stmt;
};

struct s_AST_ElseStmt
{
    AST_Type stmttype;
    union
    {
        AST_CompoundBody *compoundbody;
        AST_IfStmt *ifstmt;
    } stmt;
};

struct s_AST_Expr
{
    AST_Expr *expr;
    Token exprop;
    AST_NotExpr *notexpr;
};

struct s_AST_ExprList
{
    AST_ExprList *exprlist;
    AST_Expr *expr;
};

struct s_AST_Factor
{
    AST_Type factortype;
    union
    {
        UChar *atom;
        AST_FuncCall *funccall;
        AST_Expr *expr;
        AST_Factor *factor;
    } factor;
    Token addop;
};

struct s_AST_FuncCall
{
    UChar* identifier;
    AST_ExprList *exprlist;
};

struct s_AST_FuncDef
{
    UChar* identifier;
    AST_FuncParamList *funcparamlist;
    AST_CompoundBody *compoundbody;
};

struct s_AST_FuncParamList
{
    UChar *identifier;
    AST_FuncParamList *funcparamlist;
};

struct s_AST_IfStmt
{
    AST_Expr *expr;
    AST_CompoundBody *compoundbody;
    AST_ElseStmt *elsestmt;
};

struct s_AST_MinorExpr
{
    AST_MinorExpr *minorexpr;
    Token addop;
    AST_Term *term;
};

struct s_AST_NotExpr
{
    int tnot;
    AST_CompareExpr *compareexpr;
};

struct s_AST_Program
{
    AST_StmtList *stmtlist;
};

struct s_AST_ReturnStmt
{
    AST_Expr *expr;
};

struct s_AST_SimpleStmt
{
    AST_Type stmttype;
    union
    {
        AST_AssignStmt *assignstmt;
        AST_ReturnStmt *returnstmt;
        AST_VarStmt *varstmt;
        AST_Expr * expr;
    } stmt;
};

struct s_AST_Stmt
{
    AST_Type stmttype;
    union
    {
        AST_ComplexStmt *complexstmt;
        AST_SimpleStmt *simplestmt;
    } stmt;
};

struct s_AST_StmtList
{
    AST_StmtList *stmtlist;
    AST_Stmt *stmt;
};

struct s_AST_Term
{
    AST_Term *term;
    Token mulop;
    AST_Factor *factor;
};

struct s_AST_VarStmt
{
    AST_VarStmtList *varstmtlist;
};

struct s_AST_VarStmtList
{
    AST_VarStmtList *varstmtlist;
    AST_Type stmttype;
    union
    {
        UChar *identifier;
        AST_AssignStmt *assignstmt;
    } stmt;
};

struct s_AST_WhileStmt
{
    AST_Expr *expr;
    AST_CompoundBody *compoundbody;
};

#endif
