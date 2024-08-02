About Golang and things around...

# Building
- Mage - building tool - https://pkg.go.dev/github.com/magefile/mage#section-readme

# Code

- Milliseconds in log timeformat like `2023-11-20 23:15:10.123456 Hello` ... `log.SetFlags(log.LstdFlags | log.Lmicroseconds)`

# On Windows
- [Cygwin](https://www.cygwin.com/) - Linux environment
- [ConEmu](https://conemu.github.io/) - Terminal

> How to add `apt-cyg` package manager:  
> Require installation of Cygwin with packages `wget` or `lynx`.  
> `curl -O https://raw.githubusercontent.com/transcode-open/apt-cyg/master/apt-cyg`  
> `mv apt-cyg /usr/local/bin`  
> `apt-cyg update`


# External links
- [Vulkan bindings](https://github.com/goki/vulkan) - maintained fork of xlab's go-vulkan
- [Vulkan GPU framework](https://github.com/goki/vgpu) - like xlab's asche
- [Cogent Core](https://github.com/cogentcore/core) - Vulkan Go framework
- [Gio UI](https://gioui.org/) - multiplatform GUI for Go
- [MoltenVK](https://github.com/KhronosGroup/MoltenVK) - Vulkan on Apple devices
- [pkg.go.dev](https://pkg.go.dev) - Go package listing and search engine.
- [libs.garden](https://libs.garden/go) - Search engine and ranking watchdog for Go packages.
- [gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile) - A wrapper command around the gobind package's functionality.
- [golang/go/Mobile](https://github.com/golang/go/wiki/Mobile) - An article on the official Go repository's Wiki explaining various approaches to building and binding Go code for mobile app deployment.
- [Getting Started With WebAssembly](https://medium.com/swlh/getting-started-with-webassembly-and-go-by-building-an-image-to-ascii-converter-dea10bdf71f6) - Good article explaining the basics of building Go source to a WebAssembly binary and consuming it with HTML5.
- [Building shared libraries in Go: Part 1]() - Article showing how to create a C-style dynamic library with Go that is then consumed by a Python script.
- [GUI | Learn Go Programming](https://golangr.com/gui/) - Helpful list of Go packages that can be used to create GUIs in plain Go code.
