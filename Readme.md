# Gost

File System abstraction layer for Golang, that works with Local file system 
and Amazon S3 with a unified API. FTP, Dropbox etc. will follow soon.

```
%[[Build Status](https://travis-ci.org/usmanhalalit/gost.svg?branch=master)](https://travis-ci.org/usmanhalalit/gost)
```


Quick Example:

% Maybe add aGIF?

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
err := fs.File("test.txt").Write("sample content")
```

Create and Delete

Listing

Information

