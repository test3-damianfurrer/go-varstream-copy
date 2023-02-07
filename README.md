# go-varstream-copy

you can send text with `nc -U echo.sock`

## install

`go install github.com/test3-damianfurrer/go-varstream-copy@latest`

### purpose
the intention is to use this for my streaming setup.
I can quickly cancel a video(+ audio) stream and replace it with a new instance, 
if I want to switch the audio.

At least when I am really live streaming my desktop.

I want to be able to stream a pre recorded video too tho.
I still want to be able to switch the audio without restarting the video.

for this I intend to stream the audio from a unix socket (out. unix)

The default socket should be there for the basic audio(e.g. the videos orig. audio)
It will be a fallback for when no stream is coming to the overlay socket.
It is mandatory and will only allow one, the initial connection.

The overlay socket is intened to replace the default stream, when it is available.
I intend to dynamicly switch audio to this socket.
There shall only be one connection active at a time.
To switch most elegantly(without processing the audio), 
the active connection should automaticly be closed and the new used, when a new connection is being built.


#### Since this is a very simplistic stream copy system, it could be used for all kinds of other applications tho.
It's basicly: backup system + replaceable live system.
If you wanted to connect your live web-page-server or similar behind a unix socket, 
you could drop in replace the old version just by starting up the new one
It could of course also just be switched to be ip based
