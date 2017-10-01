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
- Delete, ex: `delete 123` removes the row with the 123 id
- Delete  all rows, ex: `delete database` emptys the cache and removes the data file
- Open test database, ex `./database test` opens seperate new database for testing
- Caching

## Current Row Data Structure
| ID      | textLength | text           |
|:-------:|:----------:|:--------------:|
| 63 bits | 16 bits    | var bit length |

*note zero length text is not an option*

*1 bit in ID bit for identification*

## To Do
- Test Suite

## Feature Ideas
- Load length of text into memory to optimize crawling data
- Text Compression
- Improved paging, to be optomized for the user's PC page size
- BST trees for efficiency
- Fork and create a db server
