# Book Manager

## Book Manager CLI 

`Usage: bm [COMMAND]... [OPTIONS]...`

| Command            	| Arguments            	| Options                                                                                               	| Output                                                    	| Error                                    	|
|--------------------	|----------------------	|-------------------------------------------------------------------------------------------------------	|-----------------------------------------------------------	|------------------------------------------	|
| add book           	| -isbn -title -author 	| -published -description -genre                                                                        	| book [title] successfully added                           	| - if isbn already exists                 	|
| edit book          	| -isbn                	| -title  -author -published -description -genre                                                        	| book [title] successfully updated                         	| - if isbn does not exist                 	|
| remove books       	|                      	| -isbn -title  -author -published -description -genre                                                  	| [# of books] successfully removed                         	|                                          	|
| detail book        	| -isbn                	|                                                                                                       	| [book details]                                            	| - if isbn does not exist                 	|
| add collection     	| -name                	|  -collection-description-books (comma separated isbns)                                                	| collection [name] successfully added with id [id]         	| - if collection already exists           	|
| view collection    	| -id                  	|                                                                                                       	| [collection details with a table of books]                	| - if collection id does not exist        	|
| edit collection    	| -id                  	| -remove-books (comma separated list of isbns) -add-books (comma separated list of isbns) -description 	| collection [name] successfully updated                    	| - if collection id does not exist        	|
| remove collection  	| -id                  	|                                                                                                       	| collection [name] successfully removed                    	| if collection with name does not exist   	|
| detail collection  	| -name                	|                                                                                                       	| [collection detail with table of books]                   	| - if collection with name does not exist 	|
| search books       	|                      	| -isbn -title  -author -published -description -genre                                                  	| [list of books with isbn, title, author, published date]  	| - if no search options are provided      	|
| search collections 	|                      	| -name -isbn -title -author -published -description -genre                                             	| [list of collections with name, # of books]               	| - if no search options are provided      	|

## Book Manager REST API
### Books
`HTTP POST /api/books`
```
[payload]
{
    "isbn":"",
    "title":"",
    "author":"",
    "description":"",
    "published_date":"",
    "genres":[]
}

[response]
200 OK

400 Bad Request
{"message":"isbn already exists"}
{"message":"title exceeds maximum 512 characters","field":"title"}
```

`HTTP DELETE /api/books/<isbn>/delete`
```
[response]
200 OK

400 Bad Request
{"message":"isbn not found"}
```

`HTTP GET /api/books/<isbn>/view`
```
[response]
{
    "isbn":"",
    "title":"",
    "author":"",
    "description":"",
    "published_date":"",
    "genres":[]
}
200 OK

400 Bad Request
{"message":"isbn not found"}
```

`HTTP PUT /api/books/<isbn>/edit`
```
[payload]
{
    "title":"",
    "author":"",
    "description":"",
    "published_date":"",
    "genres":[]
}

[response]
200 OK

400 Bad Request
{"message":"isbn not found"}
{"message":"author cannot be blank","field":"author"}
```

`HTTP GET /api/books?output=&title=&author=miller&description=pineapples&published_date=1988&genres=a,b,c`
```
[response]
200 OK
{
    "total":113,
    "results":
        [
            {
                "title":"",
                "author":"",
                "description":"",
                "published_date":"",
                "genres":[]
            },...
        ]
}

400 Bad Request
{"message":"no search terms provided"}
{"message":"too many results"}
```

### Collections

`HTTP GET /api/v1/collections`
```
[payload]
[
    {
       "name":"",
       "description":"",
       "books":[]
    },
    {
       "name":"",
       "description":"",
       "books":[]
    }...
]
[response]
200 OK
```

`HTTP POST /api/v1/collections`
```
[payload]
{
    "name":"",
    "description":"",
    "books":[]
}

[response]
200 OK

400 Bad Request
{"message":"a collection with this name already exists"}
{"message":"description exceeds maximum 512 characters","field":"title"}
{"message":"<isbn> does not exist","field":"books"}
```

`HTTP GET /api/v1/collections/<name>/view`
```
[response]

200 OK
{
    "name":"",
    "description":"",
    "books":[
        {
            "isbn":"",
            "title":"",
            "author":"",
            "published_date":""
        }
    ]
}

400 Bad Request
{"message":"this collection does not exist"}
```
`HTTP PUT /api/v1/collections/<name>/edit`
```
[payload]
{
    "name":"",
    "description":""
}

[response]
200 OK

400 Bad Request
{"message":"this collection does not exist"}
{"message":"description exceeds maximum length of 512"}
```

`HTTP POST /api/v1/collections/<name>/add-books`
```
[payload]
{
    "books":[]
}

[response]
200 OK

400 Bad Request
{"message":"this collection does not exist"}
{"message":"<isbn> does not exist"}
{"message":"no books provided"}
```

`HTTP POST /api/v1/collections/<name>/remove-books`
```
[payload]
{
    "books":[]
}

[response]
200 OK

400 Bad Request
{"message":"this collection does not exist"}
{"message":"<isbn> is not a part of this collection"}
{"message":"no books provided"}
```

`HTTP DELETE /api/v1/collections/<name>/delete`
```
[response]
200 OK

400 Bad Request
{"message":"this collection does not exist"}
```

## Data Model

### Book
```
CREATE TABLE books (
    isbn uuid PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    description TEXT,
    metadata JSONB, 
    published_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)
```
### Collection
```
CREATE TABLE collection (
    id SERIAL PRIMARY KEY,
    name TEXT,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)

CREATE TABLE book_collections (
    book_isbn UUID,
    collection_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
)

CREATE UNIQUE INDEX book_collections_primary_idx ON book_collections (isbn, collection_id);
```

## Running the Server
`make docker-compose`

## Installing the CLI
`go install github.com/john-cai/book-manager/bm`