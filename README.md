# Golang Task

Getting Package

```ch
go get -u github.com/AndrewMislyuk/golang-task/imitation
```

Import

```go
import "github.com/AndrewMislyuk/golang-task/imitation"
```

Create a configuration file "imitation.yml" in the configs folder at the root of the project

Copy the text below into the generated config file

```go
sender_time: 2s

receiver_time: 5s

stop_time: 60s
```

Usage

```go
ctx, cancel := context.WithCancel(context.Background())
wg := sync.WaitGroup{}

imit, err := imitation.New()
if err != nil {
	log.Fatal(err)
}

go imit.StopAll(cancel)

wg.Add(3)
go imit.Sender(ctx, &wg, "Hello")

go imit.Sender(ctx, &wg, "Hi") // you can add second sender if you want

go imit.Receiver(ctx, &wg)

wg.Wait()
```
