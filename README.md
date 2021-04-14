# Password Encoder

Password Encoder handles the hashing of given passwords.

## Motivation

To be able to create hashed passwords and to get the hashed passwords and statistics
about them. 

## Tech Used

- Golang 
- Gorilla Mux
- GoMock

## Tests
To run tests, enter the following terminal command:

```go test```

## API Reference

```POST /hash``` <br />
Creates a hash for a given form field password.

```GET /hash/{id}``` <br />
Returns a hash for a given id.

```GET /stats``` <br />
Returns total number of hashed passwords and the average request time to hash each.

```POST /shutdown``` <br />
Invokes graceful shutdown where server stops accepting requests and waits 
for all processes to finish before exiting.