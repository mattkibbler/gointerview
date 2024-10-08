# gointerview

A small program to rehearse potential interview questions and practice using Go.

![Demo GIF](https://raw.githubusercontent.com/mattkibbler/gointerview/refs/heads/main/demo/demo.gif)

## What?

A command line program which prompts the user with random pre-loaded questions and records their answer. After answering the user is shown the correct answer and asked to mark their own work - did they get it right or wrong. In the future I would like to integrate ChatGPT or some other AI to have this done automatically.

## Usage

| Command          | Description                      |
|------------------|----------------------------------|
| `make build`     | Builds the binary.               |
| `make run`       | Builds and runs the binary.      |

### Flags

Pass flags like to the run command like so: `make run FLAGS="..."`

#### Available flags

| Flag         | Description                             |
|--------------|-----------------------------------------|
| `db`         | The path to the sqlite DB (optional)    |




## Roadmap

With the foundation now in place I will look at tackle the following...

- [x] Add a database to store questions as well as record the user's answers.
- [x] Categories to request questions of a particular type.
- [ ] Make it look nicer... add fancy colours?
- [ ] Integrate AI to verify answers and maybe generate new ones on the fly.
