# flshell

### Usage: ./flshell imagename offset(in sectors)

FLShell is designed to decrease the time taken to navigate through a disk image using The Sleuthkit's fls tool by providing an interactive shell.

Use arrow keys to navigate (up/down to traverse through files in current directory) and (left/right to go up a level or into the selected directory).

ENTER will also go into the selected directory.

### Features
1. Hit TAB on a file to dump it and it's MFT information to the current directory.
2. View a disk image as if it were mounted, without having to mount it. 

### Known bugs
1. The flshell executable might not work when downloading straight from Github. Instead, install Go and run 

```
go get -u github.com/tgmars/flshell
cd your-go-workspace/src/github.com/tgmars/flshell
go build ./...
```

2. Unable to enter unallocated directories.

### Coming features/To do
1. Dump directories utilising tsk recover
2. Do error handling
3. Confirmation prior to dumping files
4. Gather more information on current directory and display it back to the user (timestamps, filesize?)
5. Huge code cleanup 

