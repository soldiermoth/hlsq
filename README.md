# hlsq

A small CLI for adding some color to your HLS manifests and some basic filtering.
This CLI is not strict in its parsing so it will still work for manifests preceeded
by a grep.

![Basic Example](images/basic.gif)

## Filtering

There are some basic filtering operations available in this CLI in the form of a single `{attribute name} {op} {value}`, this will be expanded in the future to accept more complex queries.

![Filtering Example](images/filter.gif)

Currently supported operations by value type
- Numbers: `>`, `>=`, `<`, `<=`, `=`, `!=`
- String: `=`, `!=`, `~`, `!~`, & `rlike`

## Demuxed Special Colors

As tribute to Demuxed2020 added colors matching the SWAG tshirts: `-demuxed`

~[Demuxed Flag](images/demuxed.png)
