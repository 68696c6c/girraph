package filesystem

import (
	"fmt"

	"github.com/68696c6c/girraph"
)

type File interface {
	SetName(string) File
	GetName() string
	SetExtension(string) File
	GetExtension() string
	GetFullName() string
	SetContents(string) File
	GetContents() string
}

type file struct {
	Name      string
	Extension string
	Contents  string
}

func (f *file) SetName(name string) File {
	f.Name = name
	return f
}

func (f *file) GetName() string {
	return f.Name
}

func (f *file) SetExtension(ext string) File {
	f.Extension = ext
	return f
}

func (f *file) GetExtension() string {
	return f.Extension
}

func (f *file) GetFullName() string {
	return fmt.Sprintf("%s.%s", f.Name, f.Extension)
}

func (f *file) SetContents(contents string) File {
	f.Contents = contents
	return f
}

func (f *file) GetContents() string {
	return f.Contents
}

func MakeFile(name, extension string) *file {
	return &file{
		Name:      name,
		Extension: extension,
	}
}

type Directory interface {
	SetName(string) Directory
	GetName() string
	SetFiles([]*file) Directory
	GetFiles() []*file
}

type directory struct {
	Name  string
	Files []*file
}

func (d *directory) SetName(name string) Directory {
	d.Name = name
	return d
}

func (d *directory) GetName() string {
	return d.Name
}

func (d *directory) SetFiles(files []*file) Directory {
	d.Files = files
	return d
}

func (d *directory) GetFiles() []*file {
	return d.Files
}

func MakeDirectory(name string) girraph.Tree[Directory] {
	return girraph.MakeTree[Directory]().SetMeta(&directory{
		Name:  name,
		Files: []*file{},
	})
}
