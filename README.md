# RG Downloader

* [Download (Windows)](https://github.com/AegirLeet/rg-downloader/releases/latest/download/rg-downloader.exe)
* [Download (Linux)](https://github.com/AegirLeet/rg-downloader/releases/latest/download/rg-downloader)

## Usage

Download and run the binary for your platform, select a download directory.  
Wait for the download to finish.  
<kbd>Ctrl-c</kbd> to interrupt.

## Building the binaries

### Prerequisites

* Go 1.24
* `make`
* GTK 3 dev packages (`libgtk-3-dev` or similar)

### Build

```shell
make deps          # install dependencies
make build-linux   # build linux binary
make build-windows # build windows binary
```
