# FLShell

### Usage: ./flshell -d=imagename -o=offset(in sectors)

```
Currently doesn't operate as advertised, see previous version for an operational version.
FLShell is being rewritten to use tcell rather than termbox-go, to improve modularity and code quality. 
```

FLShell is designed to decrease the time taken to navigate through a disk image using The Sleuthkit's fls tool by providing an interactive shell.

Use arrow keys to navigate (up/down to traverse through files in current directory) and (left/right to go up a level or into the selected directory).

ENTER will also go into the selected directory.

### Features
1. Hit TAB on a file to dump it and it's MFT information to the current directory.
2. View a disk image as if it were mounted, without having to mount it. 

### Known bugs
1. The FLShell executable might not work when downloading straight from Github. Instead, install Go and run 

```
go get -u github.com/tgmars/flshell
cd your-go-workspace/src/github.com/tgmars/flshell
go build ./...
```

2. Unable to enter unallocated directories - potentially a future feature.

### Coming features/To do
1. listChildren() to return on a string slice rather than string to save effort & memory expanding functionality.
2. Gather more information on current directory and display it back to the user (timestamps, filesize?).
3. Search for files by timestamps.
4. Add method to Item struct to get a 'pwd' equivalent.


