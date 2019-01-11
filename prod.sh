rm -rf '/Users/davidholtz/Desktop/go-message-bridge/dend/MessageBridge.app'
rm -rf '/Users/davidholtz/Desktop/go-message-bridge/dist/app'
go build -o app main.go && mv app dist/app && sh macapp.sh