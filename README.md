# Go English Dictionary

simple English cli dictionary, that uses `en-word.net` database of english words

## Features

- Fast English word lookup
- Supports both exact matching and prefix matching
- Displays definitions and example usage
- Offline access after initial download

## Installation

```bash
go install github.com/kamildemocko/goendic@latest
```

## Usage

```bash
endic [OPTIONS] WORD
```

### Options

- `-e`: Use exact matching (default: prefix matching)
- `-l`: Return all results (default: limited to 10)
- `-d`: Debug mode

### Examples

```bash
# Search for words starting with "happy"
goendic happy

# Search for exact word "happy"
goendic -e happy

# Show all results for words starting with "happy" instead of showing just first 10
goendic -l happy
```

## How it works

1. Downloads English dictionary data from en-word.net
2. Stores data in a local SQLite database with FTS5 for fast text search
3. Provides command-line interface for word lookups

## License

[MIT](LICENSE)