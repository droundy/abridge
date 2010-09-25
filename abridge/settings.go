package main

import (
	"fmt"
	"io"
	"http"
	"template"
	"json"
	"strconv"
	"github.com/droundy/bridge"
)

type Settings struct {
	Style string
	Card bridge.ConventionCard
}

var DefaultSettings = Settings{
Style: "two color",
Card: bridge.DefaultConvention,
}

func getSettings(req *http.Request) (p Settings) {
	req.ParseForm()
	p = DefaultSettings
	prefstr, _ := req.Header["Cookie"] // I don't care about errors!
	// fmt.Println("Trying to unmarshall string", prefstr)
	json.Unmarshal([]byte(prefstr), &p) // I don't care about errors!
	for k,v := range bridge.DefaultConvention.Pts {
		if _,ok := p.Card.Pts[k]; !ok {
			p.Card.Pts[k] = v
		}
	}
	for k,v := range bridge.DefaultConvention.Options {
		if _,ok := p.Card.Options[k]; !ok {
			p.Card.Options[k] = v
		}
	}
	for k,v := range bridge.DefaultConvention.Radio {
		if _,ok := p.Card.Radio[k]; !ok {
			p.Card.Radio[k] = v
		}
	}
	// fmt.Println("Unmarshal gave: %v", e)
	return p
}

func (s *Settings) Set(c *http.Conn) {
	bytes,_ := json.Marshal(*s)
	c.SetHeader("Set-Cookie", string(bytes))
}

func checkRadio(c bool) string {
	if c {
		return ` checked="checked"`
	}
	return ""
}

func settings(c *http.Conn, req *http.Request) {
	p := getSettings(req)
	// Since we may have cached things based on old card (or someone
	// else's card!), we should clear out the cache.
	bridge.ClearCache()
	if _,ok := req.Form["amsubmitting"]; ok {
		if s,ok := req.Form["style"]; ok {
			p.Style = s[0]
		}
		//for k,v := range req.Form {
		//	fmt.Println("Got key", k, "and value", v)
		//}
		if x,ok := req.Form["GeneralApproach"]; ok {
			p.Card.GeneralApproach = x[0]
		}
		for k := range bridge.DefaultConvention.Pts {
			if x,ok := req.Form[k]; ok {
				pts,e := strconv.Atoi(x[0])
				if e == nil && bridge.Points(pts) != p.Card.Pts[k] {
					fmt.Println("Got pts", k, "of", pts)
					p.Card.Pts[k] = bridge.Points(pts)
				}
			}
		}
		if x,ok := req.Form["Jacobi"]; ok {
			switch len(x) {
			case 2: p.Card.Options["Jacobi"] = true;
			case 0: p.Card.Options["Jacobi"] = false;
			}
			p.Card.Options["Jacobi"] = !p.Card.Options["Jacobi"]
		} else {
			p.Card.Options["Jacobi"] = false
		}
		for k := range bridge.DefaultConvention.Options {
			// There are two Jacobi checkboxes, so I treat it specially...
			if k != "Jacobi" {
				_,ok := req.Form[k]
				p.Card.Options[k] = ok
			}
		}
		for k := range bridge.DefaultConvention.Radio {
			if x,ok := req.Form[k]; ok {
				p.Card.Radio[k] = x[0]
			}
		}
	}
	if _,ok := req.Form["revert"]; ok {
		fmt.Println("Reverting to defaults...")
		p = DefaultSettings
	}

	p.Set(c)
	defer header(c, req, "aBridge settings")()
	fmt.Fprintln(c, `<div class="textish">`)
	fmt.Fprint(c, `

<p> I am planning to add various configuration options here.  Ideally,
we'll even support a nice convention card input interface.  For now,
I'll probably just add a toggle to switch between two-color and
four-color suits (which will also test that I'm using CSS consistently).

</p>

`)
	fmt.Fprintf(c, `<form method="post" action="%s"><div>`, req.URL.Path)
	fmt.Fprintln(c, `<fieldset><legend>Suit color style</legend>`)
	for _,s := range []string{"two color", "four color"} {
		fmt.Fprintf(c, `<input type="radio" name="style" title="foobar" value="%s"%s/>%s`,
			s, checkRadio(p.Style == s), s)
	}
	fmt.Fprintln(c, `</fieldset>`)

	conventionCard(c, p)

	fmt.Fprintln(c, `<input type="submit" value="Save settings" />`)
	fmt.Fprintln(c, `<input type="submit" name="revert" value="Revert to default" />`)
	fmt.Fprintln(c, `</div></form>`)
	fmt.Fprintln(c, `</div>`)
}

func conventionCard(c *http.Conn, p Settings) {
	cctemplate := `
<input type="hidden" name="amsubmitting" value="true"/>

<table class="cc" width="100%"><tr>
<td style="width:50%" class="cc"><table width="100%" class="cc"><tr>

	<td width="50%" class="cc">
	<center><strong>Special doubles</strong></center>
	After overcall:<br/>
	Negative:<br/>
	Responsive:<br/>
	</td><td class="cc">
	<center><strong class="unimplemented">Notrump overcalls</strong></center>

  <strong class="unimplemented">Direct:</strong>
{.section Pts}
  <input type="text" name="DirectOvercallNTmin" maxlength="2" size="2" value="{DirectOvercallNTmin}"/> to
  <input type="text" name="DirectOvercallNTmax" maxlength="2" size="2" value="{DirectOvercallNTmax}"/>
{.end}
{.section Options}
  Systems on <input type="checkbox" name="NTOvercallSystemsOn" {NTOvercallSystemsOn}/>
{.end}
  <br/>

  <strong class="unimplemented">Balancing:</strong>
{.section Pts}
  <input type="text" name="BalancingOvercallNTmin" maxlength="2" size="2" value="{BalancingOvercallNTmin}"/> to
  <input type="text" name="BalancingOvercallNTmax" maxlength="2" size="2" value="{BalancingOvercallNTmax}"/>
  <br/>
{.end}

  Jump to 2NT:
  Minors <input type="checkbox" name="OvercallJump2NT_minors"/>
  2 lowest <input type="checkbox" name="OvercallJump2NT_two_lowest"/>
	</td></tr>

	<tr><td class="cc">
	<center><strong>Simple overcall</strong></center>
  {.section Pts}
  1 level <input type="text" name="Overcallmin" maxlength="2" size="2" value="{Overcallmin}"/> to
  <input type="text" name="Overcallmax" maxlength="2" size="2" value="{Overcallmax}"/> HCP
  <br/>
  {.end}
  {.section Options}
  often 4 cards <input type="checkbox" name="FourCardOvercalls" {FourCardOvercalls}/>
  very light style <input type="checkbox" name="VeryLightOvercalls" {VeryLightOvercalls}/>
  {.end}{.section Radio}
  <center><strong class="unimplemented">Responses</strong></center>
  New suit: 
  Forcing <input type="radio" name="OvercallNewSuit" value="Force" {OvercallNewSuit|Force}/>
  NF <input type="radio" name="OvercallNewSuit" value="NotForce" {OvercallNewSuit|NotForce}/>
  {.end}
  <br/>

	</td><td class="cc">
	<center><strong class="unimplemented">Defense vs notrump</strong></center>
	</td></tr>

	<tr><td class="cc">
	<center><strong class="unimplemented">Jump overcall</strong></center>
  {.section Radio}
  Strong <input type="radio" name="JumpOvercall" value="Strong" {JumpOvercall|Strong}/>
  Intermediate
    <input type="radio" name="JumpOvercall" value="Intermediate" {JumpOvercall|Intermediate}/>
  Weak <input type="radio" name="JumpOvercall" value="Weak" {JumpOvercall|Weak}/>
  {.end}

	
	</td><td rowspan="2" class="cc">
	
	<center><strong>over opps t/o double</strong></center>
	</td></tr>
	<tr><td class="cc">
	<center><strong>Opening preempts</strong></center>
  {.section Radio}
  <table border="0" width="100%">
  <tr><td></td>
      <td align="center">Sound</td>
      <td align="center" >Light</td>
      <td align="center" >Very light</td>
  </tr><tr><td>3/4-bids</td>
      <td align="center"><input type="radio" name="WeakThree" value="Sound" {WeakThree|Sound}/></td>
      <td align="center"><input type="radio" name="WeakThree" value="Light" {WeakThree|Light}/></td>
      <td align="center"><input type="radio" name="WeakThree" value="VeryLight" {WeakThree|VeryLight}/></td>
  </tr></table>
  {.end}
	</td></tr>

	<tr><td class="cc">
	<center><strong class="unimplemented">Direct cuebid</strong></center>
  {.section Radio}
  <table border="0" width="100%">
  <tr><td>Over:</td>
      <td align="center">Minor</td>
      <td align="center" >Major</td>
  </tr><tr><td>Natural</td>
      <td align="center"><input type="radio" name="MinorCuebid" value="Natural" {MinorCuebid|Natural}/></td>
      <td align="center"><input type="radio" name="MajorCuebid" value="Natural" {MajorCuebid|Natural}/></td>
  </tr><tr><td>Strong T/O</td>
      <td align="center"><input type="radio" name="MinorCuebid" value="StrongTO" {MinorCuebid|StrongTO}/></td>
      <td align="center"><input type="radio" name="MajorCuebid" value="StrongTO" {MajorCuebid|StrongTO}/></td> 
  </tr><tr><td>Michaels</td>
      <td align="center"><input type="radio" name="MinorCuebid" value="Michaels" {MinorCuebid|Michaels}/></td>
      <td align="center"><input type="radio" name="MajorCuebid" value="Michaels" {MajorCuebid|Michaels}/></td>
  </tr></table>
  {.end}
	
	</td><td class="cc">
	<center><strong>vs opening preempts double is</strong></center>
	</td></tr>

	<tr><td colspan="2" class="cc">
	<center><strong>slam conventions</strong></center>
	</td></tr>
	<tr><td colspan="2" class="cc">
	<center><strong>leads and carding</strong></center>
	</td></tr>
	</table></td>
	
	<td class="cc"><table width="100%" class="cc">
  <tr class="cc"><td colspan="2" class="cc"><strong>Names:</strong> {Name}</td></tr>
	<tr class="cc"><td colspan="2" class="cc">

	<center><strong>General approach</strong></center>
  <input type="text" name="GeneralApproach" size="50" value="{GeneralApproach}"/><br/>
  {.section Options}
  <strong>Very light:</strong>
    <span  class="unimplemented">Openings <input type="checkbox" name="VeryLightOpenings" {VeryLightOpenings}/>
    3rd hand <input type="checkbox" name="VeryLightThirdHand" {VeryLightThirdHand}/>
    Overcalls <input type="checkbox" name="VeryLightOvercalls" {VeryLightOvercalls}/></span>
    Preempts <input type="checkbox" name="VeryLightPreempts" {VeryLightPreempts}/>
  <br/>
  <strong>Forcing opening:</strong>
    `+htmlbid("2C")+` <input type="checkbox" name="StrongTwoClubs" {StrongTwoClubs}/>
    Natural 2 bids <input type="checkbox" name="StrongTwos" {StrongTwos}/>
  <br/>
  {.end}
	
	</td></tr>

	<tr class="cc"><td colspan="2" class="cc">
	<center><strong>No trump opening bids</strong></center>
	<table width="100%" border="0">
    <tr><td width="33%" align="center">
          1NT
        </td><td width="33%"></td><td>
           <strong>2NT</strong>
{.section Pts}
           <input type="text" name="TwoNTmin" maxlength="2" size="2" value="{TwoNTmin}"/> to
           <input type="text" name="TwoNTmax" maxlength="2" size="2" value="{TwoNTmax}"/>
{.end}
        </td>
    </tr><tr>
      <td align="center">
{.section Pts}
        <input type="text" name="OneNTmin" maxlength="2" size="2" value="{OneNTmin}"/> to
        <input type="text" name="OneNTmax" maxlength="2" size="2" value="{OneNTmax}"/>
{.end}
      </td><td>
        ` + htmlbid("3C") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
      </td><td>
        Transfer responses:
      </td>
    </tr><tr>
        <td align="center">
          <input type="text" disabled="disabled" maxlength="2" size="2"/> to
          <input type="text" disabled="disabled" maxlength="2" size="2"/>
        </td><td>
           ` + htmlbid("3D") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
        </td><td>{.section Options}
              <span class="unimplemented">Jacobi <input type="checkbox" name="JacobiTransfer2NT" {JacobiTransfer2NT}/>
           Texas <input type="checkbox" name="Texas" {Texas}/></span>
         </td>{.end}
    </tr><tr>
      <td>5-card major common:{.section Options}
        <input type="checkbox" name="OneNT5CardMajor" {OneNT5CardMajor}/>
      </td><td>{.end}
        ` + htmlbid("3H") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
      </td>
    </tr><tr>
      <td>{.section Options}
        `+htmlbid("2C")+` Stayman
        <input type="checkbox" name="Stayman" {Stayman}/>
      </td><td>{.end}
        ` + htmlbid("3S") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
      </td>
    </tr><tr>
      <td>{.section Options}
        `+htmlbid("2D")+` transfer to `+bridge.SuitColorHTML[bridge.Hearts]+ `
        <input type="checkbox" name="Jacobi" value="2D" {Jacobi}/>
      </td>{.end}
    </tr><tr>
      <td>{.section Options}
        `+htmlbid("2H")+` transfer to `+bridge.SuitColorHTML[bridge.Spades]+ `
        <input type="checkbox" name="Jacobi" value="2H" {Jacobi}/>
      </td><td>{.end}
        `+htmlbid("4D")+`,`+htmlbid("4H")+` transfer:
        <input type="checkbox" disabled="disabled"/>
      </td><td>
           <strong class="unimplemented">3NT</strong>
           <input type="text" name="ThreeNTmin" maxlength="2" size="2" value=""/> to
           <input type="text" name="ThreeNTmax" maxlength="2" size="2" value=""/>
       </td>
    </tr><tr>
      <td>
        ` + htmlbid("2S") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
      </td><td>
      </td><td>{.section Options}
           Gambling <input type="checkbox" name="Gambling3NT" {Gambling3NT}/>
       </td>{.end}
    </tr><tr>
      <td>
        ` + htmlbid("2N") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
      </td>
    </tr>
    </table>
	</td></tr>

	<tr class="cc"><td width="50%" class="cc">
	
	<center><strong>Major opening</strong></center>
  <table border="0" width="100%">
  <tr><td>Expected Min. Length</td>
      <td align="center">4</td>
      <td align="center" >5</td></tr>
  <tr><td>1st/2nd</td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
      <td align="center"><input type="checkbox" checked="checked" disabled="disabled"/></td></tr>
  <tr><td>3rd/4th</td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
      <td align="center"><input type="checkbox" checked="checked" disabled="disabled"/></td></tr>
  </table>
	<center><strong class="unimplemented">Responses</strong></center>
  Double raise:{.section Radio}
  Force <input type="radio" name="MajorDoubleRaise" value="Force" {MajorDoubleRaise|Force}/>
  Inv. <input type="radio" name="MajorDoubleRaise" value="Invitational" {MajorDoubleRaise|Invitational}/>
  Weak <input type="radio" name="MajorDoubleRaise" value="Weak" {MajorDoubleRaise|Weak}/>
  <br/>{.end}
  After overcall:{.section Radio}
  Force <input type="radio" name="MajorAfterOvercall" value="Force" {MajorAfterOvercall|Force}/>
  Inv. <input type="radio" name="MajorAfterOvercall" value="Invitational" {MajorAfterOvercall|Invitational}/>
  Weak <input type="radio" name="MajorAfterOvercall" value="Weak" {MajorAfterOvercall|Weak}/>
  <br/>{.end}
  Conv. Raise:{.section Options}
  Jacobi 2NT <input type="checkbox" name="Jacobi2NT" {Jacobi2NT}/>
  Splinter <input type="checkbox" name="Splinter" {Splinter}/>
  <br/>{.end}

	</td>
	<td class="cc">
	
	<center><strong>Minor opening</strong></center>
  <table border="0" width="100%">
  <tr><td>Expected Min. Length</td>
      <td align="center">4</td>
      <td align="center">5</td>
      <td align="center">0-2</td>
      <td align="center">Conv</td>
  </tr>
  <tr><td>` + htmlbid("1C") + `</td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
  </tr>
  <tr><td>` + htmlbid("1D") + `</td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
      <td align="center"><input type="checkbox" disabled="disabled"/></td>
  </tr>
  </table>
	<center><strong class="unimplemented">Responses</strong></center>
  Double raise:{.section Radio}
  Force <input type="radio" name="MinorDoubleRaise" value="Force" {MinorDoubleRaise|Force}/>
  Inv. <input type="radio" name="MinorDoubleRaise" value="Invitational" {MinorDoubleRaise|Invitational}/>
  Weak <input type="radio" name="MinorDoubleRaise" value="Weak" {MinorDoubleRaise|Weak}/>
  <br/>{.end}
  After overcall:{.section Radio}
  Force <input type="radio" name="MinorAfterOvercall" value="Force" {MinorAfterOvercall|Force}/>
  Inv. <input type="radio" name="MinorAfterOvercall" value="Invitational" {MinorAfterOvercall|Invitational}/>
  Weak <input type="radio" name="MinorAfterOvercall" value="Weak" {MinorAfterOvercall|Weak}/>
  <br/>{.end}
  {.section Options}
  Frequently bypass 4+`+bridge.SuitColorHTML[bridge.Diamonds]+`
  <input type="checkbox" name="Bypass4diamonds" {Bypass4diamonds}/>
  <br/>{.end}
  {.section Pts}
  1NT/`+htmlbid("1C")+`
  <input type="text" name="OneNTover1Cmin" maxlength="2" size="2" value="{OneNTover1Cmin}"/> to
  <input type="text" name="OneNTover1Cmax" maxlength="2" size="2" value="{OneNTover1Cmax}"/>
  <br/>
  {.end}
	
	</td></tr>

	<tr class="cc"><td colspan="2" class="cc">
	
	<center><strong>Describe</strong></center>
	
	</td></tr>

	</table></td>
	
</tr></table>
`
	e := template.MustParse(cctemplate, myformatter).Execute(p.Card, c)
	if e != nil {
		fmt.Println("Template error:", e)
	}
}

var myformatter = template.FormatterMap(map[string]func(io.Writer, interface{}, string){
	"checked": func(c io.Writer, v interface{}, format string) {
		if b,ok := v.(bool); ok && b {
			fmt.Fprint(c, `checked="checked"`)
		}
	},
  "html": template.HTMLFormatter,
  "str":  template.StringFormatter,
	"": func(c io.Writer, v interface{}, format string) {
		if b,ok := v.(bool); ok {
			if b {
				fmt.Fprint(c, `checked="checked"`)
			}
			return
		}
		template.StringFormatter(c, v, format)
	},
	"Sound": compareStringThing,
	"Light": compareStringThing,
	"VeryLight": compareStringThing,
	"Natural": compareStringThing,
	"StrongTO": compareStringThing,
	"Michaels": compareStringThing,
	"NotForce": compareStringThing,
	"Force": compareStringThing,
	"Invitational": compareStringThing,
	"Weak": compareStringThing,
	"Intermediate": compareStringThing,
	"Strong": compareStringThing,
})

func compareStringThing(c io.Writer, v interface{}, format string) {
	if b,ok := v.(string); ok && b == format {
		fmt.Fprint(c, `checked="checked"`)
	}
}
