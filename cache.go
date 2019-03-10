package main

import (
	"sort"
	"strings"
)

// Entry ... Interface for any object part of a file system."
type Entry interface {
	addChild()
	hasChildren()
	sortChildrenByAlphaDescending()
	listChildren()
}

// Item ... Basic struct for all items in file system.
type Item struct {
	Type   string
	Inode  string
	Name   string
	Parent *Folder
}

// Folder ... extension of Item. Has a single parent folder & a slice of child folders and files.
type Folder struct {
	Item
	Children []*Item
}

// addFolder ... Append a child folder to a current folder object.
// Assign the parent of the new folder as the current folder.
func (f *Folder) addChild(newItem *Item) {
	f.Children = append(f.Children, newItem)
	newItem.Parent = f
}

// hasChildren ... Returns true if the specified pointer to a Folder
// object has either ChildFiles or ChildFolders.
func (f *Folder) hasChildren() bool {
	if len(f.Children) > 0 {
		return true
	}
	return false
}

// sortChildrenByAlphaDescending ... Sorts the Children slice of a Folder type
// by alphabetical descending order (case insensitive). If matches are the same, falls back
// to the Inode value.
func (f *Folder) sortChildrenByAlphaDescending() {
	var children = f.Children
	//Assign
	sort.Slice(children, func(i int, j int) bool {
		if strings.ToLower(children[i].Name) < strings.ToLower(children[j].Name) {
			return true
		}
		if strings.ToLower(children[i].Name) > strings.ToLower(children[j].Name) {
			return false
		}
		return children[i].Inode < children[j].Inode
	})
}

// listChildren ... Return a string containing the
// children of a Folder specified by a pointer to it. Returns in
// a format that mimics FLS output.
func (f *Folder) listChildren() string {
	var msg strings.Builder
	for _, child := range f.Children {
		msg.WriteString(child.Type + " " + child.Inode + ":\t" + child.Name + "\n")
	}
	return msg.String()

}
