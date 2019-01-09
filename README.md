# go-message-bridge

*This project is in Alpha but works as intended. *

Welcome to go message bridge. This is an app that allows you to connect Slack and Messages. It is a continuation of work done here in the following repos:

https://github.com/drbh/MessageBridge  
https://github.com/drbh/imessage-exporter  
https://github.com/drbh/imessage-anywhere  

and the following article:

https://medium.com/@david.richard.holtz/blue-green-texts-and-a-simple-solution-1c1981a00430

This is the most stable of all the mentioned projects and is not dependent on 3rd party packages for core functionality.

## Running on 
- Mac OSX Sierra 10.12.6 (16G1408)

## To Do
- [X] Rewrite Node to Golang
- [X] Handle double messages
- [X] Handle polling to be near realtime
- [X] Use a soild storage system for handle to channel resolution
- [X] Auto add bot to chat
- [X] Auto get the users id on startup
- [X] Removed prefix string on Slack->Messages messages
- [X] Resolve channel names via Address Book!!
- [ ] Test on OSX +10.12
- [ ] Setup data folder for db and config files
- [ ] Better channel namming for contacts
- [ ] Add channel clean up features (deleting)
- [ ] Add better UI (not console)


Slack setup instructions here   
https://docs.google.com/presentation/d/1YJmuXqQlD0wbIsd4XuyL3cFhSdeHqcsRGdoUgg36x4U/edit?usp=sharing

## To run

clone the repo. 
```
git clone https://github.com/drbh/go-message-bridge.git
```
cd into the directory. 
```
cd go-message-bridge
```
follow Slack Setup Instructions and 
```
{
    "bot_name": "your bots name",
    "bot_token": "your bot token from slack app",
    "user_token": "your slack legacy user token"
}
```
update config.json with Slack keys.  
run 
```
go run main.go
```

Now send someone a message in Messages and a channel should be made for that person.
