# golang_concurrency

Personal notes and exercises on Golang concurrency

### General thoughts

- Whenever I want to write to a channel I should ask "What happens other goroutine is not reading from this channel and it blocks the current goroutine"? This is very important because sometimes we expect something to read from a channel but there's no guarantee. WRITING TO A CHANNEL IS A BLOCKING OPERATION.

- Closing channels is not mandatory but is a general best practice. Channels are cleaned up whenever the GC realizes nothing holds references to them. Thus, closing a channel doesn't really cleanup any resources, it just adds a termination message to the channel.

- Closing channels is a tool to let know readers when the stream is over. It is important to close them because sthat's the way some goroutines are notified when they can stop working. If channels aren't properly closed, goroutines could be leaked.

- To end a goroutine use the context package, don't use done channels.

- Writer goroutines are responsible of creating, writing-to and closing channels. If there's a writing worker pool, closing can be synced with a WaitGroup.

- Be very careful when doing case ch1 <- <-ch2 in a select statement. If nothing is reading from ch1, the element that is read from ch2 will get lost and will never be added to ch1.

- Little's Law shows that buffered channels are not necessarily going to increase the performance of a pipeline. Buffered channels should be the last optimization done to a system because they tend to hide deadlocks and concurrency design issues.
