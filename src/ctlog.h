/**
 * Tokenized logging framework.
 */

#ifndef CTLOG_H
#define CTLOG_H

#include <stdbool.h>
#include <stdint.h>

#include "cmodule.h"

/*==============================================================================
 *                                   Defines
 *============================================================================*/
/*============================================================================*/
// Version tokenized logging lines in case we need to change the output format.
#define CTLOG_VERSION  "00"

/*============================================================================*/
// Logging levels. These definitions must not change or else it will break
// compatibility with existing libraries/tools that parse tokenized log streams.
// These also shouldn't be needed outside of this file, use the enable definitions
// below instead.
#define _CTLOG_LEVEL_ERROR_BIT   (0x00)
#define _CTLOG_LEVEL_ERROR_CHAR  'E'
#define _CTLOG_LEVEL_INFO_BIT    (0x01)
#define _CTLOG_LEVEL_INFO_CHAR   'I'
#define _CTLOG_LEVEL_DEBUG_BIT   (0x02)
#define _CTLOG_LEVEL_DEBUG_CHAR  'D'
#define _CTLOG_LEVEL_WARN_BIT    (0x03)
#define _CTLOG_LEVEL_WARN_CHAR   'W'

/*============================================================================*/
// Enable specific logging levels using these definitions.
#define CTLOG_LEVEL_ENABLE_ERROR  (1 << _CTLOG_LEVEL_ERROR_BIT)
#define CTLOG_LEVEL_ENABLE_INFO   (1 << _CTLOG_LEVEL_INFO_BIT)
#define CTLOG_LEVEL_ENABLE_DEBUG  (1 << _CTLOG_LEVEL_DEBUG_BIT)
#define CTLOG_LEVEL_ENABLE_WARN   (1 << _CTLOG_LEVEL_WARN_BIT)

#ifndef CTLOG_LEVELS_ENABLED
#define CTLOG_LEVELS_ENABLED  (CTLOG_LEVEL_ENABLE_ERROR | CTLOG_LEVEL_ENABLE_INFO | CTLOG_LEVEL_ENABLE_WARN)
#endif

/*============================================================================*/
// Type definitions used by ctlog_printf to identify what type a variadic
// value argument should be cast to.
#define _CTLOG_TYPE_BOOL    (0x00)
#define _CTLOG_TYPE_CHAR    (0x01)
#define _CTLOG_TYPE_INT     (0x02)
#define _CTLOG_TYPE_STRING  (0x03)
#define _CTLOG_TYPE_UINT    (0x04)

/*============================================================================*/
#define CTLOG_TYPE_BOOL( _val )    _CTLOG_TYPE_BOOL,   (uint8_t)(_val)
#define CTLOG_TYPE_CHAR( _val )    _CTLOG_TYPE_CHAR,   (uint8_t)(_val)
#define CTLOG_TYPE_INT( _val )     _CTLOG_TYPE_INT,    (int32_t)(_val)
#define CTLOG_TYPE_STRING( _val )  _CTLOG_TYPE_STRING, (_val)
#define CTLOG_TYPE_UINT( _val )    _CTLOG_TYPE_UINT,   (uint32_t)(_val)



/*==============================================================================
 *                                   Macros
 *============================================================================*/
/*============================================================================*/
// Helper macros for building the log level macros below. Not intended for use
// outside of this file.
#define CTLOG_BASE( _level, _nArgs, ... )  (ctlog_printf( _level, g_cmodule_index, __LINE__, _nArgs, __VA_ARGS__ ))
#define CTLOG_NO_ARGS( _level )            (CTLOG_BASE( _level, 0, NULL ))

/*============================================================================*/
#if (CTLOG_LEVELS_ENABLED & CTLOG_LEVEL_ENABLE_ERROR)
  #define CTLOG_ERROR( _str )                   (CTLOG_NO_ARGS( _CTLOG_LEVEL_ERROR_CHAR ))
  #define CTLOG_VAR_ERROR( _str, _nArgs, ... )  (CTLOG_BASE( _CTLOG_LEVEL_ERROR_CHAR, _nArgs, __VA_ARGS__ ))
#else
  #define CTLOG_ERROR( _str )
  #define CTLOG_VAR_ERROR( _str, _nArgs, ... )
#endif

/*============================================================================*/
#if (CTLOG_LEVELS_ENABLED & CTLOG_LEVEL_ENABLE_INFO)
  #define CTLOG_INFO( _str )                   (CTLOG_NO_ARGS( _CTLOG_LEVEL_INFO_CHAR ))
  #define CTLOG_VAR_INFO( _str, _nArgs, ... )  (CTLOG_BASE( _CTLOG_LEVEL_INFO_CHAR, _nArgs, __VA_ARGS__ ))
#else
  #define CTLOG_INFO( _str )
  #define CTLOG_VAR_INFO( _str, _nArgs, ... )
#endif

/*============================================================================*/
#if (CTLOG_LEVELS_ENABLED & CTLOG_LEVEL_ENABLE_DEBUG)
  #define CTLOG_DEBUG( _str )                   (CTLOG_NO_ARGS( _CTLOG_LEVEL_DEBUG_CHAR ))
  #define CTLOG_VAR_DEBUG( _str, _nArgs, ... )  (CTLOG_BASE( _CTLOG_LEVEL_DEBUG_CHAR, _nArgs, __VA_ARGS__ ))
#else
  #define CTLOG_DEBUG( _str )
  #define CTLOG_VAR_DEBUG( _str, _nArgs, ... )
#endif

/*============================================================================*/
#if (CTLOG_LEVELS_ENABLED & CTLOG_LEVEL_ENABLE_WARN)
  #define CTLOG_WARN( _str )                   (CTLOG_NO_ARGS( _CTLOG_LEVEL_WARN_CHAR ))
  #define CTLOG_VAR_WARN( _str, _nArgs, ... )  (CTLOG_BASE( _CTLOG_LEVEL_WARN_CHAR, _nArgs, __VA_ARGS__ ))
#else
  #define CTLOG_WARN( _str )
  #define CTLOG_VAR_WARN( _str, _nArgs, ... )
#endif

/*============================================================================*/
// When adding/changing new log macros keep in mind that some tools (e.g. tokenlog)
// use these macro names to create the tokenized log strings file. Make sure to
// update those tools as necessary.



/*==============================================================================
 *                               Public Functions
 *============================================================================*/
/*============================================================================*/
void
ctlog_setEnabled( bool enable );

/*============================================================================*/
bool
ctlog_isEnabled( void );

/*============================================================================*/
void
ctlog_printf( char level, cmodule_index_t moduleIndex, uint32_t line, int nArgs, ... );

/*============================================================================*/
void
ctlog_flush( void );


#endif
