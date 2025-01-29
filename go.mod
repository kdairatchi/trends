module github.com/kdairatchi/trends

go 1.20  # Ensure correct Go version

require (
    github.com/mmcdole/gofeed v1.1.0  # RSS Feed Parsing
    github.com/sirupsen/logrus v1.8.1  # Logging
    golang.org/x/net v0.8.0  # Networking
)

replace golang.org/x/net => golang.org/x/net v0.8.0
