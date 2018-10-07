# gotransition

A lightweight, object-oriented finite state machine implementation in Golang

> the trigger function execute processing

`
machine trigger -> event trigger -> transition execute -> state onEnter and onExit
     |                   |
     |                   |
     V                   |
 new machine             V
                     eventData
`
