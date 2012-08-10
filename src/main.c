#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "scanner.h"
#include "y.tab.h"
#include "unicode/ustdio.h"

extern void scanner_parse(Scanner *s); // y.tab.c

int scanner_error(Scanner *s, const char *str) {
  fprintf(stderr, "%s line %d\n", str, s->lineno);
  return 1;
}

#define perror_exit(f) \
    perror(f); \
    exit(EXIT_FAILURE)

UFILE *ustdout;

int scanner_lex(YYSTYPE *lvalp, Scanner *s)
{
    scanner_token(s);
    // force re-read on comments
    if (s->name == T_COMMENT) {
        return scanner_lex(lvalp, s);
    }
    else {
/*
        if (s->name == T_NUMBER) {
            yylval.number = ...
        }
        else if (s->name == T_IDENTIFIER || s->name == T_STRING) {
            yylval.string = ...
        }
*/
        return s->name;
    }
}

int main()
{
    Scanner *s;

    ustdout = u_finit(stdout, NULL, NULL);

    s = scanner_init();
    scanner_parse(s);
    scanner_free(s);

    return 0;
}

