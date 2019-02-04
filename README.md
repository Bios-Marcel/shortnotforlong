# Short but not for long

This is a small link shortener written in golang. It's not meant for
permanently shortening link. So in case you don't want dead links after a
reboot, this is not for you.

On top of forgetting everything on reboot, it can't hold many links and the
speed of generating new links will probably get slower after some thousand
shortened urls.

## Usage example

```go
func main() {
	shortener := NewShortener(1234)
	fmt.Println(shortener.Shorten("https://google.com"))

	blocker := make(chan struct{})
	go func() {
		log.Fatalln(shortener.Start())
		blocker <- struct{}{}
	}()

	<-blocker
}
```

Run it, open your browser and visit the link that the main-function spits out.