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
- Insert, ex: `insert tablename 123 hello world` inserts a row with id:123 and text:"hello world"
- Select, ex: `select tablename 123` prints hello world
- Delete, ex: `delete tablename 123` removes the row with id equal to 123
- Delete  database, ex: `delete database confirm` emptys the cache and removes the data file
- Create table, ex: `create tablename` creates a table with that name
- Delete table, ex: `delete tablename` deletes a table with that name
- Caching
- External Api

## To Do
- Test Suite

## Feature Ideas
- Option for table with fixed length text
- Text Compression
- BST trees for efficiency
- Use in another project to create a db server

## Current Data Structure
Metadata Row

| ID    | length | data offset |
|:-----:|:------:|:-----------:|
| uin32 | uint16 | uint32      |


*all offset are from the end of metadata table* | *length doesn't include 0 value* | *zero is an invalid ID*

At the end of the metadata table

| next metatable offset |
|:---------------------:|
| uint32                |


The text data is inbetween the current metatable and the next metatable

## Previous Data structure
concurrent rows of:

| row identifier | ID      | text length | text           |
|:--------------:|:-------:|:-----------:|:--------------:|
| 1 bit          | 63 bits | 16 bits     | var bit length |


*note: text length in bytes, zero length text is not an option*

## Road Map
I have had 2 ideas to optomize my database
- BST
- Seperate the Index from the data

For those famillair with computer science BSTs are a pretty obvious and efficient improvement to a DB.
Seperating the Index from the data was an idea I came up with on my own that I'm interested in exploring.
A Row while be split into seperate tables: metadata and data. The metadata table will be the size of the OS's page with an offset in the last spot for the next metadata table.

Pros
- No need for in memory pager, Metadata rows describe empty locations
- Reduced number of reads from disk, as metadata table includes all critical data for navigating the data table.
- Deletion consists exclusively of removing the Metadata
- Could use BST in Metatable if text length was fixed length

Cons
- Forms a complex relationship between the OS page size, the max length of the variable size text, data offset and next metadata table offset, ex: # of rows in metadata x row data length (must be) < next metatable offset or the next metadata table would be unreachable.
- Inserting is still pretty intensive as it requires looping through the whole table to find holes and requires 2 writes.
- Cannot figure out how to use BSTs with variable length text.
