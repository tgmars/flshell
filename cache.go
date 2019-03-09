package main

import "fmt"

// Entry ... Interface for any object part of a file system."
type Entry interface {
}

// Item ... Basic struct for all items in file system.
type Item struct {
	name  string
	inode string
}

// Folder ... extension of Item. Has a single parent folder & a slice of child folders and files.
type Folder struct {
	Item
	Parent       *Folder
	ChildFolders []*Folder
	ChildFiles   []*File
	NumItems     int
}

// File ... extension of Item to represent a file. Should be extended to include data from istat such as a creation timestamp.
type File struct {
	Item
	Parent *Folder
}

// AddFolder ... Append a child folder to a current folder object.
// Increase the count of the number of childen - this is probably redundant.
// Assign the parent of the new folder as the current folder.
func (f *Folder) addfolder(newfolder *Folder) {
	f.ChildFolders = append(f.ChildFolders, newfolder)
	f.NumItems++
	newfolder.Parent = f
}

func (f *Folder) addfile(newfile *File) {
	f.ChildFiles = append(f.ChildFiles, newfile)
	newfile.Parent = f
}

/*
func main() {

	foo := File{Item{"name:a.txt", "inode:4"}, nil}
	bar := Folder{Item{"name:User", "inode:2"}, nil, nil, nil, 0}

	fmt.Println("Bar:\t", bar)
	fubar := Folder{Item{"name:Documents", "inode:3"}, nil, nil, nil, 0}
	fubar3 := Folder{Item{"name:Pictures", "inode:5"}, nil, nil, nil, 0}

	bar.addfolder(&fubar)
	bar.addfolder(&fubar3)
	fubar3.addfile(&foo)
*/
}
