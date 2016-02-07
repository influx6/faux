# Loop
 Loop provides a game loop for different rendering backends(browers RAF,gl,etc).
 It provides a central repository for a central working game loop regardless of
 specific platform/backend

# Install

  ```bash

    > go get -u github.com/influx6/faux/loop/...

  ```

# API
  The loop API is simple, because for all package there exists a central global
  function call

  - Loop(func(delta time.Time)) => Looper
     Assigns a function to be called on each loop, whilst returning a subscriber
    which stops the function assigned to the specific loop.

# Usage



  ```go

    import "github.com/influx6/faux/loop/gl"

    func main(){

        // Subscribe a function into the gameloop.
        gls := gl.Loop(func(delta float64){
          //.......
        })

        // End the function recall.
        gls.End()
    }

  ```
