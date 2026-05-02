package models

type Depot struct {
	ID            int `json:"ID"`
	MechanicHours int `json:"MechanicHours"`
}

type DepotsResponse struct {
	Depots []Depot `json:"depots"`
}

type Vehicle struct {
	TaskID   string `json:"TaskID"`
	Duration int    `json:"Duration"`
	Impact   int    `json:"Impact"`
}

type VehiclesResponse struct {
	Vehicles []Vehicle `json:"vehicles"`
}
