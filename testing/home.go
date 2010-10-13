package main

func Home(evt []string) string {
	return `
<h1> Intro to aBridge</h1>

This is a neat thing.
<br/>

  <input type='submit' onclick="say('hello world')" value='Hello.'/> 
  <input type='submit' onclick="say('goodbye world')" value='Goodbye.'/> 
<br/>
(Event is "` + evt[0] + `")
<br/>
`
}
