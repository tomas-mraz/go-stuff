About Golang and things around...

# Docs
- [Mermaid](https://mermaid-js.github.io/mermaid/) - Graphs in markdown
- https://stackedit.io - Markdown online editor

# Code
- [Swagger/OpenAPI](https://swagger.io/) - API Designer 

# Building
- Mage - building tool - https://pkg.go.dev/github.com/magefile/mage#section-readme


## On Windows
- [Cygwin](https://www.cygwin.com/) - Linux environment
- [ConEmu](https://conemu.github.io/) - Terminal

> How to add `apt-cyg` package manager:  
> Require installation of Cygwin with packages `wget` or `lynx`.  
> `curl -O https://raw.githubusercontent.com/transcode-open/apt-cyg/master/apt-cyg`  
> `mv apt-cyg /usr/local/bin`  
> `apt-cyg update`


# External links

- [pkg.go.dev](https://pkg.go.dev) - Go package listing and search engine.
- [libs.garden](https://libs.garden/go) - Search engine and ranking watchdog for Go packages.
- [gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile) - A wrapper command around the gobind package's functionality.
- [golang/go/Mobile](https://github.com/golang/go/wiki/Mobile) - An article on the official Go repository's Wiki explaining various approaches to building and binding Go code for mobile app deployment.
- [Getting Started With WebAssembly](https://medium.com/swlh/getting-started-with-webassembly-and-go-by-building-an-image-to-ascii-converter-dea10bdf71f6) - Good article explaining the basics of building Go source to a WebAssembly binary and consuming it with HTML5.
- [Building shared libraries in Go: Part 1]() - Article showing how to create a C-style dynamic library with Go that is then consumed by a Python script.
- [GUI | Learn Go Programming](https://golangr.com/gui/) - Helpful list of Go packages that can be used to create GUIs in plain Go code.


### Grafana

https://github.com/lukasmalkmus/rpi_exporter

# Signing commits

I am using app called [SmartGit](https://www.syntevo.com/smartgit/) and [manual](https://docs.syntevo.com/SmartGit/HowTos/Sign-Tags-and-Commits.html) says
- if you are using Windows, please install [Gpg4win](https://gpg4win.org) (verified with Gpg4win 3.1.16)
- run gpa.exe (usually found in `C:\Program Files (x86)\Gpg4win\bin`) and create a key pair securing it with a passphrase
- in SmartGit Repository | Settings, tab Signing configure the full path to the gpg.exe and enter the Key ID of your created key pair
- if necessary, select Sign all commits
- when committing a file or tagging, a popup of GPG will occur and ask you for the key's passphrase (actually just the first time in the app's lifetime)

GitHub [documentation](https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key):

- run `gpg --full-generate-key` to generate key pair
- run `gpg --list-secret-keys --keyid-format=long` to get key ID 
- run `gpg --armor --export $USE_ID_YOU_GET` to get public key to copy into GitHub profile settings

