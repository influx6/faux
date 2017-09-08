Mongo DB API
===============================

DB API attempts to provide a simple basic documentation which details the package expose.

## WithIndex
WithIndex is a function which the Mongo DB package expose which will set the provided index slice
into the configuration of the giving collection.

```go
WithIndex(ctx context.Context, col string, indexes ...mgo.Index) error) error
```

## Exec

Exec is a function which the Mongo DB package expose which will retrieved the needed
collection for the name provided and will call the provided function to perform its operation.

```go
Exec(ctx context.Context, col string, fx func(col *mgo.Collection) error) error
```


## Example

```go

var (
	events = metrics.New(stdout.Stdout{})

	config = mongo.Config{
		Mode:     mgo.Monotonic,
		DB:       os.Getenv("dap_MONGO_DB"),
		Host:     os.Getenv("dap_MONGO_HOST"),
		User:     os.Getenv("dap_MONGO_USER"),
		AuthDB:   os.Getenv("dap_MONGO_AUTHDB"),
		Password: os.Getenv("dap_MONGO_PASSWORD"),
	}

)

func main() {
	col := "ignitor_collection"

	ctx := context.New()
	api := mongo.New(testCol, events, mongo.New(config))

	elem, err := loadJSONFor(ignitorCreateJSON)
	if err != nil {
    panic(err)
	}

	err := api.Get(ctx, col, func(col *mgo.Collection) error {
    // Do something
  });

	if err != nil {
    panic(err)
	}

}
```
