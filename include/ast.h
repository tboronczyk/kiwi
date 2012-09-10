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

typedef int token_t;

typedef struct _astnode_node astnode_node_t;

typedef struct _astnode_assignstmt astnode_assignstmt_t;
typedef struct _astnode_atom astnode_atom_t;
typedef struct _astnode_compareexpr astnode_compareexpr_t;
typedef struct _astnode_complexstmt astnode_complexstmt_t;
typedef struct _astnode_compoundbody astnode_compoundbody_t;
typedef struct _astnode_compoundbodylist astnode_compoundbodylist_t;
typedef struct _astnode_compoundstmt astnode_compoundstmt_t;
typedef struct _astnode_elsestmt astnode_elsestmt_t;
typedef struct _astnode_expr astnode_expr_t;
typedef struct _astnode_exprlist astnode_exprlist_t;
typedef struct _astnode_factor astnode_factor_t;
typedef struct _astnode_funccall astnode_funccall_t;
typedef struct _astnode_funcdef astnode_funcdef_t;
typedef struct _astnode_funcparamlist astnode_funcparamlist_t;
typedef struct _astnode_ifstmt astnode_ifstmt_t;
typedef struct _astnode_minorexpr astnode_minorexpr_t;
typedef struct _astnode_notexpr astnode_notexpr_t;
typedef struct _astnode_program astnode_program_t;
typedef struct _astnode_returnstmt astnode_returnstmt_t;
typedef struct _astnode_simplestmt astnode_simplestmt_t;
typedef struct _astnode_stmt astnode_stmt_t;
typedef struct _astnode_stmtlist astnode_stmtlist_t;
typedef struct _astnode_term astnode_term_t;
typedef struct _astnode_varstmt astnode_varstmt_t;
typedef struct _astnode_varstmtlist astnode_varstmtlist_t;
typedef struct _astnode_whilestmt astnode_whilestmt_t;

typedef enum
{
    ASTNODE_ASSIGNSTMT,
    ASTNODE_ATOM,
    ASTNODE_COMPAREEXPR,
    ASTNODE_COMPLEXSTMT,
    ASTNODE_COMPOUNDBODY,
    ASTNODE_COMPOUNDBODYLIST,
    ASTNODE_COMPOUNDSTMT,
    ASTNODE_ELSESTMT,
    ASTNODE_EXPR,
    ASTNODE_EXPRLIST,
    ASTNODE_FACTOR,
    ASTNODE_FUNCCALL,
    ASTNODE_FUNCDEF,
    ASTNODE_FUNCPARAMLIST,
    ASTNODE_IFSTMT,
    ASTNODE_MINOREXPR,
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
astnode_type_t;

astnode_node_t *astnode_init(astnode_type_t);
void astnode_free(astnode_node_t *);

struct _astnode_node
{
    astnode_type_t nodetype;
};

struct _astnode_assignstmt
{
    astnode_type_t nodetype;
    token_t identifier;
    token_t assignop;
    astnode_expr_t *expr;
};

struct _astnode_atom
{
    astnode_type_t nodetype;
    astnode_type_t atomtype;
    union
    {
        token_t token;
        astnode_funccall_t *funccall;
        astnode_expr_t *expr;
    }
    atom;
};

struct _astnode_compareexpr
{
    astnode_type_t nodetype;
    astnode_compareexpr_t *compareexpr;
    token_t compop;
    astnode_minorexpr_t *minorexpr;
};

struct _astnode_complexstmt
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_compoundstmt_t *compoundstmt;
        astnode_funcdef_t *funcdef;
    }
    stmt;
};

struct _astnode_compoundbody
{
    astnode_type_t nodetype;
    astnode_compoundbodylist_t *compoundbodylist;
};

struct _astnode_compoundbodylist
{
    astnode_type_t nodetype;
    astnode_compoundbodylist_t *compoundbodylist;
    astnode_stmt_t *stmt;
};

struct _astnode_compoundstmt
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_ifstmt_t *ifstmt;
        astnode_whilestmt_t *whilestmt;
    }
    stmt;
};

struct _astnode_elsestmt
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_compoundbody_t *compoundbody;
        astnode_ifstmt_t *ifstmt;
    }
    stmt;
};

struct _astnode_expr
{
    astnode_type_t nodetype;
    astnode_expr_t *expr;
    token_t expop;
    astnode_notexpr_t *notexpr;
};

struct _astnode_exprlist
{
    astnode_type_t nodetype;
    astnode_exprlist_t *exprlist;
    astnode_expr_t *expr;
};

struct _astnode_factor
{
    astnode_type_t nodetype;
    astnode_type_t factortype;
    union
    {
        astnode_atom_t *atom;
        astnode_factor_t *factor;
    }
    factor;
    token_t addop;
};

struct _astnode_funccall
{
    astnode_type_t nodetype;
    token_t identifier;
    astnode_exprlist_t *exprlist;
};

struct _astnode_funcdef
{
    astnode_type_t nodetype;
    token_t identifier;
    astnode_funcparamlist_t *funcparamlist;
    astnode_compoundbody_t *compoundbody;
};

struct _astnode_funcparamlist
{
    astnode_type_t nodetype;
    token_t identifier;
    astnode_funcparamlist_t *funcparamlist;
};

struct _astnode_ifstmt
{
    astnode_type_t nodetype;
    astnode_expr_t *expr;
    astnode_compoundbody_t *compoundbody;
    astnode_elsestmt_t *elsestmt;
};

struct _astnode_minorexpr
{
    astnode_type_t nodetype;
    astnode_minorexpr_t *minorexpr;
    token_t addop;
    astnode_term_t *term;
};

struct _astnode_notexpr
{
    astnode_type_t nodetype;
    token_t tnot;
    astnode_compareexpr_t *compareexpr;
};

struct _astnode_program
{
    astnode_type_t nodetype;
    astnode_stmtlist_t *stmtlist;
};

struct _astnode_returnstmt
{
    astnode_expr_t *expr;
};

struct _astnode_simplestmt
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_assignstmt_t *assignstmt;
        astnode_returnstmt_t *returnstmt;
        astnode_varstmt_t *varstmt;
        astnode_expr_t * expr;
    }
    stmt;
};

struct _astnode_stmt
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_complexstmt_t *complexstmt;
        astnode_simplestmt_t *simplestmt;
    }
    stmt;
};

struct _astnode_stmtlist
{
    astnode_type_t nodetype;
    astnode_stmtlist_t *stmtlist;
    astnode_stmt_t *stmt;
};

struct _astnode_term
{
    astnode_type_t nodetype;
    astnode_term_t *term;
    token_t mulop;
    astnode_factor_t *factor;
};

struct _astnode_varstmt
{
    astnode_type_t nodetype;
    astnode_varstmtlist_t *varstmtlist;
};

struct _astnode_varstmtlist
{
    astnode_type_t nodetype;
    astnode_varstmtlist_t *varstmtlist;
    astnode_type_t stmttype;
    union
    {
        token_t identifier;
        astnode_assignstmt_t *assignstmt;
    }
    stmt;
};

struct _astnode_whilestmt
{
    astnode_type_t nodetype;
    astnode_expr_t *expr;
    astnode_compoundbody_t *compoundbody;
};

#endif
