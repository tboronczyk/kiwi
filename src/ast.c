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

static void astnode_assignstmt_free(astnode_assignstmt_t *);
static void astnode_atom_free(astnode_atom_t *);
static void astnode_compareexpr_free(astnode_compareexpr_t *);
static void astnode_complexstmt_free(astnode_complexstmt_t *);
static void astnode_compoundbody_free(astnode_compoundbody_t *);
static void astnode_compoundbodylist_free(astnode_compoundbodylist_t *);
static void astnode_compoundstmt_free(astnode_compoundstmt_t *);
static void astnode_elsestmt_free(astnode_elsestmt_t *);
static void astnode_expr_free(astnode_expr_t *);
static void astnode_exprlist_free(astnode_exprlist_t *);
static void astnode_factor_free(astnode_factor_t *);
static void astnode_funccall_free(astnode_funccall_t *);
static void astnode_funcdef_free(astnode_funcdef_t *);
static void astnode_funcparamlist_free(astnode_funcparamlist_t *);
static void astnode_ifstmt_free(astnode_ifstmt_t *);
static void astnode_minorexpr_free(astnode_minorexpr_t *);
static void astnode_notexpr_free(astnode_notexpr_t *);
static void astnode_program_free(astnode_program_t *);
static void astnode_returnstmt_free(astnode_returnstmt_t *);
static void astnode_simplestmt_free(astnode_simplestmt_t *);
static void astnode_stmt_free(astnode_stmt_t *);
static void astnode_stmtlist_free(astnode_stmtlist_t *);
static void astnode_term_free(astnode_term_t *);
static void astnode_varstmt_free(astnode_varstmt_t *);
static void astnode_varstmtlist_free(astnode_varstmtlist_t *);
static void astnode_whilestmt_free(astnode_whilestmt_t *);

astnode_node_t *astnode_init(astnode_type_t type)
{
    astnode_node_t *n;
    size_t size;

    /* determine size to allocate for node */
    switch (type) {
        case ASTNODE_ASSIGNSTMT:
            size = sizeof(astnode_assignstmt_t);
            break;
        case ASTNODE_ATOM:
            size = sizeof(astnode_atom_t);
            break;
        case ASTNODE_COMPAREEXPR:
            size = sizeof(astnode_compareexpr_t);
            break;
        case ASTNODE_COMPLEXSTMT:
            size = sizeof(astnode_complexstmt_t);
            break;
        case ASTNODE_COMPOUNDBODYLIST:
            size = sizeof(astnode_compoundbodylist_t);
            break;
        case ASTNODE_COMPOUNDBODY:
            size = sizeof(astnode_compoundbody_t);
            break;
        case ASTNODE_COMPOUNDSTMT:
            size = sizeof(astnode_compoundstmt_t);
            break;
        case ASTNODE_ELSESTMT:
            size = sizeof(astnode_elsestmt_t);
            break;
        case ASTNODE_EXPRLIST:
            size = sizeof(astnode_exprlist_t);
            break;
        case ASTNODE_EXPR:
            size = sizeof(astnode_expr_t);
            break;
        case ASTNODE_FACTOR:
            size = sizeof(astnode_factor_t);
            break;
        case ASTNODE_FUNCCALL:
            size = sizeof(astnode_funccall_t);
            break;
        case ASTNODE_FUNCDEF:
            size = sizeof(astnode_funcdef_t);
            break;
        case ASTNODE_FUNCPARAMLIST:
            size = sizeof(astnode_funcparamlist_t);
            break;
        case ASTNODE_IFSTMT:
            size = sizeof(astnode_ifstmt_t);
            break;
        case ASTNODE_MINOREXPR:
            size = sizeof(astnode_minorexpr_t);
            break;
        case ASTNODE_NOTEXPR:
            size = sizeof(astnode_notexpr_t);
            break;
        case ASTNODE_PROGRAM:
            size = sizeof(astnode_program_t);
            break;
        case ASTNODE_RETURNSTMT:
            size = sizeof(astnode_returnstmt_t);
            break;
        case ASTNODE_SIMPLESTMT:
            size = sizeof(astnode_simplestmt_t);
            break;
        case ASTNODE_STMTLIST:
            size = sizeof(astnode_stmtlist_t);
            break;
        case ASTNODE_STMT:
            size = sizeof(astnode_stmt_t);
            break;
        case ASTNODE_TERM:
            size = sizeof(astnode_term_t);
            break;
        case ASTNODE_VARSTMTLIST:
            size = sizeof(astnode_varstmtlist_t);
            break;
        case ASTNODE_VARSTMT:
            size = sizeof(astnode_varstmt_t);
            break;
        case ASTNODE_WHILESTMT:
            size = sizeof(astnode_whilestmt_t);
            break;
        /* should never reach this */
        default:
            perror("Invalid astnode type passed on init");
            exit(EXIT_FAILURE);
    }

    if ((n = (astnode_node_t *)calloc(1, size)) == NULL) {
        perror("Allocate astnode failed");
        exit(EXIT_FAILURE);
    }
    n->nodetype = type;

    return n;
}

void astnode_free(astnode_node_t *n)
{
    /* free node with appropriate function */
    switch (n->nodetype) {
        case ASTNODE_ASSIGNSTMT:
            astnode_assignstmt_free((astnode_assignstmt_t *)n);
            break;
        case ASTNODE_ATOM:
            astnode_atom_free((astnode_atom_t *)n);
            break;
        case ASTNODE_COMPAREEXPR:
            astnode_compareexpr_free((astnode_compareexpr_t *)n);
            break;
        case ASTNODE_COMPLEXSTMT:
            astnode_complexstmt_free((astnode_complexstmt_t *)n);
            break;
        case ASTNODE_COMPOUNDBODY:
            astnode_compoundbody_free((astnode_compoundbody_t *)n);
            break;
        case ASTNODE_COMPOUNDBODYLIST:
            astnode_compoundbodylist_free((astnode_compoundbodylist_t *)n);
            break;
        case ASTNODE_COMPOUNDSTMT:
            astnode_compoundstmt_free((astnode_compoundstmt_t *)n);
            break;
        case ASTNODE_ELSESTMT:
            astnode_elsestmt_free((astnode_elsestmt_t *)n);
            break;
        case ASTNODE_EXPR:
            astnode_expr_free((astnode_expr_t *)n);
            break;
        case ASTNODE_EXPRLIST:
            astnode_exprlist_free((astnode_exprlist_t *)n);
            break;
        case ASTNODE_FACTOR:
            astnode_factor_free((astnode_factor_t *)n);
            break;
        case ASTNODE_FUNCCALL:
            astnode_funccall_free((astnode_funccall_t *)n);
            break;
        case ASTNODE_FUNCDEF:
            astnode_funcdef_free((astnode_funcdef_t *)n);
            break;
        case ASTNODE_FUNCPARAMLIST:
            astnode_funcparamlist_free((astnode_funcparamlist_t *)n);
            break;
        case ASTNODE_IFSTMT:
            astnode_ifstmt_free((astnode_ifstmt_t *)n);
            break;
        case ASTNODE_MINOREXPR:
            astnode_minorexpr_free((astnode_minorexpr_t *)n);
            break;
        case ASTNODE_NOTEXPR:
            astnode_notexpr_free((astnode_notexpr_t *)n);
            break;
        case ASTNODE_PROGRAM:
            astnode_program_free((astnode_program_t *)n);
            break;
        case ASTNODE_RETURNSTMT:
            astnode_returnstmt_free((astnode_returnstmt_t *)n);
            break;
        case ASTNODE_SIMPLESTMT:
            astnode_simplestmt_free((astnode_simplestmt_t *)n);
            break;
        case ASTNODE_STMT:
            astnode_stmt_free((astnode_stmt_t *)n);
            break;
        case ASTNODE_STMTLIST:
            astnode_stmtlist_free((astnode_stmtlist_t *)n);
            break;
        case ASTNODE_TERM:
            astnode_term_free((astnode_term_t *)n);
            break;
        case ASTNODE_VARSTMT:
            astnode_varstmt_free((astnode_varstmt_t *)n);
            break;
        case ASTNODE_VARSTMTLIST:
            astnode_varstmtlist_free((astnode_varstmtlist_t *)n);
            break;
        case ASTNODE_WHILESTMT:
            astnode_whilestmt_free((astnode_whilestmt_t *)n);
            break;
        /* should never reach this */
        default:
            perror("Invalid astnode type passed on free");
            exit(EXIT_FAILURE);
    }
}

static void astnode_assignstmt_free(astnode_assignstmt_t *n)
{
}

static void astnode_atom_free(astnode_atom_t *n)
{
}

static void astnode_compareexpr_free(astnode_compareexpr_t *n)
{
}

static void astnode_complexstmt_free(astnode_complexstmt_t *n)
{
}

static void astnode_compoundbody_free(astnode_compoundbody_t *n)
{
}

static void astnode_compoundbodylist_free(astnode_compoundbodylist_t *n)
{
}

static void astnode_compoundstmt_free(astnode_compoundstmt_t *n)
{
}

static void astnode_elsestmt_free(astnode_elsestmt_t *n)
{
}

static void astnode_expr_free(astnode_expr_t *n)
{
}

static void astnode_exprlist_free(astnode_exprlist_t *n)
{
}

static void astnode_factor_free(astnode_factor_t *n)
{
}

static void astnode_funccall_free(astnode_funccall_t *n)
{
}

static void astnode_funcdef_free(astnode_funcdef_t *n)
{
}

static void astnode_funcparamlist_free(astnode_funcparamlist_t *n)
{
}

static void astnode_ifstmt_free(astnode_ifstmt_t *n)
{
}

static void astnode_minorexpr_free(astnode_minorexpr_t *n)
{
}

static void astnode_notexpr_free(astnode_notexpr_t *n)
{
}

static void astnode_program_free(astnode_program_t *n)
{
}

static void astnode_returnstmt_free(astnode_returnstmt_t *n)
{
}

static void astnode_simplestmt_free(astnode_simplestmt_t *n)
{
}

static void astnode_stmt_free(astnode_stmt_t *n)
{
}

static void astnode_stmtlist_free(astnode_stmtlist_t *n)
{
}

static void astnode_term_free(astnode_term_t *n)
{
}

static void astnode_varstmt_free(astnode_varstmt_t *n)
{
}

static void astnode_varstmtlist_free(astnode_varstmtlist_t *n)
{
}

static void astnode_whilestmt_free(astnode_whilestmt_t *n)
{
}
