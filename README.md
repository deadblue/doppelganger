# Doppelganger

A task dispatcher written in Golang.

## Introduction

**W.I.P.**

Doppelganger runs tasks in sepcified frequence, and retries when task failed.

Task should be one of following types:

* Command
* HTTP

After a task executed, the result can be sent to a callback endpoint.

## License

MIT
