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
- [ ] ~~create new session with the file field~~
- [x] create a input component
- [x] generate sql from input
- [x] handle key bindings in the vim way
- [x] allow to return to browser without creating a new session
- [x] clear input after pressing enter or returning to browser
- [x] allow to show related windows
- [x] allow to manage related windows

## General

- [x] manage properly the log
- [x] manage properly the db connection
- [ ] ~~create a log file in the .local directories~~
- [x] handle help text for each component
- [x] create a log file in the temp directory
- [x] display a message box for errors and confirms
- [ ] ~~during editing or creation the tree should not be displayed~~

## Window Management

- [x] Enable editing
- [x] Enable deleting

## Pane Management

- [x] Enable basic crud operations

## DB

- [x] implement a really lightweight orm with entities

## Errors

- [x] display a simple message to the user to notify errors

## Help

- [x] Help text should have his own separate component
- [x] When changing tab the help text should change too
- [x] After changing to windows you have to press `?` again to see the help text
- [x] In the editor the help text should be shown by default

## Bugs

- [x] The first time sessions are loaded the blank pwd is not translated to $HOME
- [x] Creating a new session does not save the name nor the pwd
- [x] Open and not selected session has excessive padding
- [x] Sessions are loading with the wrong order (should be based on id asc)
- [x] Editing a session creates a new session and does not update the existing one
- [x] Up and Down keys do not trigger the `UpDownMsg`
- [x] Extending the help text pushes the detail to the right

## Refactoring

- [x] Editor don't cycle properly at the end or start
- [x] Use interfaces and changing elements in containers
