# serverChat
Simple golang chat server

# TCP chat with tls protocol
____
## Running

**Сertificate:**  
go run .\tls-self-signed-cert.go

**Server:**  
go run .\server.go

**Client:**  
go run .\client.go --name=sasha

## Using chat

**To send a message to a client, you need to know his name:**  

**pasha**: ___sasha hello___

**To send a message to all clients:**  
___all hello___
