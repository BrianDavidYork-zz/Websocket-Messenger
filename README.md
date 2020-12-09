# Websocket-Messenger

Websocket-Messenger is the back-end api for an instant messager/chat feature.  This was an educational project that allowed me to build a small API in Golang and learn about websockets.  Because this is an educational project, a nummber of things could be added to make it more production ready -- the option to send the full message through the websocket, pagination to the messages returned for each conversation, adding multiple users for each conversation as opposed to one-on conversations implemented here, optomizing the HTTP response by creating a custom response struct, conversations could be generated when a message first message is sent between two users not already in a conversation,  etc.

# API Endpoints

The following are all the endpoints in the API:

### POST - /user --- create a new user
### GET - /user/{username} --- fetch a user's profile
### POST - /user/login --- user login
### POST - /user/logout --- user logout

### POST - /conversation --- create a new conversation
### GET - /conversation --- get all conversation user is a member of

### POST - /message --- create a new message
### PUT - /message --- edit a message
### DELETE - /message{MessageId} --- delete a message
### GET - /message/{ConversationId} --- get all messages in a conversation

### GET - /websocket --- opens a websocket with the server


# How To Use The API

All endpoints return JSON with 2 fields:
  Message - a message telling you about success or error
  Data - any requested data
  
All endpoints (excpet user create, login, and websocket) are authenticated through JSON web tokens.
All requests must contain a request Header in the following format:
  Authorization: Bearer {Jwt Token}

The following is an explanation of each endpoint:

### POST - /user --- create a new user

Endpoint accepts a JSON body containing 2 fields:
  Username
  Password
  
Creating a new user returns a JWT in the Data field of the response JSON.  The user is logged in on creation.  This JWT will remain active until user is logged out.
  
### GET - /user/{username} --- fetch a user's profile

Returns a user's profile containing the following fields:
  Username
  LoggedOn
  LastOnline
  Created

### POST - /user/login --- user login

Endpoint accepts a JSON body containing 2 fields:
  Username
  Password
  
### POST - /user/logout --- user logout

No informtion is needed in the request besided the Authorization header.

### POST - /conversation --- create a new conversation

Endpoint accepts a JSON body containing 1 field:
  Recipient  -- this is a username
  
Returns the conversationId of the newly created conversation.
  
### GET - /conversation --- get all conversation user is a member of

No informtion is needed in the request besided the Authorization header.

### POST - /message --- create a new message

Endpoint accepts a JSON body containing 2 fields:
  ConversationId
  Message
  
Returns the MessageId of the newly created message.

### PUT - /message --- edit a message

Endpoint accepts a JSON body containing 2 fields:
  MessageId
  Message
  
### DELETE - /message{MessageId} --- delete a message

No informtion is needed in the request besided the Authorization header.

### GET - /message/{ConversationId} --- get all messages in a conversation

No informtion is needed in the request besided the Authorization header.

### GET - /websocket --- opens a websocket with the server

This cannot be sent with an http request but with a ws request.  If the API is running n localhost port 8000, it would be sent like this:
  ws://localhost:8000/websocket
  
When running locally, the API does not accept secure websockets (wss://) only unsecured (ws://).

# Websocket Use

After the websocket is opened you will receive a JSON response:
  {"Message": "Websocket Open"}
  
After the websocket is opened, the user must authenticate by sending the JWT in the websocket.

After authentication you will receive a JSON response:
  {"Message": "Websocket Authenticated"}
  
If the websocket is not authenticated, the websocket will immediately close.

This api supports showing users a "{username} is typing..." message when the other user of a conversation is writing a message.
The typing user sends the following JSON through the websocket:
  {
    "User": "{typing user's username}",
    "Message": "Typing",
    "ConvId": "{ConversationId}"
  }
  
This same object is then sent to the other member of the conversation through the websocket.
All of this would normally be handles by the front-end head, and based on the above object,
a message could be shown letting user1 of a conversation that user2 is typing on a second to second basis.



## Created By BRIAN YORK
