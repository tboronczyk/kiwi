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

#include <check.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unicode/ustdio.h>
#include "ast.h"
#include "scanner.h"
#include "parser.h"

#define SCAN_AND_ASSERT(s,t) do { \
    scanner_token((s)); \
    u_fprintf(ustdout, "%d %S\n", s->name, s->tbuf); \
    ck_assert_int_eq((s)->name, (t)); \
} while (0)

Suite *scanner_suite(void);

START_TEST (test_scanner_tokenize)
{
    UFILE *ustdout = u_finit(stdout, NULL, NULL);

    Scanner *s;
    scanner_init(&s, "./extra/scanner-in01.dat");
    SCAN_AND_ASSERT(s, T_ADD);
    SCAN_AND_ASSERT(s, T_SUBTRACT);
    SCAN_AND_ASSERT(s, T_MULTIPLY);
    SCAN_AND_ASSERT(s, T_DIVIDE);
    SCAN_AND_ASSERT(s, T_MODULO);
    SCAN_AND_ASSERT(s, T_ASSIGN);
    SCAN_AND_ASSERT(s, T_ADD_ASSIGN);
    SCAN_AND_ASSERT(s, T_SUBTRACT_ASSIGN);
    SCAN_AND_ASSERT(s, T_MULTIPLY_ASSIGN);
    SCAN_AND_ASSERT(s, T_DIVIDE_ASSIGN);
    SCAN_AND_ASSERT(s, T_MODULO_ASSIGN);
    SCAN_AND_ASSERT(s, T_EQUAL);
    SCAN_AND_ASSERT(s, T_NOT_EQUAL);
    SCAN_AND_ASSERT(s, T_IS);
    SCAN_AND_ASSERT(s, T_LESS);
    SCAN_AND_ASSERT(s, T_LESS_EQUAL);
    SCAN_AND_ASSERT(s, T_GREATER);
    SCAN_AND_ASSERT(s, T_GREATER_EQUAL);
    SCAN_AND_ASSERT(s, T_TRUE);
    SCAN_AND_ASSERT(s, T_FALSE);
    SCAN_AND_ASSERT(s, T_AND);
    SCAN_AND_ASSERT(s, T_OR);
    SCAN_AND_ASSERT(s, T_NOT);
    SCAN_AND_ASSERT(s, T_IF);
    SCAN_AND_ASSERT(s, T_ELSE);
    SCAN_AND_ASSERT(s, T_FUNC);
    SCAN_AND_ASSERT(s, T_RETURN);
    SCAN_AND_ASSERT(s, T_WHILE);
    SCAN_AND_ASSERT(s, T_VAR);
    SCAN_AND_ASSERT(s, T_DOT);
    SCAN_AND_ASSERT(s, T_CONCAT);
    SCAN_AND_ASSERT(s, T_COMMA);
    SCAN_AND_ASSERT(s, T_COLON);
    SCAN_AND_ASSERT(s, T_LBRACE);
    SCAN_AND_ASSERT(s, T_RBRACE);
    SCAN_AND_ASSERT(s, T_LPAREN);
    SCAN_AND_ASSERT(s, T_RPAREN);
    SCAN_AND_ASSERT(s, T_COMMENT);
    SCAN_AND_ASSERT(s, T_COMMENT);
    SCAN_AND_ASSERT(s, T_COMMENT);
    SCAN_AND_ASSERT(s, T_STRING);
    SCAN_AND_ASSERT(s, T_STRING);
    SCAN_AND_ASSERT(s, T_STRING);
    SCAN_AND_ASSERT(s, T_STRING);
    SCAN_AND_ASSERT(s, T_STRING);
    SCAN_AND_ASSERT(s, T_STRING);
    SCAN_AND_ASSERT(s, T_NUMBER);
    SCAN_AND_ASSERT(s, T_NUMBER);
    SCAN_AND_ASSERT(s, T_NUMBER);
    SCAN_AND_ASSERT(s, T_NUMBER);
    SCAN_AND_ASSERT(s, T_NUMBER);
    SCAN_AND_ASSERT(s, T_IDENTIFIER);
    SCAN_AND_ASSERT(s, T_IDENTIFIER);
    SCAN_AND_ASSERT(s, T_EOF);
    scanner_free(s);
}
END_TEST

Suite *scanner_suite(void)
{
    Suite *s = suite_create("scanner");
    TCase *tc_core = tcase_create("Core");
    tcase_add_test(tc_core, test_scanner_tokenize);
    suite_add_tcase(s, tc_core);

    return s;
}

int main(void)
{
    Suite *s = scanner_suite();
    SRunner *sr = srunner_create(s);
    srunner_run_all(sr, CK_NORMAL);
    int fails = srunner_ntests_failed(sr);
    srunner_free(sr);
    return (fails == 0) ? EXIT_SUCCESS : EXIT_FAILURE;
}

