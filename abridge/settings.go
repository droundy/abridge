package main

import (
	"fmt"
	"http"
	"template"
	"json"
	"github.com/droundy/bridge"
)

type Settings struct {
	Style string
}

func getSettings(req *http.Request) (p Settings) {
	req.ParseForm()
	p = Settings{"two color"}
	prefstr, _ := req.Header["Cookie"] // I don't care about errors!
	// fmt.Println("Trying to unmarshall string", prefstr)
	json.Unmarshal([]byte(prefstr), &p) // I don't care about errors!
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
	if s,ok := req.Form["style"]; ok {
		p.Style = s[0]
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

	conventionCard(c, req)

	fmt.Fprintln(c, `<input type="submit" value="Save settings" />`)
	fmt.Fprintln(c, `</div></form>`)
	fmt.Fprintln(c, `</div>`)
}

type ConventionCard struct {
	Name string
	GeneralApproach string
}

func conventionCard(c *http.Conn, req *http.Request) {
	cctemplate := `
<table class="cc" width="100%"><tr>
<td style="width:50%" class="cc"><table width="100%" class="cc"><tr>

	<td width="50%" class="cc">
	<center><strong>Special doubles</strong></center>
	After overcall:<br/>
	Negative:<br/>
	Responsive:<br/>
	</td><td class="cc">
	<center><strong>Notrump overcalls</strong></center>

  <strong>Direct:</strong>
  <input type="text" name="DirectOvercallNTmin" maxlength="2" size="2" value="15"/> to
  <input type="text" name="DirectOvercallNTmax" maxlength="2" size="2" value="17"/>
  Systems on <input type="checkbox" name="NTOvercallSystemsOn" checked="checked"/>
  <br/>

  <strong>Balancing:</strong>
  <input type="text" name="BalancingOvercallNTmin" maxlength="2" size="2" value="15"/> to
  <input type="text" name="BalancingOvercallNTmax" maxlength="2" size="2" value="17"/>
  <br/>

  Jump to 2NT:
  Minors <input type="checkbox" name="OvercallJump2NT_minors"/>
  2 lowest <input type="checkbox" name="OvercallJump2NT_two_lowest"/>
	</td></tr>

	<tr><td class="cc">
	<center><strong>Simple overcall</strong></center>
	1 level<br/>
	often 4 cards<br/>
	</td><td class="cc">
	<center><strong>Defense vs notrump</strong></center>
	</td></tr>

	<tr><td class="cc">
	<center><strong>Jump overcall</strong></center>
	
	</td><td rowspan="2" class="cc">
	
	<center><strong>over opps t/o double</strong></center>
	<p>foo</p>
	</td></tr>
	<tr><td class="cc">
	<center><strong>Opening preempts</strong></center>
	</td></tr>

	<tr><td class="cc">
	<center><strong>Direct cuebid</strong></center>
	
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
  <input type="text" name="GeneralApproach" size="50" value="{GeneralApproach}"/>
	
	</td></tr>

	<tr class="cc"><td colspan="2" class="cc">
	<center><strong>No trump opening bids</strong></center>
	<table width="100%" border="0">
    <tr><td width="33%" align="center">
          1NT
        </td><td width="33%"></td><td>
           <strong>2NT</strong>
           <input type="text" name="TwoNTmin" maxlength="2" size="2" value="20"/> to
           <input type="text" name="TwoNTmax" maxlength="2" size="2" value="22"/>
        </td>
    </tr><tr>
      <td align="center">
        <input type="text" name="OneNTmin" maxlength="2" size="2" value="15"/> to
        <input type="text" name="OneNTmax" maxlength="2" size="2" value="17"/>
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
        </td><td>
              Jacobi <input type="checkbox" checked="checked" disabled="disabled"/>
           Texas <input type="checkbox" checked="checked" disabled="disabled"/>
         </td>
    </tr><tr>
      <td>5-card major common:
        <input type="checkbox" disabled="disabled"/>
      </td><td>
        ` + htmlbid("3H") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
      </td>
    </tr><tr>
      <td>
        `+htmlbid("2C")+` Stayman
        <input type="checkbox" checked="checked" disabled="disabled"/>
      </td><td>
        ` + htmlbid("3S") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
      </td>
    </tr><tr>
      <td>
        `+htmlbid("2D")+` transfer to `+bridge.SuitColorHTML[bridge.Hearts]+ `
        <input type="checkbox" checked="checked" disabled="disabled"/>
      </td>
    </tr><tr>
      <td>
        `+htmlbid("2H")+` transfer to `+bridge.SuitColorHTML[bridge.Spades]+ `
        <input type="checkbox" checked="checked" disabled="disabled"/>
      </td><td>
        `+htmlbid("4D")+`,`+htmlbid("4H")+` transfer:
        <input type="checkbox" disabled="disabled"/>
      </td><td>
           <strong>3NT</strong>
           <input type="text" name="ThreeNTmin" maxlength="2" size="2" value=""/> to
           <input type="text" name="ThreeNTmax" maxlength="2" size="2" value=""/>
       </td>
    </tr><tr>
      <td>
        ` + htmlbid("2S") + `<input type="text" disabled="disabled" maxlength="10" size="10"/>
      </td><td>
      </td><td>
           Gambling <input type="checkbox" checked="checked" disabled="disabled"/>
       </td>
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
	<center><strong>Responses</strong></center>
  Double raise:
  Force <input type="radio" name="MajorDoubleRaise" disabled="disabled"/>
  Inv. <input type="radio" name="MajorDoubleRaise" checked="checked" disabled="disabled"/>
  Weak <input type="radio" name="MajorDoubleRaise" disabled="disabled"/>
  <br/>

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
	<center><strong>Responses</strong></center>
	
	</td></tr>

	<tr class="cc"><td colspan="2" class="cc">
	
	<center><strong>Describe</strong></center>
	
	</td></tr>

	</table></td>
	
</tr></table>
`
	mycard := ConventionCard{Name:""}
	template.MustParse(cctemplate, nil).Execute(mycard, c)
}
