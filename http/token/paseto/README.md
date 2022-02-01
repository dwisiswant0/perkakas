# PASETO (Platform-Agnostic SEcurity TOkens)
Paseto is specification for implementation of stateless secure token, similar to JWT

## What Will Give You Motivation To Move From JWT to PASETO
JWT has critical security issue:
https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/

## About This Library
This paseto library:
* Will help encrypt and decrypt paseto token. Symmetric and asymmetric key are supported.
* Support only paseto token v2, since v1 token are deprecated

## Symmetric Or Asymmetric
Symmetric token are intended to create token for local usage; and asymmetric token are intended
to create token for public use. Following is local and public token example:

**Local paseto token**
```
v2.local.vx3Xd8jVUVJlsYwkQzEQ_XqjbuE3ejO49QnkBpL3imOGRoaKS-sXGZcJCfZekI1beLjzlYdRqnVKYDl_rC1Wp1rFbRym-a_HDM7V1Lk7xoksoxEwTOBEA4C7bX463odkc3Vge0rOudWsAOx8vQ9rV8LZtBGYqR4NVhPxdE4BdUytRXuUpHf09udPaL5gTmNGIpFFkg4MqAkqkWApyxla1Z-_mRWIe5h7wH_NE3l432s4aGnHk0iI3W3V7yVsMDRrmCDyndJI27jf1FIJ-6tRcn9KFA6s-48ZHVU4xzVjad0U6wnU61Ycz3ld3fNbz7_UjjP1bpBIzqmUAaCUPrvhaJuI0_CNqQWxdqxS1LBCUehfsebh-KfrPQUEpNbB7nOU-gAqYon8P9jyd4l9fIS6NJa4edXrv12cvkPx1BmGWzHbLtXoSX0ETA.S2l0YWJpc2EuY29t
```

**Public paseto token**
```
v2.public.eyJhdWQiOiJLaXRhYmlzYSBzZXJ2aWNlcyIsImVtYWlsIjoiYnVkaUBraXRhYmlzYS5jb20iLCJleHAiOiIyMDIxLTA4LTIxVDE2OjMyOjE3KzA3OjAwIiwiaWF0IjoiMjAyMS0wOC0yMFQxNjozMjoxNyswNzowMCIsImlzcyI6Ikt1bG9udXd1biIsImp0aSI6IjcwNmNiZmNlLWMwMzEtNGE0NC04MTVlLTAzMGY5NjNmN2Q0ZSIsIm5hbWUiOiJCdWRpIEFyaXlhbnRvIiwibmJmIjoiMjAyMS0wOC0yMFQxNjozMjoxNyswNzowMCIsInN1YiI6ImNhYzJlZTdlLTcwZDAtNDIyMC1iYWRkLTdiNTY5NWY1M2FkOCJ9gGMdg3eH8wFcIZ5tn3_87e7Z7jwfW-Q6vZqmkXijRTRGCFb171Vl4Rc8PRvXirvV6dj9Ns8xTBY0ksmRskz9DA.S2l0YWJpc2EuY29t
```

## Paseto Token Format
```<token_version>.<purpose>.<payload>.<optional_footer>```

* **token_version**: `v1` or `v2`. `v2` is recommended since `v1` is deprecated.
* **purpose**: `local` or `public`
* **payload**: token payload. Consist of token expiry time, audience, etc, and also your data
* **footer**: optional. Usually footer is identity of token issuer. I.e: company name.

## Example
### Symmetric Token Encryption
```go
// Symmetric key use 32 bytes characters long. Please create a secure one.
symmetricPaseto, err := NewSymmetric("PrU5AbXJawKJIUOJFmd4f6ZwmifLvvoF")
now := time.Now()

// Standard token data
token := paseto.JSONToken{
    Audience:   "Kitabisa services",
    Issuer:     "Kulonuwun",
    Jti:        "706cbfce-c031-4a44-815e-030f963f7d4e",
    Subject:    "cac2ee7e-70d0-4220-badd-7b5695f53ad8",
    Expiration: now.Add(24 * time.Hour),
    IssuedAt:   now,
    NotBefore:  now,
}

// Set your custom data here
token.Set("email", "budi@kitabisa.com")
token.Set("name", "Budi Ariyanto")

// Encrypt token with footer Kitabisa.com
encryptedToken, err = symmetricPaseto.Encrypt(token, "Kitabisa.com")
fmt.Println(encryptedToken)
```

### Symmetric Token Decryption
```go
symmetricPaseto, _ := NewSymmetric("PrU5AbXJawKJIUOJFmd4f6ZwmifLvvoF")

encToken := "v2.local.vx3Xd8jVUVJlsYwkQzEQ_XqjbuE3ejO49QnkBpL3imOGRoaKS-sXGZcJCfZekI1beLjzlYdRqnVKYDl_rC1Wp1rFbRym-a_HDM7V1Lk7xoksoxEwTOBEA4C7bX463odkc3Vge0rOudWsAOx8vQ9rV8LZtBGYqR4NVhPxdE4BdUytRXuUpHf09udPaL5gTmNGIpFFkg4MqAkqkWApyxla1Z-_mRWIe5h7wH_NE3l432s4aGnHk0iI3W3V7yVsMDRrmCDyndJI27jf1FIJ-6tRcn9KFA6s-48ZHVU4xzVjad0U6wnU61Ycz3ld3fNbz7_UjjP1bpBIzqmUAaCUPrvhaJuI0_CNqQWxdqxS1LBCUehfsebh-KfrPQUEpNbB7nOU-gAqYon8P9jyd4l9fIS6NJa4edXrv12cvkPx1BmGWzHbLtXoSX0ETA.S2l0YWJpc2EuY29t"
decToken, footer, err := symmetricPaseto.Decrypt(encToken)
if err != nil {
    return err
}

fmt.Printf("Decrypted token: %+v\n", decToken)
fmt.Println("Footer:", footer)
```

## Asymmetric Key Creation
```go
// You can create public and private keypair with ed25519 library from golang
// and encode it to hex string
import (
    "crypto/ed25519"
    "encoding/hex"
)

func main() {
    pubkey, privKey, _ := ed25519.GenerateKey(nil)
    pub := hex.EncodeToString(pubkey)
    priv := hex.EncodeToString(privKey)

    fmt.Println("Public key:", pub)
    fmt.Println("Private key:", priv)
}
```

### Asymmetric Token Ecryption
```go
pub := "07ddbb390ad3096b51138358c8e082686e8023853616768d0cd60cdb5ab68b58"
priv := "211eec1bbe7d6e4779542fb14591f8e2c0ef57607bff958b32267f1b9667bf0007ddbb390ad3096b51138358c8e082686e8023853616768d0cd60cdb5ab68b58"

ast, err := NewAsymmetric(pub, priv)
now := time.Now()

// Standard token data
token := paseto.JSONToken{
    Audience:   "Kitabisa services",
    Issuer:     "Kulonuwun",
    Jti:        "706cbfce-c031-4a44-815e-030f963f7d4e",
    Subject:    "cac2ee7e-70d0-4220-badd-7b5695f53ad8",
    Expiration: now.Add(24 * time.Hour),
    IssuedAt:   now,
    NotBefore:  now,
}

// Set your custom data here
token.Set("email", "budi@kitabisa.com")
token.Set("name", "Budi Ariyanto")

// Encrypt token with footer Kitabisa.com
encryptedToken, err := ast.Encrypt(token, "Kitabisa.com")
fmt.Println(encryptedToken)
```

### Asymmetric Token Deryption
```go
pub := "07ddbb390ad3096b51138358c8e082686e8023853616768d0cd60cdb5ab68b58"
priv := "211eec1bbe7d6e4779542fb14591f8e2c0ef57607bff958b32267f1b9667bf0007ddbb390ad3096b51138358c8e082686e8023853616768d0cd60cdb5ab68b58"

ast, _ := NewAsymmetric(pub, priv)

encToken := "v2.public.eyJhdWQiOiJLaXRhYmlzYSBzZXJ2aWNlcyIsImVtYWlsIjoiYnVkaUBraXRhYmlzYS5jb20iLCJleHAiOiIyMDIxLTA4LTIxVDE2OjMyOjE3KzA3OjAwIiwiaWF0IjoiMjAyMS0wOC0yMFQxNjozMjoxNyswNzowMCIsImlzcyI6Ikt1bG9udXd1biIsImp0aSI6IjcwNmNiZmNlLWMwMzEtNGE0NC04MTVlLTAzMGY5NjNmN2Q0ZSIsIm5hbWUiOiJCdWRpIEFyaXlhbnRvIiwibmJmIjoiMjAyMS0wOC0yMFQxNjozMjoxNyswNzowMCIsInN1YiI6ImNhYzJlZTdlLTcwZDAtNDIyMC1iYWRkLTdiNTY5NWY1M2FkOCJ9gGMdg3eH8wFcIZ5tn3_87e7Z7jwfW-Q6vZqmkXijRTRGCFb171Vl4Rc8PRvXirvV6dj9Ns8xTBY0ksmRskz9DA.S2l0YWJpc2EuY29t"
decToken, footer, err := ast.Decrypt(encToken)
if err != nil {
    return err
}

fmt.Printf("Decrypted token: %+v\n", decToken)
fmt.Println("Footer:", footer)
```