# Bot for Block Battle on The AI Games

This is my final project for my COMPSCI 383 AI class at UMass.

It contains both the logic for training a bot using reinforcement learning
and for creating a bot that can be played in the competition.

# Structure

* `./game`: Core game logic, about the board and moves.
* `./player`: Translates between the text game engine instruction
  and the `game` structures.
* `./bots`: Creates a `player.Player` that interacts with stdout and stdin
   and defines a `bots.Bot` that can execute moves on this player.
* `./train`: Runs the local java game engine to simulate games and training
  a model on them.

## `./bots`
Creating a Go bot for The AI Games isn't straightforward because they require
all go files to be in one (a `main`) package. This prevents any sort of
imports of external files or code organization.

So I got a bit creative to enable us to use the nice module structure of Go.
The `./bots/process` is basically the `main` module that we want the site to
run. If you run `go run ./bots/process/*.go` it will start up a little bot
that responds to `stdin` and outputs actions to `stdout`. But we can't just
upload this directory, because the server won't be able to import the required
things. So instead, we build it into a binary that will run on their server
(64 bit linux) and create a little C program in `./bots/binary/main.c` that will
`exec` that binary when it runs.

All we have to do is run `make bot.zip` and it will build the Go binary
and zip that up with the C code.

## `./train/engine`
In this directory we have some self contained code to create two `player.Player`s
by running a java engine in another process and getting back the two players
for this.

The Java engine takes in two commands and runs those as the two player bots.
So in order to communicate from our Go code to this process I set up the
following:

1. For each player/bot we create two temporary files on disk, `input` and `output`.
   We will write our actions to the `output` file and continuously read the
   `input` file to monitor for changes the engine sends there.
2. We start up process running the Java engine and tell it to use the
   `./train/engine/pipe.bash <output file> <input file>` command for each player.
   This script is very simple. It takes its `stdin` and pipes that to the
   `input` file and reads the `input` file to stdout. While this might seem
   backward (why are we writing to input?) the terms are relative. The `input`
   file is input to our bot and output from the engine. So when the engine
   writes to this process (like `update game ...`) it will be written to the
   `input` file and our go process will pick it up.
3. When the Java engine stops, we delete the temporary files. This cancels
    the input and output channels in Go and we know that the game has ended.



## Methods

Reinforcement learning

Policy: State -> Action
