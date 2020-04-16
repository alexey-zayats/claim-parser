package model

// OutKind ...
type OutKind int

const (
	// OutUnknown ...
	OutUnknown OutKind = iota
	// OutVehicleClaim ...
	OutVehicleClaim
	// OutVehicleRegistry ...
	OutVehicleRegistry
	// OutPeopleClaim ...
	OutPeopleClaim
	// OutPeopleRegistry ...
	OutPeopleRegistry
)

// Out ...
type Out struct {
	Kind  OutKind
	Event *Event
	Value interface{}
}
