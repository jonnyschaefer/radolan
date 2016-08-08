# radolan
### Go package radolan parses the DWD RADOLAN / RADVOR radar composite format.
This format is used by the [DWD](http://www.dwd.de/DE/leistungen/radolan/radolan.html)
for weather radar data.

The obtained results can be processed and visualized with additional functions.
The example program `radolan2png` is included to quickly convert composite files to png images.

Currently the national and the extended european grids are supported.
Tested input products are PG, FZ, SF, RW, RX and EX. Those can be considered working with
sufficient accuracy.

### Documentation
Documentation is included in the corresponding source files and also available at
https://godoc.org/gitlab.cs.fau.de/since/radolan

### Installation
```
mkdir -p ~/go/src ~/go/pkg ~/go/bin
GOPATH="~/go" go get gitlab.cs.fau.de/since/radolan/radolan2png
```

### Sample image
This image shows the radar reflectivity (dBZ) captured 31.07.2016 18:50 CEST
![alt text](https://gitlab.cs.fau.de/since/radolan/raw/master/assets/31-07-2016-1850.png)
```Datenbasis: Deutscher Wetterdienst, Radardaten bildlich wiedergegeben```
