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
static bool ctlog_enable = false;

/*============================================================================*/
// Sequence number to keep track of which log event this was, so we can know if
// one was dropped or is missing.
uint16_t ctlog_sequenceNumber = 0;



/*==============================================================================
 *                                  Functions
 *============================================================================*/
/*============================================================================*/
void
ctlog_fputc_json( char c, FILE* stream )
{
    if ( iscntrl(c) || (c == '"') || (c == '\\') )
    {
        fputc( '\\', stream );

        switch (c)
        {
            case '"':
                fputc( '"', stream );
                break;
            case '\\':
                fputc( '\\', stream );
                break;
            case '/':
                fputc( '/', stream );
                break;
            case '\b':
                fputc( 'b', stream );
                break;
            case '\f':
                fputc( 'f', stream );
                break;
            case '\n':
                fputc( 'n', stream );
                break;
            case '\r':
                fputc( 'r', stream );
                break;
            case '\t':
                fputc( 't', stream );
                break;
            default:
                fprintf( stream, "u%04X", c );
        }
    }
    else
    {
        fputc( c, stream );
    }
}



/*==============================================================================
 *                               Public Functions
 *============================================================================*/
/*============================================================================*/
void
ctlog_setEnabled( bool enable )
{
    ctlog_enable = enable;
}

/*============================================================================*/
bool
ctlog_isEnabled( void )
{
    return ctlog_enable;
}


/*============================================================================*/
void
ctlog_fprintf( FILE* stream, char level, cmodule_index_t moduleIndex, uint32_t line, int nArgs, ... )
{
    if ( ctlog_enable )
    {
        fprintf( stream, "$TL" "%"PRIu16 "," "%"PRIu16 ",%c," "%"PRIu32 "," "%"PRIu32 ",%d,", CTLOG_VERSION, ctlog_sequenceNumber, level, moduleIndex, line, nArgs );

        if ( nArgs > 0 )
        {
            int i;
            va_list vl;

            va_start( vl, nArgs );
            for ( i = 0; i < (2*nArgs); i += 2 )
            {
                uint8_t type = (uint8_t)va_arg( vl, int );
                fprintf( stream, "%"PRIu8",", type );

                switch ( type )
                {
                    case _CTLOG_TYPE_UINT: fprintf( stream, "%"PRIu32, (uint32_t)va_arg( vl, int ) ); break;
                    case _CTLOG_TYPE_INT:  fprintf( stream, "%"PRId32, (int32_t)va_arg( vl, int ) ); break;

                    case _CTLOG_TYPE_STRING:
                    {
                        fputc( '^', stream );
                        fputc( '\x00', stream );
                        fprintf( stream, "%s", va_arg( vl, char* ) );
                        fputc( '$', stream );
                        fputc( '\x00', stream );
                    }
                    break;

                    case _CTLOG_TYPE_BOOL: fprintf( stream, "%"PRIu8, (uint8_t)va_arg( vl, int ) ); break;
                    case _CTLOG_TYPE_CHAR: fprintf( stream, "%"PRIu8, (uint8_t)va_arg( vl, int ) ); break;
                    default: assert( false ); break;
                }

                fputc( ',', stream );
            }
            va_end( vl );
        }

        fputs( "\n", stream );
    }

    ctlog_sequenceNumber += 1;
}


/*============================================================================*/
void
ctlog_json_fprintf( FILE* stream, char level, cmodule_index_t moduleIndex, uint32_t line, int nArgs, ... )
{
    if ( ctlog_enable )
    {
        fprintf( stream, "{\"ctlog\":" "%"PRIu16 ",\"seq\":" "%"PRIu16 ",\"lvl\":\"%c\",\"mi\":" "%"PRIu32 ",\"ml\":" "%"PRIu32 ",\"args\":[", CTLOG_VERSION, ctlog_sequenceNumber, level, moduleIndex, line );

        if ( nArgs > 0 )
        {
            int i;
            int n = 2*nArgs;
            va_list vl;

            va_start( vl, nArgs );
            for ( i = 0; i < n; i += 2 )
            {
                uint8_t type = (uint8_t)va_arg( vl, int );
                fprintf( stream, "{\"t\":" "%"PRIu8 ",\"v\":", type );

                switch ( type )
                {
                    case _CTLOG_TYPE_UINT: fprintf( stream, "%"PRIu32, (uint32_t)va_arg( vl, int ) ); break;
                    case _CTLOG_TYPE_INT:  fprintf( stream, "%"PRId32, (int32_t)va_arg( vl, int ) ); break;

                    case _CTLOG_TYPE_STRING:
                    {
                        char* s = va_arg( vl, char* );

                        fputc( '"', stream );
                        while ( *s != '\0' )
                        {
                            ctlog_fputc_json( *s, stream );
                            s++;
                        }
                        fputc( '"', stream );
                    }
                    break;

                    case _CTLOG_TYPE_BOOL:
                    {
                        if ( va_arg( vl, int ) )
                        {
                            fputs( "true", stream );
                        }
                        else
                        {
                            fputs( "false", stream );
                        }
                    }
                    break;

                    case _CTLOG_TYPE_CHAR:
                    {
                        fputc( '"', stream );
                        ctlog_fputc_json( (char)va_arg( vl, int ), stream );
                        fputc( '"', stream );
                    }
                    break;

                    default: assert( false ); break;
                }

                fputc( '}', stream );
                if ( i + 2 < n )
                {
                    fputc( ',', stream );
                }
            }
            va_end( vl );
        }

        fputs( "]}\n", stream );
    }

    ctlog_sequenceNumber += 1;
}
