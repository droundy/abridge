package bridge

type ConventionCard struct {
	Name string
	GeneralApproach string
	Pts map[string]Points
	Options map[string]bool
	Radio map[string]string
}

var DefaultConvention = ConventionCard {
Name: "Default",
GeneralApproach: "David's standard American",
Pts: map[string]Points{
		"DirectOvercallNTmin": 15,
		"DirectOvercallNTmax": 17,
		"BalancingOvercallNTmin": 15,
		"BalancingOvercallNTmax": 17,
		"OneNTmin": 15,
		"OneNTmax": 17,
		"TwoNTmin": 20,
		"TwoNTmax": 22,
		"OneNTover1Cmin": 6,
		"OneNTover1Cmax": 9,
		"Overcallmin": 10,
		"Overcallmax": 18,
	},
Options: map[string]bool{
		"Stayman": true,
		"Jacobi": true,
		"NTOvercallSystemsOn": true,
		"OneNT5CardMajor": false,
		"JacobiTransfer2NT": true,
		"Texas": false,
		"Splinter": true,
		"Jacobi2NT": false,
		"Bypass4diamonds": true,
		"VeryLightOpenings": false,
		"VeryLightThirdHand": false,
		"VeryLightOvercalls": false,
		"VeryLightPreempts": true,
		"StrongTwoClubs": true,
		"StrongTwos": false,
		"FourCardOvercalls": false,
		"Gambling3NT": true,
	},
Radio: map[string]string{
		"MajorDoubleRaise": "Invitational",
		"MajorAfterOvercall": "Invitational",
		"MinorDoubleRaise": "Invitational",
		"MinorAfterOvercall": "Invitational",
		"OvercallNewSuit": "Force",
		"MinorCuebid": "Michaels",
		"MajorCuebid": "Michaels",
		"WeakThree": "Light",
		"JumpOvercall": "Weak",
	},
}
