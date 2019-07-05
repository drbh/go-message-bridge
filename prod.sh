rm -rf 'dend/MessageBridge.app'
rm -rf 'dist/app'
rm -rf 'dist/minimal'

go build -o dist/minimal viewer/window.go

go build -o app main.go && mv app dist/app && sh macapp.sh