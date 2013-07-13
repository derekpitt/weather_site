# Weather site

Just the code for my personal weather site using my other weather packages for go.

This site includes the server and the pusher (the station where it gets the data from the console and posts the data to the server)

The pusher and server share a key that they use to verify identity. Don't want someone else inserting their own data!

## Server

Dead simple searver that displays and inserts new data.

### TODO:
 - Cache the results in a goroutine
 - Finish trend charts
 - Finish History charts
 - Add more data to view

## Pusher

Polls the console, signs the data and sends it off to the server. It will hold in memory any resends (just incase a post failed)

### TODO:
 - Also save resends to files


That's it for now.. Will soon be up at http://weatheratdereks.com


