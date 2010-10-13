package main

func Home(dat *ClientData, evt []string) string {
	return `
<div class="textish">
<h1>aBridge</h1>
<p>
The aBridge program is a hacky little program that I wrote in <a
href="http://golang.org">go</a> to play with.  It's pretty easy to
write servers in go, so that's what I've done, and I figured that I
may as well make it available for others to play with as well.

</p><p>
The source code is available at <a
href="http://github.com/droundy/abridge">github</a> if you're
interested (and I hope you are).
</p><p>

For the curious, there was a former bridge program named abridge that
I wrote, which was quite different, and which is now quite dead.  So I
figured I may as well reuse the name, particularly as I still own the
domain.
</p>
</div>
<br/>
(Event is "` + evt[0] + `")
<br/>
`
}
