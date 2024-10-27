# Gotots

Gotots is a tool for generating TypeScript types from Go definitions.


## Installation

### Installing the tool
To install the Gotots, clone the repository and run `go install` inside it.
This will install Gotots to $GOPATH/bin. Make sure that $GOPATH/bin is added to your path in order to be able to run it.
On Mac and Linux this can be done by adding `export PATH="$(go env GOPATH)/bin:$PATH"` to your `.bashrc` or `.zshrc`.

### Adding it to your project
The first step is to add the following line to the go files containing the types you want to be generated.
```go
//go:generate gotots
```
This enables the use of the `go generate` command to generate the types.
Before this works you also need a config file to configure things like where the generated files should go. 
The config file is called `gotots.yaml` and can be placed anywhere in your go project.

The config file should look something like this:
```yaml
out-path: "relative/path/to/output/folder" # A path relative to the config file where the generated files should be put
use-comments: false # Whether the generated types should also use the comments from the go type
```
You can have multiple config files and each file will use the first one found.

Now you should be able to run `go generate ./...` and the files should all be generated.
