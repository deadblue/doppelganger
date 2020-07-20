# Doppelganger

A task dispatcher written in Golang.

## Introduction

Doppelganger runs tasks in sepcified frequence, and can retry on failed.

Task should be one of following types:

* Command
* HTTP

## Usage

Doppelganger supports two entrypoints:

* HTTP: JSON & Protocol Buffers 
* Unix Domain Socket: Protocol Buffers

## License

MIT
