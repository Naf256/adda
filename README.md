# Adda
Adda is a multiuser TCP-based chat application that allows multiple clients to communicate in real-time. It is designed with basic security features to protect the server and ensure smooth communication.

### Features
Multiuser Chatting: Connect multiple clients over TCP for real-time messaging.

Rate Limiting: Prevents message spamming by enforcing a delay between messages.

Security Against Bots: Implements a strike and ban system to discourage bots and abusive behavior.

Message Validation: Ensures all incoming messages are valid UTF-8 encoded strings to prevent server crashes or exploits.

### How It Works
 #### Message Rate Limiting
   - Clients must wait at least 1 second between two consecutive messages.

   - Violating this adds a strike.

   - 10 strikes = temporary ban for 10 minutes.

#### Message Validation
   - All messages are checked for UTF-8 validity.

   - Invalid messages result in a strike.

   - As above, 10 strikes from a single IP will result in a 10-minute ban.

## Getting Started
#### Requirements
- go 1.23 or higher

### Run Server
  ```bash
    go run main.go
  ```

### Run Client
  ```bash
    telnet 127.0.0.1 6969
  ```
  You can run as many client as you can with the same computer, telnet will assign a port for
  each one of them, pretending to be a different client

  - In case, you want to do dDos attack
    ```bash
    cat /dav/urandom | nc 127.0.0.1 6969
    ```
