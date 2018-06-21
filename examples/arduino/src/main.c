/**
 * Ctlog exampled for the Arduino Uno.
 * Copyright (C) 2016 Josh Lubawy <jlubawy@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

// Comment out if you want to disable tokenized logging
#define ENABLE_CTLOG

#include <avr/io.h>

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "ctlog.h"

// Define the C module (the file basename should usually be used)
CMODULE_DEFINE( main );

static FILE g_uart_stream = {0};

static char g_line_buffer[256];

static void
uart_init( uint32_t baud )
{
    /* See table 20-1 for baud rate calculations */
    uint16_t ubrr = (F_CPU / (16*baud)) - 1;

    UBRR0H = (uint8_t)(ubrr >> 8);
    UBRR0L = (uint8_t)ubrr;

    /* Disable 2x TX speed */
    UCSR0A &= ~_BV(U2X0);

    /* Enable RX and TX */
    UCSR0B = _BV(RXEN0) | _BV(TXEN0);

    /* Set 8-N-1 frame */
    UCSR0C = (3 << UCSZ00);
}


static int
uart_getc( FILE* stream )
{
    while ( (UCSR0A & _BV(RXC0)) == 0 ); /* wait for a character to be received */
    return UDR0;
}


static int
uart_putc( char c, FILE* stream )
{
    while ( (UCSR0A & _BV(UDRE0)) == 0 ); /* wait for the TX buffer to be ready */
    UDR0 = c; /* put char into TX buffer */
    return 0;
}


int
main( void )
{
    int i;
    int c;
    int empty;

    ctlog_setStream( stdout );

    DDRB |= _BV(5);

    uart_init( 9600 );
    fdev_setup_stream( &g_uart_stream, uart_putc, uart_getc, _FDEV_SETUP_RW );

    stdin  = &g_uart_stream;
    stdout = &g_uart_stream;
    stderr = &g_uart_stream;

#ifdef ENABLE_CTLOG
    ctlog_setStream( stdout );
#endif

    // Log a ~1KB string to demonstrate the memory savings that tokenized
    // logging can provide.
#ifdef ENABLE_CTLOG
    CTLOG_INFO( "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris lacus ligula, ultrices sed condimentum ac, aliquet in nulla. Ut lobortis pulvinar dui, auctor consectetur nulla. Suspendisse id malesuada neque. Cras pretium nisl quis felis hendrerit tristique. Phasellus sed porttitor dui. Phasellus aliquam fermentum elit at aliquet. Nullam porta, tortor vitae sagittis dapibus, felis libero dictum nunc, eu tincidunt orci diam at risus. Donec quis bibendum turpis. Maecenas ultrices imperdiet nulla non laoreet. Sed euismod rhoncus lorem, porttitor varius nunc tempus at. Aenean dignissim fringilla dui ac commodo. Mauris iaculis et ipsum id malesuada. Donec semper magna a malesuada dictum. Aliquam vehicula ligula vitae venenatis elementum. Phasellus congue eleifend viverra. Suspendisse potenti. Fusce aliquet, massa ac tristique egestas, dui tellus molestie mi, quis accumsan lacus eros quis tellus. Nulla ipsum nulla, dapibus in purus sed, pellentesque volutpat tortor. Aliquam tincidunt interdum arcu ac maximus." );
#else
    fputs( "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris lacus ligula, ultrices sed condimentum ac, aliquet in nulla. Ut lobortis pulvinar dui, auctor consectetur nulla. Suspendisse id malesuada neque. Cras pretium nisl quis felis hendrerit tristique. Phasellus sed porttitor dui. Phasellus aliquam fermentum elit at aliquet. Nullam porta, tortor vitae sagittis dapibus, felis libero dictum nunc, eu tincidunt orci diam at risus. Donec quis bibendum turpis. Maecenas ultrices imperdiet nulla non laoreet. Sed euismod rhoncus lorem, porttitor varius nunc tempus at. Aenean dignissim fringilla dui ac commodo. Mauris iaculis et ipsum id malesuada. Donec semper magna a malesuada dictum. Aliquam vehicula ligula vitae venenatis elementum. Phasellus congue eleifend viverra. Suspendisse potenti. Fusce aliquet, massa ac tristique egestas, dui tellus molestie mi, quis accumsan lacus eros quis tellus. Nulla ipsum nulla, dapibus in purus sed, pellentesque volutpat tortor. Aliquam tincidunt interdum arcu ac maximus.\n", stdout );
#endif

    for (;;)
    {
        // Ready to receive input
        fputs( "> ", stdout );

        // Read input
        i = 0;
        do
        {
            c = fgetc( stdin );
            if ( !ferror( stdin ) )
            {
                // If no error then process the character
                if ( (char)c == '\n' )
                {
                    // If newline then stop reading and process the line buffer
                    g_line_buffer[ i ] = '\0';
                    goto EXIT_INNER;
                }
                else if ( (char)c == '\r' )
                {
                    // Else if carriage-return, drop the character
                    continue;
                }
                else
                {
                    // Else echo any other character and add it to the line buffer
                    fputc( (char)c, stdout );
                    g_line_buffer[ i ] = (char)c;
                }

                i += 1;
            }
            else
            {
                // Else if there is an error, log it
                fputs( "\nError!\n> ", stdout );
                i = 0;
            }
        }
        while ( i < (sizeof(g_line_buffer)-1) );

EXIT_INNER:
        g_line_buffer[ sizeof(g_line_buffer)-1 ] = '\0';

        fputc( '\n', stdout );

#ifdef ENABLE_CTLOG
        CTLOG_VAR_INFO( "line=%s", 1, CTLOG_TYPE_STRING( g_line_buffer ) );
#else
        fputs( g_line_buffer, stdout );
        fputc( '\n', stdout );
#endif

        fflush( stdout );
    }

EXIT:
    return 0;
}
