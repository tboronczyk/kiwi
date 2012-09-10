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

typedef enum
{
    ASTNODE_PROGRAM,
    ASTNODE_STMTLIST,
    ASTNODE_STMT,
    ASTNODE_COMPLEXSTMT,
    ASTNODE_COMPOUNDSTMT,
    ASTNODE_IFSTMT,
    ASTNODE_EXPR,
    ASTNODE_COMPAREEXPR,
    ASTNODE_MINOREXPR,
    ASTNODE_TERM,
    ASTNODE_FACTOR,
    ASTNODE_ATOM,
    ASTNODE_FUNCCALL,
    ASTNODE_EXPRLIST,
    ASTNODE_ELSESTMT,
    ASTNODE_COMPOUNDBODY,
    ASTNODE_COMPOUNDBODYLIST,
    ASTNODE_WHILESTMT,
    ASTNODE_FUNCDEF,
    ASTNODE_FUNCPARAMLIST,
    ASTNODE_SIMPLESTMT,
    ASTNODE_ASSIGNSTMT,
    ASTNODE_ASSIGNOP,
    ASTNODE_RETURNSTMT,
    ASTNODE_VARSTMT,
    ASTNODE_VARSTMTLIST
}
astnode_type_t;

typedef struct _astnode_varstmtlist
{
    astnode_type_t nodetype;
    _astnode_varstmtlist *varstmtlist;
    astnode_type_t stmttype;
    union
    {
        token_t identifier;
        astnode_assignstmt_t *assignstmt;
    }
}
astnode_varstmtlist_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_varstmtlist_t *varstmtlist;
}
astnode_varstmt_t;

typedef struct
{
    astnode_expr_t *expr;
}
astnode_returnstmt_t;

typedef struct
{
    astnode_type_t nodetype;
    token_t identifier;
    token_t assignop;
    astnode_expr_t *expr;
}
astnode_assignstmt_t;

typedef struct _astnode_funcparamlist
{
    astnode_type_t nodetype;
    token_t identifier;
    struct _astnode_funcparamlist *funcparamlist;
}
astnode_funcparamlist_t;

typedef struct
{
    astnode_type_t nodetype;
    token_t identifier;
    astnode_funcparamlist_t *funcparamlist;
    astnode_compoundbody_t *compoundbody;
}
astnode_funcdef_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_expr_t *expr;
    astnode_compoundbody_t *compoundbody;
}
astnode_whilestmt_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_compoundbodylist_t *compoundbodylist;
    astnode_stmt_t *stmt;
}
astnode_compoundbodylist_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_compoundbodylist_t *compoundbodylist;
}
astnode_compoundbody_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_compoundbody_t *compoundbody;
        astnode_ifstmt_t *ifstmt;
    };
}
astnode_elsestmt_t;

typedef struct _astnode_exprlist
{
    astnode_type_t nodetype;
    struct _astnode_exprlist *exprlist;
    astnode_expr_t *expr;
}
astnode_exprlist_t;

typedef struct
{
    astnode_type_t nodetype;
    token_t identifier;
    astnode_exprlist_t *exprlist;
}
astnode_funccall_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_type_t atomtype;
    union
    {
        token_t token;
        astnode_funccall_t *funccall;
        astnode_expr_t *expr;
    }
}
astnode_atom_t;

typedef struct _astnode_factor
{
    astnode_type_t nodetype;
    astnode_type_t factortype;
    union
    {
        astnode_atom_t *atom;
        struct _astnode_factor *factor;
    }
    token_t addop;
}
astnode_factor_t;

typedef struct _astnode_term
{
    astnode_type_t nodetype;
    struct _astnode_term *term;
    token_t mulop;
    astnode_factor_t *factor;
}
astnode_term_t;

typedef struct _astnode_minorexpr
{
    astnode_type_t nodetype;
    struct _astnode_minorexpr *minorexpr;
    token_t addop;
    astnode_term_t *term;
}
astnode_minorexpr_t;

typedef struct _astnode_compareexpr
{
    astnode_type_t nodetype;
    struct _astnode_compareexpr *compareexpr;
    token_t compop;
    astnode_minorexpr_t *minorexpr;
}
astnode_compareexpr_t;

typedef struct
{
    astnode_type_t nodetype;
    token_t tnot;
    astnode_compareexpr_t *compareexpr;
}
astnode_notexpr_t;

typedef struct _astnode_expr
{
    astnode_type_t nodetype;
    struct _astnode_expr *expr;
    token_t expop;
    astnode_notexpr_t *notexpr;
}
astnode_expr_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_expr_t *expr;
    astnode_compoundbody_t *compoundbody;
    astnode_elsestmt_t *elsestmt;
}
astnode_ifstmt_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_ifstmt_t *ifstmt;
        astnode_whilestmt_t *whilestmt;
    };
}
astnode_compoundstmt_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_compoundstmt_t *compoundstmt;
        astnode_funcdef_t *funcdef;
    };
}
astnode_complexstmt_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_type_t stmttype;
    union
    {
        astnode_complexstmt_t *complexstmt;
        astnode_simplestmt_t *simplestmt;
    };
}
astnode_stmt_t;

typedef struct _astnode_stmtlist
{
    astnode_type_t nodetype;
    struct _astnode_stmtlist *stmtlist;
    astnode_stmt_t *stmt;
}
astnode_stmtlist_t;

typedef struct
{
    astnode_type_t nodetype;
    astnode_stmtlist_t *stmtlist;
}
astnode_program_t;

#endif AST_H