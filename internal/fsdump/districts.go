package fsdump

import "github.com/alexey-zayats/claim-parser/internal/formstruct"

// District ...
type District struct {
	ID    int
	Title string
}

// Districts ...
var Districts = map[string]District{
	"5e851500497393b7298b4577": {ID: 2, Title: "Карасунский внутригородской округ"},
	"5e85150d497393b7298b4578": {ID: 4, Title: "Центральный внутригородской округ"},
	"5e8515194973933f2a8b4582": {ID: 5, Title: "Берёзовский сельский округ"},
	"5e85151e497393702d8b457b": {ID: 6, Title: "Елизаветинский сельский округ"},
	"5e851523497393b7298b4579": {ID: 7, Title: "Калининский сельский округ"},
	"5e8515294973933f2a8b4583": {ID: 8, Title: "Пашковский сельский округ"},
	"5e851530497393702d8b457c": {ID: 9, Title: "Старокорсунский сельский округ"},
	"5e85a1b8497393702d8b458b": {ID: 10, Title: "Департамент транспорта и дорожного хозяйства"},
}

// Forms ...
var Forms = map[string]map[string]formstruct.State{
	// Карасунский внутригородской округ ...
	"5e851500497393b7298b4577": {
		"5e851adf6239dt5": formstruct.StateKind,
		"5e851b61bebdet1": formstruct.StateName,
		"5e851b6b850a5t1": formstruct.StateAddress,
		"5e851b7409a23t1": formstruct.StateINN,
		"5e851b7b4a1e4t1": formstruct.StateFIO,
		"5e851b82a9b36t1": formstruct.StatePhone,
		"5e85b13da00c8t1": formstruct.StateEMail,
		"5e8583f01e863t2": formstruct.StateCars,
		"5e851ba90e1a7t4": formstruct.StateAgreement,
		"5e851bc2eb2d1t4": formstruct.StateReliability,
	},
	// Центральный внутригородской округ ...
	"5e85150d497393b7298b4578": {
		"5e851e7b2e992t5": formstruct.StateKind,
		"5e851e81e73d0t1": formstruct.StateName,
		"5e851e880123at1": formstruct.StateAddress,
		"5e851e8ca3567t1": formstruct.StateINN,
		"5e851e9166b2ct1": formstruct.StateFIO,
		"5e851e9808846t1": formstruct.StatePhone,
		"5e85b16b87f92t1": formstruct.StateEMail,
		"5e8585a2cc48ft2": formstruct.StateCars,
		"5e851eb80a60dt4": formstruct.StateAgreement,
		"5e851ecca2b5et4": formstruct.StateReliability,
	},
	// Берёзовский сельский округ ...
	"5e8515194973933f2a8b4582": {
		"5e851eecd79b0t5": formstruct.StateKind,
		"5e85200cbad6ct1": formstruct.StateName,
		"5e8520147d9fat1": formstruct.StateAddress,
		"5e85201ace44ct1": formstruct.StateINN,
		"5e8520214a7fbt1": formstruct.StateFIO,
		"5e852028d76d3t1": formstruct.StatePhone,
		"5e85b17c82904t1": formstruct.StateEMail,
		"5e8585c80afcct2": formstruct.StateCars,
		"5e8520565c3b2t4": formstruct.StateAgreement,
		"5e85206da8ec3t4": formstruct.StateReliability,
	},
	// Елизаветинский сельский округ ...
	"5e85151e497393702d8b457b": {
		"5e85208d5263ct5": formstruct.StateKind,
		"5e8520937a4fbt1": formstruct.StateName,
		"5e85209940f83t1": formstruct.StateAddress,
		"5e8520a05a890t1": formstruct.StateINN,
		"5e8520a7c4a33t1": formstruct.StateFIO,
		"5e8520adcff38t1": formstruct.StatePhone,
		"5e85b18ad1b3ft1": formstruct.StateEMail,
		"5e85860c2c2cbt2": formstruct.StateCars,
		"5e8520cc09037t4": formstruct.StateAgreement,
		"5e8520da4c994t4": formstruct.StateReliability,
	},
	// Калининский сельский округ ...
	"5e851523497393b7298b4579": {
		"5e8520fd4e0e3t5": formstruct.StateKind,
		"5e852104bf478t1": formstruct.StateName,
		"5e85210cbcdc7t1": formstruct.StateAddress,
		"5e852114030edt1": formstruct.StateINN,
		"5e85211a0cd32t1": formstruct.StateFIO,
		"5e8521212c652t1": formstruct.StatePhone,
		"5e85b1a5c9b9ft1": formstruct.StateEMail,
		"5e85863d87091t2": formstruct.StateCars,
		"5e85213ea2a5ct4": formstruct.StateAgreement,
		"5e85215277878t4": formstruct.StateReliability,
	},
	// Пашковский сельский округ ...
	"5e8515294973933f2a8b4583": {
		"5e85216e8794et5": formstruct.StateKind,
		"5e8521750f423t1": formstruct.StateName,
		"5e85217b3d7b1t1": formstruct.StateAddress,
		"5e85218292dc5t1": formstruct.StateINN,
		"5e8521881c103t1": formstruct.StateFIO,
		"5e852191c3a3bt1": formstruct.StatePhone,
		"5e85b1b24d297t1": formstruct.StateEMail,
		"5e85866962147t2": formstruct.StateCars,
		"5e8521b12bafdt4": formstruct.StateAgreement,
		"5e8521c0cbba0t4": formstruct.StateReliability,
	},
	// Старокорсунский сельский округ ...
	"5e851530497393702d8b457c": {
		"5e8521d97f66at5": formstruct.StateKind,
		"5e8521e0553a2t1": formstruct.StateName,
		"5e8521e637ec6t1": formstruct.StateAddress,
		"5e8521ed5b0d1t1": formstruct.StateINN,
		"5e8521f3e4d57t1": formstruct.StateFIO,
		"5e8521fc017f7t1": formstruct.StatePhone,
		"5e85b1be280d1t1": formstruct.StateEMail,
		"5e85868085e0dt2": formstruct.StateCars,
		"5e85221c2f151t4": formstruct.StateAgreement,
		"5e85222f6c2e8t4": formstruct.StateReliability,
	},
	// Старокорсунский сельский округ ...
	"5e85a1b8497393702d8b458b": {
		"5e85a232baa8ct5": formstruct.StateKind,
		"5e85a2433eba0t1": formstruct.StateName,
		"5e85a2476f0b7t1": formstruct.StateAddress,
		"5e85a2542637bt1": formstruct.StateINN,
		"5e85a25e8398dt1": formstruct.StateFIO,
		"5e85a265437bft1": formstruct.StatePhone,
		"5e85b1c948b73t1": formstruct.StateEMail,
		"5e85a26ee99f0t1": formstruct.StateCars,
		"5e85a285af770t4": formstruct.StateAgreement,
		"5e85a2966f555t4": formstruct.StateReliability,
	},
}
