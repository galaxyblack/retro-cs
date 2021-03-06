# monitor

Enable the monitor by using -m on the command line.

*NOTE*: This document needs an update as it is no longer correct.

## Arguments

The arguments for *address* and *value* are decimal values, or other values using the following prefixes:

    - `$` or `0x`: hexadecimal
    - `%` or `0b`: binary

Examples:
```
monitor> poke $1234 10
monitor> poke $1234 $0a
monitor> poke $1234 %1010
```

## Conversions
Typing in a number at the monitor prompt will show the value in decimal,
hexadecimal, and binary.

Examples:
```
monitor> 129
129 $81 %10000001
```

## Commands

### b[reak] [list]

List all active breakpoint addresses.

### b[reak] clear *address*

Clear the breakpoint at *address*.

### b[reak] clear-all

Clear all breakpoints.

### b[reak] set *address*

Set a breakpoint at *address*. The CPU will be stopped before executing the instruction at this address.

### cpu

Show the CPU status (registers and flags)

### cpu reg

List available registers

### cpu reg *name*

Show the value for the register with the given *name*

### cpu reg *name* *value*

Set the *value* for the register with the given *name*

### cpu flag

List available flags

### cpu flag *name*

Show the value for the flag with the given *name*

### cpu flag *name* *value*

Set the *value* for the flag with the given name

### d[asm] [list] [*start_address*] [*end_address*]

Disassemble code from *start_address* to *end_address*. If *end_address* is not specified, disassemble an amount specified with the `dasm lines` command. If *start_address* is not specified, continue from the last disassembly.

### dasm lines

Show the number of lines disassembled when an end address is not specified. A value of 0 means to disassemble an amount of lines that fit on the screen.

### dasm lines *count*

Set the number of lines disassembled to *count* when an end address is not specified. A value of 0 means to disassemble an amount of lines that fit on the screen.

### g[o]

Go. Start execution of the processors.

### load [*name*]

Load state that was saved with the `save` command with the given *name*. If name isn't specified, `state` is used.

### m[em] [dump] [*start_address*] [*end_address*]

Dump memory from *start_address* to *end_address*. If *end_address* is not specified, sump an amount specified with the `mem lines` command. If *start_address* is not specified, continue from the last dump.

### mem encoding

List the character encodings available for display when dumping memory.

### mem encoding *name*

Set the character encoding, with the given *name*, used when dumping memory.

### mem fill *start_address* *end_address* *value*

Fill memory from *start_address* to *end_address* with *value*.

### mem lines

Show the number of lines dumped when an end address is not specified. The default value is to dump a page.

### mem lines *count*

Set the number of lines dumped to *count* when an end address is not specified.

### p[ause]

Pause the execution of all processors.

### poke *address* *value*

Set the memory *address* with the given *value*

### peek *address*

Show the memory value at *address*

### q[uit]

Exit.

### r

Display the CPU status (registers and flags)

### save [*name*]

Save the current state with the given *name*. If *name* is not specified, `state` is used. Use load to restore to this state.

### t[race]

Toggle the tracing of instruction execution.

### w[atch] [list]

Watch all active memory watches.

### w[atch] clear *address*

Clear the watch at *address*.

### w[atch] clear-all

Clear all watches.

### w[atch] set *address* *mode*

Set a watch for *address*. Mode is either *r* for reads, *w* for writes, *rw* for reads and writes.

### x

Stop all logging output.

