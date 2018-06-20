/**
 * Tokenized logging framework.
 */

#include <assert.h>
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
        fputs( "$TL" CTLOG_VERSION ",", stream );
        fprintf( stream, "%"PRIu16",%c,%"PRIu32",%"PRIu32",%d,", ctlog_sequenceNumber, level, moduleIndex, line, nArgs );

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
