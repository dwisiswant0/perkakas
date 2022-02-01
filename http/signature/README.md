# Signature

This module will help you creata HMAC signature based on sha256 & secret key.

## Usage

```go
package main

import "github.com/kitabisa/perkakas/signature"

func main() {
    message := "this is my message"
	secretKey := "123-qwe"

    // generate signature
    strSignature := signature.GenerateHmac(message, secretKey)
    
    // check is signature is match
    isSignatureMatch := signature.IsMatchHmac(message, strSignature, secretKey)
}
```

