module github.com/lukasld/gocry

go 1.22.0

require (
	github.com/awnumar/memcall v0.2.0 // indirect
	github.com/awnumar/memguard v0.22.5 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	internal/cliSimpleCall v0.0.0 // indirect
)

replace internal/cliSimpleCall => ./internal/cliSimpleCall

require internal/cliCalls v0.0.0

replace internal/cliCalls => ./internal/cliCalls

replace internal/encryptionTools => ./internal/encryptionTools
