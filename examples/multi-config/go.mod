module multi-config-example

go 1.24.5

require github.com/blackorder/configwatcher v0.0.0

require (
	github.com/blackorder/chanhub v0.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/blackorder/configwatcher => ../../
