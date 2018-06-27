# examples/basic

This example is a simple program that can be run on a PC, it demonstrates the
difference between tokenized logging on and off.

```c
#include <stdbool.h>
#include <stdio.h>

int
main( void )
{
    fputs( "Long string\n", stdout );
    printf( "%u\n", 123 );
    printf( "%u\n", 456 );
    printf( "%u\n", 789 );
    printf( "%d\n", -123 );
    printf( "%d\n", -456 );
    printf( "%d\n", -678 );
    printf( "%s\n", "Hello\tWorld" );
    printf( "%s\n", (true) ? "true" : "false" );
    printf( "%c\n", 'J' );
    return 0;
}
```

With the output:

    Long string
    123
    456
    789
    -123
    -456
    -678
    Hello   World
    true
    J

The tokenized logging equivalent looks like the following:

```c
#include <stdbool.h>
#include <stdio.h>

#include "ctlog.h"

CMODULE_DEFINE( main );

int
main( void )
{
    ctlog_setStream( stdout );
    CTLOG_INFO( "Long string" );
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_UINT( 123 ) );
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_UINT( 456 ) );
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_UINT( 789 ) );
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_INT( -123 ) );
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_INT( -456 ) );
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_INT( -678 ) );
    CTLOG_VAR_INFO( "%s", 1, CTLOG_TYPE_STRING( "Hello\tWorld" ) );
    CTLOG_VAR_INFO( "%t", 1, CTLOG_TYPE_BOOL( true ) );
    CTLOG_VAR_INFO( "%c", 1, CTLOG_TYPE_CHAR( 'J' ) );
    return 0;
}
```

The raw output looks like the following:

    {"ctlog":0,"seq":0,"lvl":"I","mi":0,"ml":13,"args":[]}
    {"ctlog":0,"seq":1,"lvl":"I","mi":0,"ml":14,"args":[{"t":4,"v":123}]}
    {"ctlog":0,"seq":2,"lvl":"I","mi":0,"ml":15,"args":[{"t":4,"v":456}]}
    {"ctlog":0,"seq":3,"lvl":"I","mi":0,"ml":16,"args":[{"t":4,"v":789}]}
    {"ctlog":0,"seq":4,"lvl":"I","mi":0,"ml":17,"args":[{"t":2,"v":-123}]}
    {"ctlog":0,"seq":5,"lvl":"I","mi":0,"ml":18,"args":[{"t":2,"v":-456}]}
    {"ctlog":0,"seq":6,"lvl":"I","mi":0,"ml":19,"args":[{"t":2,"v":-678}]}
    {"ctlog":0,"seq":7,"lvl":"I","mi":0,"ml":20,"args":[{"t":3,"v":"Hello\tWorld"}]}
    {"ctlog":0,"seq":8,"lvl":"I","mi":0,"ml":21,"args":[{"t":0,"v":true}]}
    {"ctlog":0,"seq":9,"lvl":"I","mi":0,"ml":22,"args":[{"t":1,"v":"J"}]}

This JSON output can then be run through the ```ctlog``` tool like so:

    ./examples/basic/main_ctlog | ctlog log examples/basic/ctlog_dict.json

Which produces identical output as the original program:

    Long string
    123
    456
    789
    -123
    -456
    -678
    Hello   World
    true
    J
