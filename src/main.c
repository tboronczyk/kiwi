#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "scanner.h"
#include "y.tab.h"
#include "unicode/ustdio.h"

extern int yyparse(void);

int yyerror(char *s) {
  fprintf(stderr, "%s\n", s);
  return 1;
}

#define perror_exit(f) \
    perror(f); \
    exit(EXIT_FAILURE)

UFILE *ustdout;
Scanner *s;

int yylex()
{
    scanner_token(s);
    // force re-read on comments
    if (s->name == T_COMMENT) {
        return yylex();
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

int main(int argc, char **argv)
{
    ustdout = u_finit(stdout, NULL, NULL);
    s = scanner_init("stdin");
    yyparse();
    scanner_free(s);
    return 0;
}

