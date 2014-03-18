#ifndef ASTDUMP_H
#define ASTDUMP_H

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

#include "ast.h"

void astdump_assignstmt(AST_AssignStmt *);
void astdump_compareexpr(AST_CompareExpr *);
void astdump_complexstmt(AST_ComplexStmt *);
void astdump_compoundbody(AST_CompoundBody *);
void astdump_compoundbodylist(AST_CompoundBodyList *);
void astdump_compoundstmt(AST_CompoundStmt *);
void astdump_elsestmt(AST_ElseStmt *);
void astdump_expr(AST_Expr *);
void astdump_exprlist(AST_ExprList *);
void astdump_factor(AST_Factor *);
void astdump_funccall(AST_FuncCall *);
void astdump_funcdef(AST_FuncDef *);
void astdump_funcparamlist(AST_FuncParamList *);
void astdump_ifstmt(AST_IfStmt *);
void astdump_minorexpr(AST_MinorExpr *);
void astdump_notexpr(AST_NotExpr *);
void astdump_program(AST_Program *);
void astdump_returnstmt(AST_ReturnStmt *);
void astdump_simplestmt(AST_SimpleStmt *);
void astdump_stmt(AST_Stmt *);
void astdump_stmtlist(AST_StmtList *);
void astdump_term(AST_Term *);
void astdump_varstmt(AST_VarStmt *);
void astdump_varstmtlist(AST_VarStmtList *);
void astdump_whilestmt(AST_WhileStmt *);

#endif
