# mendix-logfilter

This is based on https://github.com/42wim/rl.

Rate limit messages from stdin. (keep or drop those exceeding the limit)

# Use case

* When logging it can be convenient to limit the amount of data that can be sent.

# Building
Go 1.6+ is required. Make sure you have [Go](https://golang.org/doc/install) properly installed, including setting up your [GOPATH] (https://golang.org/doc/code.html#GOPATH
)

```
$ make
```

or build with a specific version of golang:

```
$ docker run --rm -it -v "$PWD":/build -w /build golang:1.9-stretch make
```

# Usage
```
Usage of ./mendix-logfilter:
  -k    keep the messages instead of dropping them
  -r int
        limit to r messages per second (drops those exceeding the limit) (default 5)


(the number of dropped messages will be sent to stderr, when not using the -k switch)
```


# Example
```
tail -f runtime.log | mendix-logfilter -k -r 5
```
