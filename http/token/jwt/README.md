# Package JWT

This package contains tools for JWT token creation and validation.
Supports JWT token creation and validation using symmetric and asymmetric key.

# JWT Creation Using Symmetric Key
```go
	signKey := []byte("abcde")

	claims := UserClaim{}
	claims.UserID = 12345
	claims.SecondaryID = "8fae85be-e441-4344-8634-d41f23684146"
	claims.Scopes = []string{"read"}
	claims.ClientID = "apdifuoqpweyr9823u"
	claims.Id = "63410cd1-110b-4a2c-8c3f-ae1535eda9a1"
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(48 * time.Hour).Unix()

	jwt := NewJWT(signKey)
    token, _ := jwt.Create(claims)
    fmt.Println(token)
```

# JWT Creation Using Asymmetric Key
```go
	claims := UserClaim{}
	claims.UserID = 12345
	claims.SecondaryID = "8fae85be-e441-4344-8634-d41f23684146"
	claims.Scopes = []string{"read"}
	claims.ClientID = "apdifuoqpweyr9823u"
	claims.Id = "63410cd1-110b-4a2c-8c3f-ae1535eda9a1"
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	pubKey, _ := ioutil.ReadFile("/path/to/app.rsa.pub")
	privKey, _ := ioutil.ReadFile("/path/to/app.rsa")

	jwt, _ := NewJWTRSA(pubKey, privKey)

    token, _ := jwt.Create(claims)
    fmt.Println(token)
```

For public and private key generation, you can use this `openssl` command:
```sh
# generate private key
openssl genrsa -out app.rsa 2048

# generate public key using private key
openssl rsa -in app.rsa -pubout > app.rsa.pub
```