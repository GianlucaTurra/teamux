# TEAMUX

## Session Management

- [ ] ~~display open session real time using a color and/or an asterisk~~
- [x] prevent an already open session from opening (maybe attach to it)
- [x] enable attaching to an open session
- [x] prevent attaching to a not open session (maybe open and attach it)
- [x] kill a open session
- [x] update the legend
- [x] open and attach (needs further development with the db)
- [x] add confirm message for session delete

## General

- [x] manage properly the log
- [x] manage properly the db connection
- [ ] ~~create a log file in the .local directories~~
- [x] handle help text for each component
- [x] create a log file in the temp directory
- [ ] display a message box for errors and confirms
- [ ] during editing or creation the tree should not be displayed

## Session creation

- [ ] ~~create new session with the file field~~
- [x] create a input component
- [x] generate sql from input
- [x] handle key bindings in the vim way
- [x] allow to return to browser without creating a new session
- [x] clear input after pressing enter or returning to browser

## DB

- [x] implement a really lightweight orm with entities

## Errors

- [ ] display a simple message to the user to notify errors

## Bugs

- [ ] The first time sessions are loaded the blank pwd is not translated to $HOME
- [x] Creating a new session does not save the name nor the pwd
- [x] Open and not selected session has excessive padding
- [x] Sessions are loading with the wrong order (should be based on id asc)
- [ ] Editing a session creates a new session and does not update the existing one
- [ ] Up and Down keys do not trigger the updownmessage
