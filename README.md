## Adda

Adda is a Multiuser chatting application, allowing users to communicate over tcp
It is some security features:

#### Limits the message rate

client has to wait 1sec before sending two consecutive messages otherwise he gets
a strike. 10 strikes and the client gets banned for 10 minutes. this is how we
deal with bots

#### Validates messages

It is quite reasonable that someone may try to break the server by sending invalid
text message. Here we use utf8 validation to check for invalid requests. Everytime
a request from IP sends a invalid string it will get a strike. As mentioned earlier
10 stikes will bann this user for 10 minutes

