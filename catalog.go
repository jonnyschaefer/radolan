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
var dimensionCatalog = map[string]spec{
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

type Unit int

const (
	Unit_unknown = iota
	Unit_mm      // mm/interval
	Unit_dBZ     // dBZ
	Unit_km      // km
	Unit_mps     // m/s
)

func (u Unit) String() string {
	return []string{"unknown unit", "mm", "dBZ", "km", "m/s"}[u]
}

var unitCatalog = map[string]Unit{
	"CH": Unit_mm,
	"CX": Unit_dBZ,
	"D2": Unit_mm,
	"D3": Unit_mm,
	"EA": Unit_dBZ,
	"EB": Unit_mm,
	"EC": Unit_mm,
	"EH": Unit_mm,
	"EM": Unit_mm,
	"EW": Unit_mm,
	"EX": Unit_dBZ,
	"EY": Unit_mm,
	"EZ": Unit_mm,
	"FX": Unit_dBZ,
	"FZ": Unit_dBZ,
	"HX": Unit_dBZ,
	"OL": Unit_dBZ,
	"OX": Unit_dBZ,
	"PA": Unit_dBZ,
	"PC": Unit_dBZ,
	"PD": Unit_mps,
	"PE": Unit_km,
	"PF": Unit_dBZ,
	"PG": Unit_dBZ,
	"PH": Unit_mm,
	"PI": Unit_dBZ,
	"PK": Unit_dBZ,
	"PL": Unit_dBZ,
	"PM": Unit_dBZ,
	"PN": Unit_dBZ,
	"PR": Unit_mps,
	"PU": Unit_mps,
	"PV": Unit_mps,
	"PX": Unit_dBZ,
	"PY": Unit_mm,
	"PZ": Unit_dBZ,
	"RA": Unit_mm,
	"RB": Unit_mm,
	"RE": Unit_mm,
	"RH": Unit_mm,
	"RK": Unit_mm,
	"RL": Unit_mm,
	"RM": Unit_mm,
	"RN": Unit_mm,
	"RQ": Unit_mm,
	"RR": Unit_mm,
	"RU": Unit_mm,
	"RW": Unit_mm,
	"RX": Unit_dBZ,
	"RY": Unit_mm,
	"RZ": Unit_mm,
	"S2": Unit_mm,
	"S3": Unit_mm,
	"SF": Unit_mm,
	"SH": Unit_mm,
	"SQ": Unit_mm,
	"TB": Unit_mm,
	"TH": Unit_mm,
	"TW": Unit_mm,
	"TX": Unit_dBZ,
	"TZ": Unit_mm,
	"W1": Unit_mm,
	"W2": Unit_mm,
	"W3": Unit_mm,
	"W4": Unit_mm,
	"WN": Unit_dBZ,
	"WX": Unit_dBZ,
	"YW": Unit_mm,
}
