#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "scanner.h"
#include "token.h"
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

int yylex() {
    Token *t;
    int i;

    t = (Token *)scanner_token(s);
    // force re-read on comments
    if (t->name == T_COMMENT) {
        token_free(t);
        return yylex();
    }
    else {
/*
        if (t->name == T_NUMBER) {
            yylval.number = ...
        }
        else if (t->name == T_IDENTIFIER || t->name == T_STRING) {
            yylval.string = ...
        }
*/
        i = t->name;
        token_free(t);
        return i;
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

