package main

type MessageOptions interface {
    sendMessage(host string, port string)
}
