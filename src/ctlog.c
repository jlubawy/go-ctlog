/**
 * Tokenized logging framework.
 */

#include <assert.h>
#include <ctype.h>
#include <inttypes.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

#include "ctlog.h"
#include "cmodule.h"

/*==============================================================================
 *                                   Macros
 *============================================================================*/
/*============================================================================*/


/*==============================================================================
 *                                   Globals
 *============================================================================*/
/*============================================================================*/
// Global which controls whether logging is enabled at runtime. Defaults to off.
static FILE* g_stream = NULL;

/*============================================================================*/
// Sequence number to keep track of which log event this was, so we can know if
// one was dropped or is missing.
static uint16_t g_sequence_number = 0;



/*==============================================================================
 *                                  Functions
 *============================================================================*/
/*============================================================================*/
static void
ctlog_fputc_json( char c, FILE* stream )
{
    if ( iscntrl(c) || (c == '"') || (c == '\\') )
    {
        fputc( '\\', g_stream );

        switch (c)
        {
            case '"':
                fputc( '"', g_stream );
                break;
            case '\\':
                fputc( '\\', g_stream );
                break;
            case '/':
                fputc( '/', g_stream );
                break;
            case '\b':
                fputc( 'b', g_stream );
                break;
            case '\f':
                fputc( 'f', g_stream );
                break;
            case '\n':
                fputc( 'n', g_stream );
                break;
            case '\r':
                fputc( 'r', g_stream );
                break;
            case '\t':
                fputc( 't', g_stream );
                break;
            default:
                fprintf( g_stream, "u%04X", c );
        }
    }
    else
    {
        fputc( c, g_stream );
    }
}



/*==============================================================================
 *                               Public Functions
 *============================================================================*/
/*============================================================================*/
void
ctlog_setStream( FILE* stream )
{
    g_stream = stream;
}


/*============================================================================*/
void
ctlog_fprintf( char level, cmodule_index_t moduleIndex, uint32_t line, int nArgs, ... )
{
    if ( g_stream != NULL )
    {
        fprintf( g_stream, "$TL" "%"PRIu16 "," "%"PRIu16 ",%c," "%"PRIu32 "," "%"PRIu32 ",%d,", CTLOG_VERSION, g_sequence_number, level, moduleIndex, line, nArgs );

        if ( nArgs > 0 )
        {
            int i;
            va_list vl;

            va_start( vl, nArgs );
            for ( i = 0; i < (2*nArgs); i += 2 )
            {
                uint8_t type = (uint8_t)va_arg( vl, int );
                fprintf( g_stream, "%"PRIu8",", type );

                switch ( type )
                {
                    case CTLOG_TYPE_N_UINT: fprintf( g_stream, "%"PRIu32, (uint32_t)va_arg( vl, int ) ); break;
                    case CTLOG_TYPE_N_INT:  fprintf( g_stream, "%"PRId32, (int32_t)va_arg( vl, int ) ); break;

                    case CTLOG_TYPE_N_STRING:
                    {
                        fputc( '^', g_stream );
                        fputc( '\x00', g_stream );
                        fprintf( g_stream, "%s", va_arg( vl, char* ) );
                        fputc( '$', g_stream );
                        fputc( '\x00', g_stream );
                    }
                    break;

                    case CTLOG_TYPE_N_BOOL: fprintf( g_stream, "%"PRIu8, (uint8_t)va_arg( vl, int ) ); break;
                    case CTLOG_TYPE_N_CHAR: fprintf( g_stream, "%"PRIu8, (uint8_t)va_arg( vl, int ) ); break;
                    default: assert( false ); break;
                }

                fputc( ',', g_stream );
            }
            va_end( vl );
        }

        fputs( "\n", g_stream );
    }

    g_sequence_number += 1;
}


/*============================================================================*/
void
ctlog_json_fprintf( char level, cmodule_index_t moduleIndex, uint32_t line, int nArgs, ... )
{
    if ( g_stream != NULL )
    {
        fprintf( g_stream, "{\"ctlog\":" "%"PRIu16 ",\"seq\":" "%"PRIu16 ",\"lvl\":\"%c\",\"mi\":" "%"PRIu32 ",\"ml\":" "%"PRIu32 ",\"args\":[", CTLOG_VERSION, g_sequence_number, level, moduleIndex, line );

        if ( nArgs > 0 )
        {
            int i;
            int n = 2*nArgs;
            va_list vl;

            va_start( vl, nArgs );
            for ( i = 0; i < n; i += 2 )
            {
                uint8_t type = (uint8_t)va_arg( vl, int );
                fprintf( g_stream, "{\"t\":" "%"PRIu8 ",\"v\":", type );

                switch ( type )
                {
                    case CTLOG_TYPE_N_UINT: fprintf( g_stream, "%"PRIu32, (uint32_t)va_arg( vl, int ) ); break;
                    case CTLOG_TYPE_N_INT:  fprintf( g_stream, "%"PRId32, (int32_t)va_arg( vl, int ) ); break;

                    case CTLOG_TYPE_N_STRING:
                    {
                        char* s = va_arg( vl, char* );

                        fputc( '"', g_stream );
                        while ( *s != '\0' )
                        {
                            ctlog_fputc_json( *s, g_stream );
                            s++;
                        }
                        fputc( '"', g_stream );
                    }
                    break;

                    case CTLOG_TYPE_N_BOOL:
                    {
                        if ( va_arg( vl, int ) )
                        {
                            fputs( "true", g_stream );
                        }
                        else
                        {
                            fputs( "false", g_stream );
                        }
                    }
                    break;

                    case CTLOG_TYPE_N_CHAR:
                    {
                        fputc( '"', g_stream );
                        ctlog_fputc_json( (char)va_arg( vl, int ), g_stream );
                        fputc( '"', g_stream );
                    }
                    break;

                    default: assert( false ); break;
                }

                fputc( '}', g_stream );
                if ( i + 2 < n )
                {
                    fputc( ',', g_stream );
                }
            }
            va_end( vl );
        }

        fputs( "]}\n", g_stream );
    }

    g_sequence_number += 1;
}
