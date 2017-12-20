package radolan

type spec struct {
	px int // plain data dimensions
	py int

	dx int // data (layer) dimensions
	dy int

	rx float64 // resolution
	ry float64
}

// local picture products do not provide dimensions in header
var catalog = map[string]spec{
	"OL": {200, 224, 200, 200, 2, 2},  // reflectivity (no clutter detection)
	"OX": {200, 224, 200, 200, 1, 1},  // reflectivity (no clutter detection)
	"PD": {200, 224, 200, 200, 1, 1},  // radial velocity
	"PE": {200, 224, 200, 200, 2, 2},  // echotop
	"PF": {200, 224, 200, 200, 1, 1},  // reflectivity (15 classes)
	"PH": {200, 224, 200, 200, 1, 1},  // accumulated rainfall
	"PL": {200, 224, 200, 200, 2, 2},  // reflectivity
	"PM": {200, 224, 200, 200, 2, 2},  // max. reflectivity
	"PR": {200, 224, 200, 200, 1, 1},  // radial velocity
	"PU": {200, 2400, 200, 200, 1, 1}, // 3D radial velocity
	"PV": {200, 224, 200, 200, 1, 1},  // radial velocity
	"PX": {200, 224, 200, 200, 1, 1},  // reflectivity (6 classes)
	"PY": {200, 224, 200, 200, 1, 1},  // accumulated rainfall
	"PZ": {200, 2400, 200, 200, 2, 2}, // 3D reflectivity CAPPI
}
