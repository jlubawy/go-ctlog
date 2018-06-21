# examples/arduino

This example is a simple program that can be run on an Arduino Uno. The program
reads and echos lines over the UART. To demonstrate the memory savings that
tokenized logging provides we write an arbitrary ~1kB string to stdout.

In a real-life scenario the arbitrary string wouldn't exist, but we might have
other actually useful format strings used to output runtime data to help with
debugging and knowing the status of our device.

## TL;DR

With tokenized logging enabled we save 886 bytes of flash memory used
for data (the string), and we only add 144 bytes of flash memory used for the
program. __The net savings are 742 bytes of flash memory.__ When the strings
start taking up multiple kB of memory the savings become even more significant.

## Without tokenized logging

With tokenized logging disabled, this string obviously uses up the ~1kB of flash
memory needed to store the string.

Here is the memory usage report:

    AVR Memory Usage
    ----------------
    Device: atmega328p

    Program:    4664 bytes (14.2% Full)
    (.text + .data + .bootloader)

    Data:       1310 bytes (64.0% Full)
    (.data + .bss + .noinit)

And here is the output:

    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris lacus ligula, ultrices sed condimentum ac, aliquet in nulla. Ut lobortis pulvinar dui, auctor consectetur nulla. Suspendisse id malesuada neque. Cras pretium nisl quis felis hendrerit tristique. Phasellus sed porttitor dui. Phasellus aliquam fermentum elit at aliquet. Nullam porta, tortor vitae sagittis dapibus, felis libero dictum nunc, eu tincidunt orci diam at risus. Donec quis bibendum turpis. Maecenas ultrices imperdiet nulla non laoreet. Sed euismod rhoncus lorem, porttitor varius nunc tempus at. Aenean dignissim fringilla dui ac commodo. Mauris iaculis et ipsum id malesuada. Donec semper magna a malesuada dictum. Aliquam vehicula ligula vitae venenatis elementum. Phasellus congue eleifend viverra. Suspendisse potenti. Fusce aliquet, massa ac tristique egestas, dui tellus molestie mi, quis accumsan lacus eros quis tellus. Nulla ipsum nulla, dapibus in purus sed, pellentesque volutpat tortor. Aliquam tincidunt interdum arcu ac maximus.
    > Hello world
    Hello world

## With tokenized logging

When tokenized logging is enabled, the string is no longer stored on the device
so we save the ~1kB that the string previously used.

Here is the memory usage report:

    AVR Memory Usage
    ----------------
    Device: atmega328p

    Program:    4808 bytes (14.7% Full)
    (.text + .data + .bootloader)

    Data:        424 bytes (20.7% Full)
    (.data + .bss + .noinit)

And here is the output:

    {"ctlog":0,"seq":0,"lvl":"I","mi":0,"ml":101,"args":[]}
    > Hello world
    {"ctlog":0,"seq":1,"lvl":"I","mi":0,"ml":154,"args":[{"t":3,"v":"Hello world"}]}

This output can then be used to lookup the string using the module index and line
number to reproduce the original output:

```json
{
  "date": "2018-06-21T21:55:45.537831Z",
  "modules": [
    {
      "index": 0,
      "name": "main",
      "path": "./go-ctlog/examples/arduino/src/main.c",
      "lines": [
        {
          "number": 101,
          "formatString": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris lacus ligula, ultrices sed condimentum ac, aliquet in nulla. Ut lobortis pulvinar dui, auctor consectetur nulla. Suspendisse id malesuada neque. Cras pretium nisl quis felis hendrerit tristique. Phasellus sed porttitor dui. Phasellus aliquam fermentum elit at aliquet. Nullam porta, tortor vitae sagittis dapibus, felis libero dictum nunc, eu tincidunt orci diam at risus. Donec quis bibendum turpis. Maecenas ultrices imperdiet nulla non laoreet. Sed euismod rhoncus lorem, porttitor varius nunc tempus at. Aenean dignissim fringilla dui ac commodo. Mauris iaculis et ipsum id malesuada. Donec semper magna a malesuada dictum. Aliquam vehicula ligula vitae venenatis elementum. Phasellus congue eleifend viverra. Suspendisse potenti. Fusce aliquet, massa ac tristique egestas, dui tellus molestie mi, quis accumsan lacus eros quis tellus. Nulla ipsum nulla, dapibus in purus sed, pellentesque volutpat tortor. Aliquam tincidunt interdum arcu ac maximus."
        },
        {
          "number": 154,
          "formatString": "line=%s"
        }
      ]
    }
  ]
}
```
