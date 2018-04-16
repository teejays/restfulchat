# RestfulChat API

A RESTful chat API based on Go's net/http package, httprouter, and my GoFileDb package.

## Getting Started
Setting this API up should be fairly easy to set up if you have Go set up. 
### Prerequisites:
1) Install GoLang from the official [Go website](https://golang.org/).
2) Install the following Go packages:
	* [httprouter](https://github.com/julienschmidt/httprouter): ``` go get github.com/julienschmidt/httprouter```
	* [gofiledb](https://github.com/teejays/gofiledb): ``` go get github.com/teejays/gofiledb```
    
    

### Installation:
1) Clone this repository: 

	```git clone https://github.com/teejays/restfulchat.git```
    
2) In the project folder, edit the ```GoFiledbRoot``` variable in _settings.json_ file to a folder path that you have read & write access to (this is where all the data will be stored).

3) Compile the application: 

	```go build -o server.out``` or ```go install``` or however you feel comfortable.
4) Assuming that the compiled executable is named _server.out_, start the server:
	
    ```./server.out```  

### Testing:

The unit tests for this package are located in _./tests_ folder. To run the tests, use the following command from the project root directory:

```go tests -v ./tests```

---

## Documentation
### API Endpoints:

We have the following API endpoints
* **GET /v1/chat/:userid:** Fetches the chat log of the _userid_	 	
	* CURL e.g. ```curl localhost:8080/v1/chat/someuser1```


* **POST /v1/chat/:userid:** Sends a message from _userid_. The message content and recipient is provided in the request body. 
	* CURL e.g. ```curl localhost:8080/v1/chat/someuser1 -X POST -H "Content-Type: application/json" -d '{"Content":"Hello World!", "To":"someuser2"}'```


* **PUT /v1/chat/:userid:** Edits a message previously sent from _userid_ to a given recipient. The id of the message to edit, the recipient, and the new message content are provided in the request body. 
	* CURL e.g. ```curl localhost:8080/v1/chat/someuser1 -X PUT -H "Content-Type: application/json" -d '{"MessageId": 1, Content":"Hello World! (edited)", "To":"someuser2"}'```

* **DELETE /v1/chat/:userid:** Deletes a message previously sent from _userid_ to a given recipient. The id of the message to delete and the recipientare provided in the request body. 
	* CURL e.g. ```curl localhost:8080/v1/chat/someuser1 -X DELETE -H "Content-Type: application/json" -d '{"MessageId": 1, "To":"someuser2"}'```



---
## Notes
### Data Structures
The applications is based on three objects:
1) _User_: Represents a user.
    * Structure:
    	* _UserId_ (string)
    * _Buddy_: a user that another user is interacts with.


2) _Conversation_: A conversation is stored communication between two or more users.
	* Structure: 
		* _UserIds_: an array of user ids of all the users that are a part of a conversation
		* _Messages_: An array of _Message_
		* _LastMessageId_ (int): Keeps track of the last (also largest) unique message id so the new messages can be given an appropriate id.


3) _Message_: The most basic data unit that makes a conversation.
	* Structure:
		* _Id_ (int): unique identifier of a message within a conversation
		* Content (string): the content of a message
		* TimestampCreated (time): when the message was sent
		* TimestampUpdated (time): when the message was last updated
		* From (string): contributed the message in a conversation

### Authentication
The authentication layer for this server hasn't been implemented yet. However, the API is built in a way that that Basic Auth could be incorporated easily without changing the structure of the code.

### Database
I am using my own [_GoFiledb_](https://github.com/teejays/gofiledb) package for as a database. GoFiledb is a simple, minimalistic Go client that lets applications use the filesystem as a database. The main advantage of GoFiledb is that it uses the years of optimization efforts that went into file systems to make reading and serving of data is very fast. It is very quick to set up (vs. a proper database, which are sometimes an overkill for a simple project). 

_Scalability:_
This is a minimalistic API, developed mostly for fun and experimentation reasons. In order to scale it further, a few decisions probably need to be changed. For example, the local file syetem based data storage should probably be replaced by a proper schemaless DB system.

### Contact Info
For any issues or feedback, please create an issue in Github for this repo.
