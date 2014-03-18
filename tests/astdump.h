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

void astdump_assignstmt(ASTNode_AssignStmt *);
void astdump_compareexpr(ASTNode_CompareExpr *);
void astdump_complexstmt(ASTNode_ComplexStmt *);
void astdump_compoundbody(ASTNode_CompoundBody *);
void astdump_compoundbodylist(ASTNode_CompoundBodyList *);
void astdump_compoundstmt(ASTNode_CompoundStmt *);
void astdump_elsestmt(ASTNode_ElseStmt *);
void astdump_expr(ASTNode_Expr *);
void astdump_exprlist(ASTNode_ExprList *);
void astdump_factor(ASTNode_Factor *);
void astdump_funccall(ASTNode_FuncCall *);
void astdump_funcdef(ASTNode_FuncDef *);
void astdump_funcparamlist(ASTNode_FuncParamList *);
void astdump_ifstmt(ASTNode_IfStmt *);
void astdump_minorexpr(ASTNode_MinorExpr *);
void astdump_notexpr(ASTNode_NotExpr *);
void astdump_program(ASTNode_Program *);
void astdump_returnstmt(ASTNode_ReturnStmt *);
void astdump_simplestmt(ASTNode_SimpleStmt *);
void astdump_stmt(ASTNode_Stmt *);
void astdump_stmtlist(ASTNode_StmtList *);
void astdump_term(ASTNode_Term *);
void astdump_varstmt(ASTNode_VarStmt *);
void astdump_varstmtlist(ASTNode_VarStmtList *);
void astdump_whilestmt(ASTNode_WhileStmt *);

#endif
