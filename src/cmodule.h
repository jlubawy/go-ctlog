/**
 * Framework for working with C modules. A module is defined as a C source
 * file (*.c) that has a unique filename relative to other source files.
 * For example, 'gpio.c' and 'gpio_mcu_abc.c' are two
 * distinct modules with the names 'gpio' and 'gpio_mcu_abc'
 * respectively. See the macros below for how to define and use modules.
 *
 * Rather than store the name of each module name as a string in the firmware
 * (which uses a lot of code space), we can assign an index to each module
 * and do a lookup of the module index to the module name (either manually
 * or using development tools). See tlog.h for an example of how this can be used.
 *
 * Also refer to the cmodule build tool to see how the module indices are
 * generated (basically we sort the module names alphabetically and assign the
 * index based on that order).
 */

#ifndef CMODULE_H
#define CMODULE_H

#include <stdint.h>

#include "cmodule_indices.h"

/*==============================================================================
 *                                   Defines
 *============================================================================*/
/*============================================================================*/


/*==============================================================================
 *                                   Macros
 *============================================================================*/
/*============================================================================*/
// Helper macro to get the module index from the module name.
#define CMODULE_GET_INDEX( _name )  CMODULE_INDEX_ ## _name

/*============================================================================*/
// Helper macro to define a firmware module. The cmodule build tool searches
// *.c files for this macro and extract the _name argument. The _name argument
// must match the filename minus the *.c extension. For example, to create a new
// module within the file 'gpio_mcu_abc.c' you would place the following
// at the top of the file:
//
//     CMODULE_DEFINE( gpio_mcu_abc );
//
// Once a module is defined the static variable 'g_cmodule_index' can be
// referenced in other macros, eliminating the need to modify everywhere that
// references the module index if the filename changes.
//
#define CMODULE_DEFINE( _name ) \
            static cmodule_index_t g_cmodule_index = CMODULE_GET_INDEX( _name )


/*==============================================================================
 *                                   Types
 *============================================================================*/
/*============================================================================*/
// A cmodule_index_t is defined as a fixed-width uint32_t so that the index
// works reliably across different platforms.
typedef uint32_t cmodule_index_t;


/*==============================================================================
 *                                   Globals
 *============================================================================*/
/*============================================================================*/


/*==============================================================================
 *                               Public Functions
 *============================================================================*/
/*============================================================================*/


#endif /* CMODULE_H */
