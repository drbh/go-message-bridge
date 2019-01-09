# go-message-bridge

This project is in Alpha but works as intended. 

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
- [] Test on OSX +10.12
- [] Setup data folder for db and config files
- [] Better channel namming for contacts
- [] Add channel clean up features (deleting)
- [] Add better UI (not console)


Slack setup instructions here   
https://docs.google.com/presentation/d/1YJmuXqQlD0wbIsd4XuyL3cFhSdeHqcsRGdoUgg36x4U/edit?usp=sharing

## To run

clone the repo. 
cd into the directory. 
follow Slack Setup Instructions and  
update config.json with Slack keys.  
run `go run main.go`

Now send someone a message in Messages and a channel should be made for that person. In order for the application to work - the Bot must be invited to the Slack channel. Do this by clicking "Invite person" in Slack. 

If you do not invite the Bot - no messages will show up!