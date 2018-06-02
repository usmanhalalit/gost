# Gost

File System abstraction layer for Golang, that works with Local file system 
and Amazon S3 with a unified API. FTP, Dropbox etc. will follow soon.

```
%[[Build Status](https://travis-ci.org/usmanhalalit/gost.svg?branch=master)](https://travis-ci.org/usmanhalalit/gost)
```


Quick Example:

% Maybe add a GIF?

```
fs := gost.s3.New(Config{
	ID: "aws-id",
	Key: "aws-key",
	Region: "es-west-1",
})

note := fs.File("my-note.txt").ReadString()
fs.File("another-note.txt").WriteString("another note")

movies := fs.Directory("movies")
files := movies.Files()
movies.File("Pirated-movie.mp4").Delete()
```

## Initialize

Everything is same, you just initialize the adapters differently.

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
Simple read, suitable for small files.
```
fileContent, err := fs.File("test.txt").ReadString()
```

Bytes read, compatible with `io.Reader`
```
b := make([]byte, 3)
n, err := fs.File("test.txt").Read(b)
```

```
fs.File("test.txt").Write("sample content")
```

## Traversing

Chained in a natural was 

```
dirs, err := fs.Directory("Parent").Directory("Child").Directories()
files, err := fs.Directory("Parent").Directory("Child").Files()
```

```
dirs, err := fs.Directory("Parent").Directtory("Child").Files()
```

## Listing

```
files, err := fs.Directory("Parent").Directory("Child").Files()
for _, file := range files {
    fmt.Println(file.ReadString())
}
```

```
dirs, err := fs.Directories()
for _, dir := range dirs {
    files := dir.Files()
    fmt.Println(files)
}
```

## Stat

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

```
fs.File("test.txt").Delete()
// Delete an entrie directory, beware please!
fs.Directory("Images").Delete()
```

```
fs.Directory("Images").Create()
```

## Copy and Paste Between Different Sources
```
localFile = lfs.File("photo.jpg")
b, err := ioutil.ReadAll(f)
n, err := s3fs.File("photo.jpg").Write(b)
```
Plans to automate this 

## Custom Adapter

## Testing and Mocking

## Contributions

