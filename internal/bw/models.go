package bw

type Stats struct {
	IsQueue            bool  `json:"isqueue"`
	CrowdednessPercent uint8 `json:"percent"`
	Queue              int   `json:"queue"`
}

type BW struct {
	Domain string
	ShortName string
}
