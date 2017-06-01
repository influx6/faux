# Process
Process provide a simple structure for providing a nice means of executing shell scripts and commands easily.


## Example

- Executing simple commands synchronously
```go
src := process.SyncProcess{
    Commands: []process.Command{
        process.Command{
            Name: "echo",
            Args: []string{"New Login"},
        },
    },
}

ctx := context.Background()

var errBu, outBu bytes.Buffer
err := src.SyncExec(ctx, &outBu, &errBu)
```

- Executing simple commands asynchronously
```go
src := process.ASyncProcess{
    Commands: []process.Command{
        process.Command{
            Name: "echo",
            Args: []string{"New Login"},
        },
    },
}

ctx := context.Background()

var errBu, outBu bytes.Buffer
err := src.AsyncExec(ctx, &outBu, &errBu)
```

- Executing a shell script source

```go

src := process.ScriptProcess{
    Shell:  "/bin/bash",
    Source: `echo "New Login"`,
}

ctx := context.Background()

var errBu, outBu bytes.Buffer
err := src.Exec(ctx, &outBu, &errBu)
```