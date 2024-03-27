module github.com/lukasld/gocry

go 1.22.0

require internal/cliSimpleCall v0.0.0 // indirect
replace internal/cliSimpleCall => ./internal/cliSimpleCall

require internal/cliCalls v0.0.0
replace internal/cliCalls => ./internal/cliCalls

replace internal/encryptionTools => ./internal/encryptionTools
