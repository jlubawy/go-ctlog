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
ctlog_printf( char level, cmodule_index_t moduleIndex, uint32_t line, int nArgs, ... )
{
    if ( ctlog_enable )
    {
        fputs( "$TL" CTLOG_VERSION ",", stdout );
        fprintf( stdout, "%"PRIu16",%c,%"PRIu32",%"PRIu32",%d,", ctlog_sequenceNumber, level, moduleIndex, line, nArgs );

        if ( nArgs > 0 )
        {
            int i;
            va_list vl;

            va_start( vl, nArgs );
            for ( i = 0; i < (2*nArgs); i += 2 )
            {
                uint8_t type = (uint8_t)va_arg( vl, int );
                fprintf( stdout, "%"PRIu8",", type );

                switch ( type )
                {
                    case _CTLOG_TYPE_UINT: fprintf( stdout, "%"PRIu32, (uint32_t)va_arg( vl, int ) ); break;
                    case _CTLOG_TYPE_INT:  fprintf( stdout, "%"PRId32, (int32_t)va_arg( vl, int ) ); break;

                    case _CTLOG_TYPE_STRING:
                    {
                        fputc( '^', stdout );
                        fputc( '\x00', stdout );
                        fprintf( stdout, "%s", va_arg( vl, char* ) );
                        fputc( '$', stdout );
                        fputc( '\x00', stdout );
                    }
                    break;

                    case _CTLOG_TYPE_BOOL: fprintf( stdout, "%"PRIu8, (uint8_t)va_arg( vl, int ) ); break;
                    case _CTLOG_TYPE_CHAR: fprintf( stdout, "%"PRIu8, (uint8_t)va_arg( vl, int ) ); break;
                    default: assert( false ); break;
                }

                putchar( ',' );
            }
            va_end( vl );
        }

        fputs( "\n", stdout );
    }

    ctlog_sequenceNumber += 1;
}


/*============================================================================*/
void
ctlog_flush( void )
{
    fflush( stdout );
}
