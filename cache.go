package main

import (
	"sort"
	"strings"
)

// Item ... Basic struct for all items in file system.
type Item struct {
	Type     string
	Inode    string
	Name     string
	Parent   *Item
	Children []*Item
}

// addFolder ... Append a child folder to a current folder object.
// Assign the parent of the new folder as the current folder.
func (f *Item) addChild(newItem *Item) {
	f.Children = append(f.Children, newItem)
}

// hasChildren ... Returns true if the specified pointer to a Folder
// object has either ChildFiles or ChildFolders.
func (f *Item) hasChildren() bool {
	if len(f.Children) > 0 {
		return true
	}
	return false
}

// sortChildrenByAlphaDescending ... Sorts the Children slice of a Folder type
// by alphabetical descending order (case insensitive). If matches are the same, falls back
// to the Inode value.
func (f *Item) sortChildrenByAlphaDescending() {
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
func (f *Item) listChildren() string {
	var msg strings.Builder
	for _, child := range f.Children {
		msg.WriteString(child.Type + " " + child.Inode + ":\t" + child.Name + "\n")
	}
	return msg.String()

}

// populate ... Returns a pointer to an Item with children populated by passing
// a message delimited by newlines.
func (f *Item) populate(input string) *Item {
	//Split input on newlines and assign to a slice of strings
	s := strings.Split(input, "\n")
	//For each line in the split, add it as a child to the root folder.
	for _, line := range s {
		// Only attempt to parse values out of the line if there's stuff in it.
		if line != "" {
			itemType := dirMatcher(line)
			itemInode := inodeMatcher(line)
			itemName := nameMatcher(line)
			f.addChild(&Item{itemType, itemInode, itemName, f, nil})
		}
	}
	return f
}

// goUp ... sets the current item as the parent stored in the struct.
func (f *Item) goUp(parent Item) *Item {
	f = parent.Parent
	return f
}

// goDown ... sets the current item as the child specified by parameters.
func (f *Item) goDown(parent Item, itemType string, inode string) *Item {
	for _, child := range f.Children {
		if (child.Inode == inode) && (child.Type == itemType) {
			f = child
			f.Parent = &parent
			return f
		}
	}
	return nil
}
