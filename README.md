### GOCRY - W.I.P.

![GoCry CLI - Schema](https://i.ibb.co/rdF9cxG/Gocry-V2-draw.jpg)

This is a work-in-progress repo for a Go-1Password-CLI
wrapper. 


Welcome to the Go-1Password-CLI wrapper repository! This project is currently a work-in-progress aimed at providing seamless integration for encyption using Go and 1Password CLI. With this wrapper, users will gain functionalities such as file encryption/decryption, automatic generation of 1Password entries and passwords, all securely stored in the GoCry (1PW) vault, and retrieval of passwords for encrypted files. This takes away the hussle of manually saving Keys, but instead saving them automatically in the GoCry vault.


It will allow the user:
- To encrypt/decrypt file(s) and folders
- automatically generate 1PW - entry and password, which is saved in the GoCry vault.
- retrieve passwords for encrypted files, write passwords to encrypt files.

Benefits:
- completely hands-off approach when encrypting files and handling passwords by leveraging 1PW - CLI.
- everything is handled through 1PW-CLI and the GoCry wrapper.

External Libraries used:
- memguard


TODO:
- Unit-tests.
- Memguard critical sections.
- Encryption package.
- CLI interface-package.
- Documentation.
