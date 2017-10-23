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
- Insert, ex: `insert 123 hello world` inserts a row with id:123 and text:"hello world"
- Select, ex: `select 123` prints hello world
- Delete, ex: `delete 123` removes the row with id equal to 123
- Delete  database, ex: `delete database` emptys the cache and removes the data file
- Open test database, ex `./database test` opens seperate new database for testing
- Caching
- External Api

## Current Row Data Structure
| row identifier | ID      | text length | text           |
|:--------------:|:-------:|:-----------:|:--------------:|
| 1 bit          | 63 bits | 16 bits     | var bit length |

*note: text length in bytes, zero length text is not an option*

## To Do
- Test Suite

## Feature Ideas
- Load length of text into memory to optimize crawling data
- Text Compression
- Improved paging, to be optomized for the user's PC page size
- BST trees for efficiency
- Fork and create a db server

## Road Map
I have had 2 ideas to optomize my database that can be combined
- BST
- Seperate the Index from the data

For those famillair with computer science BSTs are a pretty obvious and efficient improvement to a DB.
Seperating the Index from the data was an idea I came up with on my own that I'm interested in exploring.
A Row while be split into seperate tables: metadata and data. The metadata table will be the size of the OS's page with an offset in the last spot for the next metadata table.

Metadata Row
| ID    | length | data offset |
|:-----:|:------:|:-----------:|
| uin32 | uint16 | uint32      |

*all offset are from the end of metadata table* | *length doesn't include 0 value* | *zero is an invalid ID*

At the end of the metadata table
| next metatable offset |
|:----------------------:|
| uint32 |

Data Row
| Text            |
|:---------------:|
| variable length |

Pros
- No need for in memory pager, Metadata rows describe empty locations
- Reduced number of reads from disk, as metadata table includes all critical data for navigating the data table.
- Deletion could consist exclusively of removing the Metadata
- Still allows for use of BSTs (would be used in metadata table)

Cons
- Forms a complex relationship between the OS page size, the max length of the variable size text, data offset and next metadata table offset, ex: # of rows in metadata x row data length (must be) < next metatable offset or the next metadata table would be unreachable.
- Inserting is still pretty intensive as it requires looping through the whole table to find holes
