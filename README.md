# flshell

### Usage: ./flshell imagename offset(in sectors)

Use arrow keys to navigate (up/down to traverse through files in current directory) and (left/right to go up a level or into the selected directory).
ENTER will also go into the selected directory.

### Features
1. Hit TAB on a file to dump it and it's MFT information to the current directory.
2. View a disk image as if it were mounted, without having to mount it. 

### Known bugs
1. Unable to enter unallocated directories.

### Coming features
Dump directories utilising tsk recover
Confirmation prior to dumping files

