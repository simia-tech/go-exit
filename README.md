# An exit strategy for go routines.

[![Build Status](https://travis-ci.org/simia-tech/go-exit.svg)](https://travis-ci.org/simia-tech/go-exit)
[![Code Coverage](http://gocover.io/_badge/github.com/simia-tech/go-exit)](http://gocover.io/github.com/simia-tech/go-exit)
[![Documentation](https://godoc.org/github.com/simia-tech/go-exit?status.svg)](https://godoc.org/github.com/simia-tech/go-exit)

The library helps to end the go routines in your program and collects potential errors.

## Install

`go get github.com/simia-tech/go-exit`

## Example

```go
func main() {
	exit.Main.SetTimeout(2 * time.Second)

	counterExitSignal := exit.Main.NewSignal("counter")
	go func() {
		counter := 0

		var reply exit.Reply
		for reply == nil {
			select {
			case reply = <-counterExitSignal.Chan:
				break
			case <-time.After(1 * time.Second):
				counter++
				fmt.Printf("%d ", counter)
			}
		}

		switch {
		case counter%5 == 0:
			// Don't send a return via reply to simulate
			// an infinite running go routine. The timeout
			// should be hit in this case.
		case counter%2 == 1:
			reply.Err(fmt.Errorf("exit on the odd counter %d", counter))
		default:
			reply.Ok()
		}
	}()

	if report := exit.Main.ExitOn(syscall.SIGINT); report != nil {
		fmt.Println()
		report.WriteTo(os.Stderr)
		os.Exit(-1)
	}
	fmt.Println()
}
```

## License

The project is licensed under [Apache 2.0](http://www.apache.org/licenses).
