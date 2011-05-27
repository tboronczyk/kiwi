/*
 * Copyright (c) 2011, Timothy Boronczyk
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

#include <ctype.h>
#include <stdio.h>
#include <string.h>
#include "uhelpers.h"

unsigned short u_getc(FILE *stream, char *bytes) {
    unsigned short i;
    // mask values for bit pattern of first byte in multi-byte UTF-8 sequences:
    // 192 - 110xxxxx - for U+0080 to U+07FF
    // 224 - 1110xxxx - for U+0800 to U+FFFF
    // 240 - 11110xxx - for U+010000 to U+1FFFFF
    // @TODO: This may move to uhelpers.h if needed elsewhere
    static unsigned short mask[] = {192, 224, 240};

    // initialize buffer
    memset(bytes, 0, U_MAX_BYTES);

    // read bytes into buffer
    bytes[0] = getc(stream);
    // end of file was reached or an error was encountered
    if (bytes[0] == EOF) {
        return 0;
    }
    else {
        i = 0;
        // read subsequent character bytes
        while ((bytes[0] & mask[i]) == mask[i] && i < U_MAX_BYTES - 1) {
            i++;
            // end of file or error here means bad UTF-8 data!
            if ((bytes[i] = getc(stream)) == EOF) {
                return 0;
            }
        }
        return i + 1;
    }
}

unsigned short u_isalpha(char *bytes) {
    // @FIXME: u_isalpha is used to identify tokens and should therefore
    // understand UTF-8. This is temporary; do not use ctype!
    return isalpha(bytes[0]);
    
}

unsigned short u_isalnum(char *bytes) {
    // @FIXME: u_isalnum is used to identify tokens and should therefore
    // understand UTF-8. This is temporary; do not use ctype!
    return isalnum(bytes[0]);
}

unsigned short u_isdigit(char *bytes) {
    // digits are 0 .. 9, which are encoded as single-byte sequences
    return memchr("0123456789", bytes[0], 10) != NULL;
}

unsigned short u_is2digit(char *bytes) {
    // binary digits are either 0 or 1, which are encoded as single-byte
    // sequences
    return bytes[0] == '0' || bytes[0] == '1';
}

unsigned short u_is8digit(char *bytes) {
    // octal digits are 0 .. 7, which are encoded as single-byte sequences
    return memchr("01234567", bytes[0], 8) != NULL;
}

unsigned short u_is16digit(char *bytes) {
    // hexadecimal digits are 0 .. 9 and case-insensitive A .. F,  which are
    // encoded as single-byte sequences
    return u_isdigit(bytes) ||
        memchr("ABCDEF", bytes[0], 6) != NULL || 
        memchr("abcdef", bytes[0], 6) != NULL;
}

unsigned short u_isnewline(char *bytes) {
    return 10 == (int)bytes[0] ||   // line feed
        13 == (int)bytes[0];        // carriage return
    // @FIXME: add checks for U+2028 LS and U+2029 PS
        
}

unsigned short u_isspace(char *bytes) {
    return u_isnewline(bytes) ||
        9 == (int)bytes[0] ||       // tab
        13 == (int)bytes[0] ||      // vertical tab
        14 == (int)bytes[0] ||      // form feed
        32 == (int)bytes[0];        // space
    // @FIXME: add checks for U+00A0 NBSP, U+FEFF BOM, and category Zs USP
}

