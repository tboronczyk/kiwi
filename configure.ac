m4_define([KIWI_MAJOR], 0)
m4_define([KIWI_MINOR], 0)
m4_define([KIWI_PATCH], 0)
m4_define([KIWI_BUGS], [])

AC_PREREQ([2.64])
AC_INIT([kiwi], [KIWI_MAJOR.KIWI_MINOR.KIWI_PATCH], [KIWI_BUGS])
AC_CONFIG_SRCDIR([src/main.c])

AM_INIT_AUTOMAKE([foreign -Wall -Werror])
AM_CONFIG_HEADER([config.h])

AC_ARG_ENABLE(debug,
AS_HELP_STRING([--enable-debug],
               [enable debugging, default: no]),
[case "${enableval}" in
             yes) debug=true ;;
             no)  debug=false ;;
             *)   AC_MSG_ERROR([bad value ${enableval} for --enable-debug]) ;;
esac],
[debug=false])

AM_CONDITIONAL(DEBUG, test x"$debug" = x"true")
if test -n "${CFLAGS}"; then
    user_set_cflags=yes
fi
AC_LANG([C])
AC_PROG_CC
AC_PROG_YACC
if test X$user_set_cflags != Xyes; then
    CFLAGS=""
fi

AC_CHECK_HEADERS([assert.h stdarg.h stdio.h stdlib.h string.h])
AC_CHECK_HEADERS([unicode/uchar.h unicode/ustdio.h unicode/ustring.h])

AC_CONFIG_FILES([Makefile src/Makefile])
AC_OUTPUT
