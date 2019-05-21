https://travis-ci.org/DUE-Info-AG-APP/DuckTalk.svg?branch=master
# Description
#### A simple messaging service.

# API
## Models
### Message model:
    id: [PK] integer
    sender: text
    recipient: text
    timestamp: timestamp with timezone
    data: text

### User model:
    username: text
    password: text
    contacts: character varying[] (255)


## Possible Requests

#### Get list of all users on the server
    GET /users

#### Create new user
    POST /users
    {
      "username": "user_x",
      "password": "1234"
    }

#### Get data of specific user
    GET /users/{username}
Required Headers:

    Authorization: BASIC username:password in BASE64


#### Add contact to user
    POST /users/{username}/contacts
    {
      "username": "nameofcontact"
    }


#### Send message
    POST /messages
    {
      "sender": "user_x",
      "recipient": "user_y",
      "timestamp": "2019-05-07T16:01:45+02:00",
      "data": "bla bla bla"
    }


#### Get messages
relating to specified sender and recipient
where authorization has to be provided with either the
senders or the recipients credentials

    GET /messages?sender="senderusername"&recipient="recipientusername"
Required Headers:

    Authorization: BASIC username:password in BASE64
