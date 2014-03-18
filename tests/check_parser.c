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
#include "astdump.h"
#include "scanner.h"
#include "parser.h"

Suite *parser_suite(void);

UFILE *ustdout;
UFILE *ustderr;

START_TEST (test_parser_parse)
{
    ustdout = u_finit(stdout, NULL, NULL);
    ustderr = u_finit(stderr, NULL, NULL);

    Scanner *s;
    scanner_init(&s, "./fib.kw");

    AST_Program *node;
    yyparse(s, &node);

    astdump_program(node);

    scanner_free(s);
}
END_TEST

Suite *parser_suite(void)
{
    Suite *s = suite_create("parser");
    TCase *tc_core = tcase_create("Core");
    tcase_add_test(tc_core, test_parser_parse);
    suite_add_tcase(s, tc_core);

    return s;
}

int main(void)
{
    Suite *s = parser_suite();
    SRunner *sr = srunner_create(s);
    srunner_run_all(sr, CK_NORMAL);
    int fails = srunner_ntests_failed(sr);
    srunner_free(sr);
    return (fails == 0) ? EXIT_SUCCESS : EXIT_FAILURE;
}
