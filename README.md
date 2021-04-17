# GoPacktpub-Downloader
[![Run Gosec](https://github.com/ne0z/GoPacktpub-Downloader/actions/workflows/security-scan.yaml/badge.svg)](https://github.com/ne0z/GoPacktpub-Downloader/actions/workflows/security-scan.yaml)

This is tool to generate Ebook (amazon kindle friendly) from the Packtpub subscription library. So, you need to subscribe Packtpub library first to execute this command.

## Dependencies
* Calibre `brew install calibre`

## How to install
```bash
$ git clone git@github.com:ne0z/GoPacktpub-Downloader.git
$ cd GoPacktpub-Downloader
$ go get .
$ make build
```

## How to use
### Login
```bash
$ ./packt login
Username : your@email.com
Password : *******
```

### Search book
```bash
$ ./packt search kubernetes
```

### Generate Ebook
```bash
$ ./packt <extension> <isbn>
$ ./packt epub 9781800207974
$ ./packt mobi 9781800207974