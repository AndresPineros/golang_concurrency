# golang_concurrency

Personal notes and exercises on Golang concurrency

### Sources

OReilly's Concurrency in Go

https://www.youtube.com/watch?v=5zXAHh5tJqQ&ab_channel=GopherAcademy

https://www.youtube.com/watch?v=hFqXgmor74k&ab_channel=GopherConIndia

https://golang.org/doc/effective_go#channels

https://www.youtube.com/watch?v=QDDwwePbDtw&ab_channel=GoogleDevelopers

https://www.youtube.com/watch?v=YEKjSzIwAdA&ab_channel=CodingTech


### General thoughts

- APIs should be synchronous. If we create a package that exposes functions, those functions should be synchronous. Go allows wrapping a syncrhonous API in goroutines, to make it asynchronous. DON'T USE CHANNELS IN PUBLIC APIS.

- Whenever I want to write to a channel I should ask "What happens if another goroutine is not reading from this channel and it blocks the current goroutine"? This is very important because sometimes we expect something to read from a channel but there's no guarantee. WRITING TO A CHANNEL IS A BLOCKING OPERATION.

- Closing channels is not mandatory but is a general best practice. Channels are cleaned up whenever the GC realizes nothing holds references to them. Thus, closing a channel doesn't really cleanup any resources, it just adds a termination message to the channel that can be read by all the readers.

- Closing channels is a tool to let readers know when the stream is over. It is important to close them because sthat's the way some goroutines are notified when they can stop working. If channels aren't properly closed, goroutines could be leaked.

- To end a goroutine use the context package, don't use done channels.

- Writer goroutines are responsible of creating, writing-to and closing channels. If there's a writing worker pool, closing can be synced with a WaitGroup.

- Be very careful when doing case ch1 <- <-ch2 in a select statement. If nothing is reading from ch1, the element that is read from ch2 will get lost and will never be added to ch1.

- Little's Law shows that buffered channels are not necessarily going to increase the performance of a pipeline. Buffered channels should be the last optimization done to a system because they tend to hide deadlocks and concurrency design issues.
