### GOCRY - W.I.P.

This is a work-in-progress repo for a Go-1Password-CLI
wrapper. 

It will allow the user:
- To encrypt/decrypt file(s)
- automatically generate 1PW - entry and password, which is saved in the GoCry vault.
- retrieve passwords for encrypted files

Benefits:
- completely hands-off approach when encrypting files and handling passwords.
- everything is handled through 1PW-CLI and GUI.

External Libraries used:
- memguard


TODO:
- Unit-tests.
- Memguard critical sections.
- Encryption package.
- CLI interface-package.
- Documentation.
