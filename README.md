# goPasetoV4

The `goPasetoV4` package provides a robust and secure implementation for creating and validating Paseto (Platform-Agnostic Security Tokens) tokens in Go. This package focuses on ease of use, security, and flexibility, making it suitable for enterprise-level applications that require reliable authentication mechanisms.

## Features

- Secure token generation using Paseto V4.
- Flexible nonce handling with environmental configuration or automatic generation.
- Comprehensive validation of tokens, including expiration and integrity checks.
- Easy integration with Go applications.

## Installation

To install `goPasetoV4`, use the `go get` command:

```bash
go get -u github.com/anfen93/goPasetoV4
```

This command retrieves the library and installs it.

## Usage

Below is a quick example of how to use `goPasetoV4`:

### Token Creation

```go
package main

import (
    "github.com/anfen93/goPasetoV4"
    "log"
    "time"
)

func main() {
    maker := goPasetoV4.NewPasetoMaker()

    // Create a token for a specific username with a 1-hour duration
    token, err := maker.CreateToken("username", time.Hour)
    if err != nil {
        log.Fatalf("Error creating token: %v", err)
    }

    log.Printf("Generated Token: %s", token)
}
```

### Token Verification

```go
package main

import (
    "github.com/anfen93/goPasetoV4"
    "log"
)

func main() {
    maker := goPasetoV4.NewPasetoMaker()

    token := "your_generated_token"
    payload, err := maker.VerifyToken(token)
    if err != nil {
        log.Fatalf("Error verifying token: %v", err)
    }

    log.Printf("Token Payload: %+v", payload)
}
```

## Configuration

You can configure the nonce used in token generation by setting the `PASETO_NONCE` environment variable. If not set, the system will generate a random nonce for each instance.

## Contributions

Contributions are welcome! Please feel free to submit a pull request or open issues to improve the functionality, documentation, or code quality of `goPasetoV4`.

## TODOS
Right now the implementation of the Payload is not so flexible. It should be possibile to personalize the payload with custom fields.
## License

This project is licensed under the [MIT License](LICENSE.md).
