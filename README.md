# GoLite
GoLite is a basic SQLite replica in Go

## To replicate

run the following cmds inside the golite folder

To build:
```
go build
```

To run:
```
./database
```

## Current Features
- Repl
- Insert, ex: `insert 123 hello world`
- Select, ex: `select 123` prints hello world
- Caching

## ToDo
- Test Suite
- Remove row

## Feature Ideas
- Load length of text into memory to optimize crawling data
- Text Compression
- Improved paging, to be optomized for the user's PC page size
- BST trees for efficiency
- Fork and create a db server
