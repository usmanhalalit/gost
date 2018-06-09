# Gost

[![Build Status](https://travis-ci.org/usmanhalalit/gost.svg?branch=master)](https://travis-ci.org/usmanhalalit/gost)

Filesystem abstraction layer for Golang, that works with Local file system 
and Amazon S3 with a unified API. You can even copy-paste files from different sources.
FTP, Dropbox etc. will follow soon.


Quick Example:

```
// Initialize a filesystem
fs := gost.s3.New(Config{
	ID: "aws-id",
	Key: "aws-key",
	Region: "es-west-1",
})

// Read
note := fs.File("my-note.txt").ReadString()
//Write
fs.File("another-note.txt").WriteString("another note")

// Traverese natuarally
movies := fs.Directory("movies")
files := movies.Files()
movies.File("Pirated-movie.mp4").Delete()

// Copy file from one source to another
localFile := lfs.File("photo.jpg")
s3Dir := fs.Directory("photos")
err := localFile.CopyTo(s3dir)
```

### Table of Contents
  * [Initialize](#initialize)
    + [S3](#s3)
    + [Local](#local)
  * [Read and Write](#read-and-write)
    + [Read](#read)
    + [Write](#write)
  * [Traversing](#traversing)
  * [Listing](#listing)
  * [Stat](#stat)
  * [Create and Delete](#create-and-delete)
  * [Copy and Paste Between Different Sources](#copy-and-paste-between-different-sources)
  * [Custom Adapter](#custom-adapter)


## Initialize

You just initialize the S3 and Local adapters differently, **everything else in the API is same**.

```bash
go get github.com/usmanhalalit/gost
``` 

### S3
```
fs := gost.s3.New(gost.s3.Config{
	ID: "aws-id",
	Key: "aws-key",
	Region: "es-west-1",
})
```

### Local
fs := gost.local.New(gost.local.Config{
	BasePath: "/home/user",
})

## Read and Write

### Read
Simple read, suitable for small files.

```
fileContent, err := fs.File("test.txt").ReadString()
```

Bytes read, compatible with `io.Reader`, so you can do buffered read.
```
b := make([]byte, 3)
n, err := fs.File("test.txt").Read(b)
```

### Write
Simple write
```
fs.File("test.txt").WriteString("sample content")
```

Byte write
```
n, err := file.Write(bytes)
// n == number of bytes written
```

## Traversing

You can explore the filesystem like you in your desktop file explorer.
File and directories are chained in a natural way. 

```
dirs, err := fs.Directory("Parent").Directory("Child").Directories()
files, err := fs.Directory("Parent").Directory("Child").Files()
```

```
dirs, err := fs.Directory("Parent").Directtory("Child").Files()
```

## Listing

Get all files and loop through them
```
files, err := fs.Directory("Parent").Directory("Child").Files()
for _, file := range files {
    fmt.Println(file.ReadString())
}
```
Get all directories and loop through them
```
dirs, err := fs.Directories()
for _, dir := range dirs {
    files := dir.Files()
    fmt.Println(files)
}
```

## Stat

Get file size and last modified timestamp:

```
stat, _ := fs.File("test.txt").Stat()
fmt.Println(stat.Size)
fmt.Println(stat.LastModifed)
```

You can get stat of directories too, but it's not available on S3.

```
fs.Directory("Downloads").File("test.txt").GetPath()
```


## Create and Delete
Delete a file and directory:
```
fs.File("test.txt").Delete()
// Delete an entrie directory, beware please!
fs.Directory("Images").Delete()
```

Create a new directory:
```
fs.Directory("Images").Create()
```

To create a new file simply write something to it:
```
fs.File("non_existant_file").WriteString("")
```  

## Copy and Paste Between Different Sources

You can copy a file to any Directory, be it in in the same filesystem or not(local or S3)

```
localFile := lfs.File("photo.jpg")
s3Dir := s3fs.Directory("photos")
err := localFile.CopyTo(s3dir)
``` 

Fun, eh? 

You can optionally provide a new filename too:
```
err := localFile.CopyTo(anotherDir, "copied_file.jpg")
```

Also there is a helper to copy file in the same Directory:
```
file.Copy("copied_file.jpg")
``` 
 

## Custom Adapter

Yes, you can write one and it'll be appreciated if you contribute back.
. `gost.go` file has all the interfaces defined. Basically you've to implement
`gost.File` and `gost.Directory` interfaces. Check the `local` adapter to get an idea. 

