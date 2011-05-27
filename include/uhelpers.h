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

#ifndef UHELPERS_H
#define UHELPERS_H

#include <stdio.h>

/**
 * UTF-8 encoding restricts character sequences to four bytes
 */
#define U_MAX_BYTES 4

/**
 * Get UTF-8 byte sequence from stream
 *
 * Reads the UTF-8 character byte sequence currently pointed to by the
 * internal file position indicator of the specified stream. The position
 * indicator is advanced to point to the next character. The read sequence
 * is placed into the specified byte buffer.
 *
 * @param stream An open file descriptor
 * @param bytes  Destination buffer for byte sequence
 *
 * @return the size in bytes of the character, 0 if EOF is reached or a
 * reading error occurs.
 */
unsigned short u_getc(FILE *, char *);

/**
 * Check if UTF-8 byte sequence is alphanumeric
 *
 * Checks if the UTF-8 character byte sequence is a decimal digit or a letter.
 *
 * @param bytes UTF-8 character byte sequence
 *
 * @return non-zero if the provided sequence is considered a digit or letter,
 * otherwise 0.
 */
unsigned short u_isalnum(char *);

/**
 * Check if UTF-8 byte sequence is a letter
 *
 * Checks if the UTF-8 character byte sequence is an alphabetic letter.
 *
 * @param bytes UTF-8 character byte sequence
 *
 * return non-zero if the provided sequence is considered a digit or letter,
 * otherwise 0.
 */
unsigned short u_isalpha(char *);

/**
 * Check if UTF-8 byte sequence is a decimal digit
 *
 * Checks if the UTF-8 character byte sequence is considered a decimal digit.
 *
 * @param bytes UTF-8 character byte sequence
 *
 * @return non-zero if the provided sequence is considered a digit, otherwise
 * 0.
 */
unsigned short u_isdigit(char *);

/**
 * Check if UTF-8 byte sequence is a binary digit
 *
 * Checks if the UTF-8 character byte sequence is considered a binary digit.
 *
 * @param bytes UTF-8 character byte sequence
 *
 * @return non-zero if the provided sequence is considered a binary digit,
 * otherwise 0.
 */
unsigned short u_is2digit(char *);

/**
 * Check if UTF-8 byte sequence is an octal digit
 *
 * Checks if the UTF-8 character byte sequence is considered an octal digit.
 *
 * @param bytes UTF-8 character byte sequence
 *
 * @return non-zero if the provided sequence is considered an octal digit,
 * otherwise 0.
 */
unsigned short u_is8digit(char *);

/**
 * Check if UTF-8 byte sequence is a hexadecimal digit
 *
 * Checks if the UTF-8 character byte sequence is considered a hexadecimal 
 * digit.
 *
 * @param bytes UTF-8 character byte sequence
 *
 * @return non-zero if the provided sequence is considered a hexadecimal digit,
 * otherwise 0.
 */
unsigned short u_is16digit(char *);

/**
 * Check if UTF-8 byte sequence is a newline
 *
 * Checks if the UTF-8 character byte sequence is considered a newline.
 *
 * @param bytes UTF-8 character byte sequence
 *
 * @return non-zero if the provided sequence is considered a newline, otherwise
 * 0.
 */
unsigned short u_isnewline(char *);

/**
 * Check if UTF-8 byte sequence is white-space
 *
 * Checks if the UTF-8 character byte sequence is considered white-space.
 *
 * @param bytes UTF-8 character byte sequence
 *
 * @return non-zero if the provided sequence is considered white-space,
 * otherwise 0.
 */
unsigned short u_isspace(char *);

#endif

