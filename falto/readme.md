# Falto
  Fault tolerant low-latency log only readers and writers that provide a fast,
  concurrent read and writes to a log based WAL system.

## Install

    ```bash

        go get -u github.com/influx6/falto

    ```

## API
  API is still under heavy development

## Example  


  ```go

    // Admin provides a structure which gets logged.
    type Admin struct{
      Name string
      Level int
    }

    redis := falto.NewRedisBackend()
    zipper := falto.New(redis)


    // get a writer from the logdb.
    adminw := zipper.Writer()
    adminw.Write(Admin{Name:"John Wax",Level: 300})

    // get a reader from the logdb
    adminr := zipper.Reader()

    data := adminr.Read() // returns a FaltoRecord{ data:[]byte{....}}
    adminw.Write(Admin{Name:"John Wax",Level: 500})

    //get a reader from the logdb at a specific position using timestamps.
    adminr2 := admin2.Reader(Falto.TimeSeek{ timestamp: (2 * time.Hour), Direction: -1 })

    //get a reader from the logdb at a specific position using timestamps.
    adminr3 := admin2.Reader(Falto.RangeSeek{ Range: 30, Direction: -1, Position: 0 })

  ```