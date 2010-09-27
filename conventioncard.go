package bridge

type ConventionCard struct {
	Name string
	GeneralApproach string
	Pts map[string]Points
	Options map[string]bool
	Radio map[string]string
}

func (x *ConventionCard) SameAs(y *ConventionCard) bool {
	for k,v := range x.Pts {
		if y.Pts[k] != v {
			return false
		}
	}
	for k,v := range x.Options {
		if y.Options[k] != v {
			return false
		}
	}
	for k,v := range x.Radio {
		if y.Radio[k] != v {
			return false
		}
	}
	return true
}

func DefaultConvention() (out ConventionCard) {
	out.Name = "Default"
	out.GeneralApproach = "David's standard American"
	out.Pts = make(map[string]Points)
	out.Pts["DirectOvercallNTmin"] = 15
	out.Pts["DirectOvercallNTmax"] = 17
	out.Pts["BalancingOvercallNTmin"] = 15
	out.Pts["BalancingOvercallNTmax"] = 17
	out.Pts["OneNTmin"] = 15
	out.Pts["OneNTmax"] = 17
	out.Pts["TwoNTmin"] = 20
	out.Pts["TwoNTmax"] = 22
	out.Pts["OneNTover1Cmin"] = 6
	out.Pts["OneNTover1Cmax"] = 9
	out.Pts["Overcallmin"] = 6
	out.Pts["Overcallmax"] = 17
	out.Options = make(map[string]bool)
	out.Options["Stayman"] = true
	out.Options["Jacobi"] = true
	out.Options["Blackwood"] = true
	out.Options["Gerber"] = false
	out.Options["NTOvercallSystemsOn"] = true
	out.Options["OneNT5CardMajor"] = false
	out.Options["JacobiTransfer2NT"] = true
	out.Options["Texas"] = false
	out.Options["Splinter"] = true
	out.Options["Jacobi2NT"] = false
	out.Options["Bypass4diamonds"] = true
	out.Options["VeryLightOpenings"] = false
	out.Options["VeryLightThirdHand"] = false
	out.Options["VeryLightOvercalls"] = false
	out.Options["VeryLightPreempts"] = true
	out.Options["StrongTwoClubs"] = true
	out.Options["StrongTwos"] = false
	out.Options["FourCardOvercalls"] = false
	out.Options["Gambling3NT"] = true
	out.Radio = make(map[string]string)
	out.Radio["MajorDoubleRaise"] = "Invitational"
	out.Radio["MajorAfterOvercall"] = "Invitational"
	out.Radio["MinorDoubleRaise"] = "Invitational"
	out.Radio["MinorAfterOvercall"] = "Invitational"
	out.Radio["OvercallNewSuit"] = "Force"
	out.Radio["MinorCuebid"] = "Michaels"
	out.Radio["MajorCuebid"] = "Michaels"
	out.Radio["WeakThree"] = "Light"
	out.Radio["JumpOvercall"] = "Weak"
	out.Radio["JumpShiftOverTOX"] = "Weak"
	out.Radio["NewSuitForcingOverTOX"] = "TwoLevel"
	return
}
