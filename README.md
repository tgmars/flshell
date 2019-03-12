# flshell

### Usage: ./flshell imagename offset(in sectors)

Use arrow keys to navigate (up/down to traverse through files in current directory) and (left/right to go up a level or into the selected directory).

ENTER will also go into the selected directory.

### Features
1. Hit TAB on a file to dump it and it's MFT information to the current directory.
2. View a disk image as if it were mounted, without having to mount it. 

### Known bugs
1. Unable to enter unallocated directories.
2. Scrolling; currently the printer goes to max of the received string to print,eg 40. It matches the current line
being printed with the 'selectedstring' and highlights that. However, selectedstring has freedom of movement throughout the entire directory.
Thus, you get a situation where upon scrolling, currentline being printed will never equal the selected string, past its maximum value.

### Coming features
Dump directories utilising tsk recover
Confirmation prior to dumping files

