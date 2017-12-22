# radolan
### Go package radolan parses the DWD RADOLAN / RADVOR radar composite format.
This format is used by the [DWD](http://www.dwd.de/DE/leistungen/radolan/radolan.html)
for weather radar data.

The obtained results can be processed and visualized with additional functions.
The example program `radolan2png` is included to quickly convert composite files to png images.

Besides local scans, the following grids are currently supported:
- National Grid (900km x 900km)
- National Picture Grid (920km x 920km)
- Extended National Grid (900km x 1100km)
- Middle-European Grid (1400km x 1500km)

Tested input products: 
| Product | Grid              | Description             |
| ------- | ----------------- | ----------------------- |
| EX      | middle-european   | reflectivity            |
| FX      | national          | nowcast reflectivity    |
| FZ      | national          | nowcast reflectivity    |
| PE      | local             | echo top                |
| PF      | local             | reflectivity            |
| PG      | national picture  | reflectivity            |
| PR      | local             | doppler radial velocity |
| PX      | local             | reflectivity            |
| PZ      | local             | 3D reflectivity CAPPI   | 
| RW      | national          | hourly accumulated      |
| RX      | national          | reflectivity            |
| SF      | national          | daily accumulated       |
| WX      | extended national | reflectivity            | 

Those can be considered working with sufficient accuracy.

### Documentation
Documentation is included in the corresponding source files and also available at
https://godoc.org/gitlab.cs.fau.de/since/radolan

### Installation
```
mkdir -p ~/go/src ~/go/pkg ~/go/bin
GOPATH="~/go" go get gitlab.cs.fau.de/since/radolan/radolan2png
```

### Sample image
This image shows radar reflectivity (dBZ) captured 31.07.2016 18:50 CEST
![alt text](https://gitlab.cs.fau.de/since/radolan/raw/master/assets/31-07-2016-1850.png)
Datenbasis: Deutscher Wetterdienst, Radardaten bildlich wiedergegeben
