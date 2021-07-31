# HLSQ 

A small CLI for adding some color to your HLS manifests and some basic filtering.
This CLI is not strict in its parsing so it will still work for manifests preceeded
by a grep. Named in tribute to the the great venerable [`jq`](https://github.com/stedolan/jq) cli.

<p align="center">
  <a href="https://github.com/soldiermoth/hlsq/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/soldiermoth/hlsq.svg?style=for-the-badge"></a>
  <a href="https://github.com/soldiermoth/hlsq/actions?workflow=Release"><img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/soldiermoth/hlsq/Release?style=for-the-badge"></a>
  <a href="https://goreportcard.com/report/github.com/soldiermoth/hlsq"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/soldiermoth/hlsq?style=for-the-badge"></a>
  <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
  <a href="https://github.com/goreleaser"><img alt="Powered By: GoReleaser" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge"></a>
  <br/>
  <a href="http://github.com/krzemienski/awesome-video"><img alt="Mentioned in Awesome" src="https://awesome.re/mentioned-badge-flat.svg"></a>
</p>

![Basic Example](images/basic.gif)

## Filtering

There are some basic filtering operations available in this CLI in the form of a single `{attribute name} {op} {value}`, this will be expanded in the future to accept more complex queries.

![Filtering Example](images/filter.gif)

Currently supported operations by value type
- Numbers: `>`, `>=`, `<`, `<=`, `=`, `!=`
- String: `=`, `!=`, `~`, `!~`, & `rlike`

## Install Instructions

### Pre-built Binary
Visit the [latest releases](https://github.com/soldiermoth/hlsq/releases) and pull a pre-built binary

### Homebrew

```
$ brew install soldiermoth/tap/hlsq
```

### From Source
Assuming a recent installation of Go is installed: [https://golang.org/doc/install](https://golang.org/doc/install)
```
$ go get github.com/soldiermoth/hlsq
```

## Demuxed Special Colors

As tribute to Demuxed2020 added colors matching the SWAG tshirts: `-demuxed`

![Demuxed Flag](images/demuxed2020.png)
