// Copyright 2015 Auburn University. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"text/template"
)

type UserGuideContent struct {
	docContent
	ManPageHTML string
	VimdocHTML  string
}

// PrintUserGuide outputs the User's Guide for the Go Doctor (in HTML).
//
// Both the godoctor man page and the Vim plugin reference are generated and
// included in the User's Guide.  The man page content is piped through groff
// to convert it to HTML.
func PrintUserGuide(aboutText string, flags *flag.FlagSet, out io.Writer) {
	PrintUserGuideAsGiven(aboutText, flags, &UserGuideContent{}, out)
}

// PrintUserGuideAsGiven outputs the User's Guide for the Go Doctor (in HTML).
// However, if the content's ManPageHTML and/or VimdocHTML is nonempty, the
// given content is used rather than generating the content.  This is used by
// the online documentation, which cannot execute groff to convert the man page
// to HTML (due to an App Engine restriction), and which uses a Vim-colored
// version of the Vim plugin documentation.
func PrintUserGuideAsGiven(aboutText string, flags *flag.FlagSet, ctnt *UserGuideContent, out io.Writer) {
	ctnt.docContent = prepare(aboutText, flags)
	if ctnt.ManPageHTML == "" {
		ctnt.ManPageHTML = extractBetween(convertManPage(aboutText, flags),
			"<body>", "</body>")
	}
	if ctnt.VimdocHTML == "" {
		ctnt.VimdocHTML = fmt.Sprintf(
			"<pre>\n%s\n</pre>",
			printVimdoc(aboutText, flags))
	}

	tmpl := template.Must(template.New("userGuide").Parse(userGuide))
	err := tmpl.Execute(out, ctnt)
	if err != nil {
		panic(err)
	}
}

func convertManPage(aboutText string, flags *flag.FlagSet) string {
	cmd := exec.Command("groff", "-t", "-mandoc", "-Thtml")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Sprintf("[ERROR] %s", err.Error())
	}

	go func() {
		defer func() {
			recover()
		}()
		PrintManPage(aboutText, flags, stdin)
		err = stdin.Close()
	}()

	result, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("[ERROR] %s", err.Error())
	}
	return string(result)
}

func extractBetween(s, from, to string) string {
	i := strings.Index(s, from)
	if i < 0 {
		return ""
	}

	j := strings.LastIndex(s, to)
	if j < 0 {
		return ""
	}

	return s[i+len(from) : j]
}

func printVimdoc(aboutText string, flags *flag.FlagSet) string {
	defer func() {
		recover()
	}()
	var b bytes.Buffer
	PrintVimdoc(aboutText, flags, &b)
	return b.String()
}

const userGuide = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
  <title>{{.AboutText}} User's Guide</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
  <style>
  html {
    font-family: Arial;
    font-size: 0.688em;
    line-height: 1.364em;
    background-color: white;
    color: black;
  }
  a {
    color: black;
    text-decoration: underline;
  }
  a: hover {
    text-decoration: none;
  }
  tt {
    font-size: 1.2em;
  }
  h1 {
    text-align: center;
    font-size: 2.6em;
    font-weight: bold;
    padding: 20px 0 20px 0;
    background-color: #e0e0e0;
  }
  h2 {
    text-align: left;
    font-size: 2.5em;
    font-weight: bold;
    padding: 10px 0 10px 0;
    background-color: #e0e0e0;
  }
  h3 {
    text-align: left;
    font-size: 1.75em;
    font-weight: bold;
    padding: 5px 0 2px 0;
    border-bottom: 2px dashed #c0c0c0;
  }
  h4 {
    text-align: left;
    font-size: 1.5em;
    font-weight: bold;
    padding: 5px 0 0 0;
  }
  .highlight {
    background-color: yellow;
  }
  .dotted {
    border: 1px dotted;
  }
 
  .clicktoshow {
    display: show;
    //float: right;
    font-size: 10px;
    font-weight: normal;
    color: #808080;
  }
  .showable {
    display: none;
  }

  #toc2col .column1 {
    width: 250px;
    padding: 0;
    position: fixed;
    right: 0px;
    top: 0px;
  }
  #toc2col .column2 {
    width: 628px;
    padding: 10px 0 10px 0;
  }

  .box {
    background-color: #c0c0c0;
    width: 210px;
  }
  .box h2 {
    text-align: center;
    font-size: 1.364em;
    line-height: 1em;
    font-weight: bold;
    padding: 3px 0 3px 0;
    margin-top: 40px;
    color:#ffffff;
    background-color: #000000;
  }
  .box ul       { list-style: none; padding: 0; margin: 0; }
  .box ul li    { padding: 5px 0 1px 15px; font-weight: bold;}
  .box ul ul li { padding: 1px 0 1px 30px; font-weight: normal;}

  .man h1 {
    text-align: center;
    font-size: 1.8em;
    line-height: 1em;
    font-weight: bold;
    padding: 3px 0 3px 0;
    margin-top: 5px;
    background-color: #ffffff;
    color: black;
  }
  .man h2 {
    text-align: left;
    font-size: 1.4em;
    line-height: 1em;
    font-weight: bold;
    padding: 3px 0 3px 0;
    margin-top: 20px;
    background-color: #ffffff;
  }

  .vimdoc pre {
    font-size:1.0em;
    margin-left: 20px;
  }
  </style>
  <script language="JavaScript">
    function setDisplay(selectors, value) {
      var divs = document.querySelectorAll(selectors);
      for (var i = 0; i < divs.length; i++) {
        divs[i].style.display = value;
      }
    }

    function show(id) {
      setDisplay('.showable', 'none');
      setDisplay('.clicktoshow', 'block');
      document.getElementById(id).style.display = 'block';
      document.getElementById(id + '-click').style.display = 'none';
    }

    function showAll() {
      setDisplay('.showable', 'block');
      setDisplay('.clicktoshow', 'none');
    }

    function hideAll() {
      setDisplay('.showable', 'none');
      setDisplay('.clicktoshow', 'block');
    }
  </script>
</head>
<body id="toc2col">
    <!-- BEGIN BODY -->
    <div id="middle">
      <div class="container">
        <div class="column1">
          <div class="box">
            <div class="corner_bottom_left">
              <div class="corner_top_right">
                <div class="corner_top_left">
                  <div class="indent">
                    <!-- BEGIN TOC -->
                    <h2>Getting Started</h2>
                    <ul class="toc1">
                      <li><a onClick="show('usage');" href="#usage">Basic Usage</a></li>
                      <ul class="toc2">
                        <li><a onClick="show('usage-vim');" href="#usage-vim">Vim Plug-in</a></li>
                        <li><a onClick="show('usage-cli');" href="#usage-cli">Command Line Tool (godoctor)</a></li>
                      </ul>
                      <li><a onClick="show('help');" href="#help">Getting Help</a></li>
                      <ul class="toc2">
                        <li><a onClick="show('documentation');" href="#documentation">Online Documentation</a></li>
                        <li><a onClick="show('mailing-list');" href="#mailing-list">Joining the Mailing List</a></li>
                      </ul>
                    </ul>
                    <h2>Refactorings</h2>
                    <ul class="toc2">
                      {{range .Refactorings}}
                      <li><a onClick="show('refactoring-{{.Key}}');" href="#refactoring-{{.Key}}">{{.Description.Name}}</a></li>
                      {{end}}
                    </ul>
                    <h2>References</h2>
                    <ul class="toc2">
                      <li><a onClick="show('godoctor-man');" href="#godoctor-man">Man Page (<tt>godoctor.1</tt>)</a></li>
                      <li><a onClick="show('godoctor-vim');" href="#godoctor-vim">Vim Plug-in Reference</a></li>
                      <li><a onClick="show('license');" href="#license">License</a></li>
                    </ul>
		    <p style="text-align: center; font-size: 10px; color: #808080;">
                      <a href="#" onClick="showAll();">Show All</a> |
                      <a href="#" onClick="hideAll();">Hide All</a>
                    </p>
                    <!-- END TOC -->
                    <br/>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="column2">
          <div class="indent">
            <!-- BEGIN CONTENT -->
<h1>{{.AboutText}} User's Guide</h1>
<a name="usage"></a>
<h2>Basic Usage</h2>
<div id="usage-click" class="clicktoshow"></div>
<div id="usage" class="showable"></div>
<a name="usage-vim"></a>
<h3>Using the Go Doctor Vim Plugin</h3>
<div id="usage-vim-click" class="clicktoshow">
  <a href="#usage-vim" onClick="show('usage-vim');">Show&nbsp;&raquo;</a>
</div>
<div id="usage-vim" class="showable">
  <p>(Coming soon)</p>
  <p>For more information about the Vim plug-in, see the <a
  onClick="show('godoctor-vim');" href="#godoctor-vim">Vim Plug-in
  Reference</a>.  If you have installed the Vim plug-in documentation, you can
  access this by starting Vim and executing the command:<br/><br/>
  <tt>:help godoctor</tt></p>
</div>
<a name="usage-cli"></a>
<h3>Using the Command Line Tool (<tt>godoctor</tt>)</h3>
<div id="usage-cli-click" class="clicktoshow">
  <a href="#usage-cli" onClick="show('usage-cli');">Show&nbsp;&raquo;</a>
</div>
<div id="usage-cli" class="showable">
  <p>(Coming soon)</p>
  <p>For more information about the <tt>godoctor</tt> command, see the <a
  onClick="show('godoctor-man');" href="#godoctor-man"><tt>godoctor</tt> man
  page</a>.  If you have installed this man page locally, you can access it
  from a shell prompt by running<br/><br/>
  <tt>man godoctor</tt></p>
</div>
<a name="help"></a>
<h2>Getting Help</h2>
<div id="help-click" class="clicktoshow"></div>
<div id="help" class="showable"></div>
<a name="documentation"></a>
<h3>Online Documentation</h3>
<div id="documentation-click" class="clicktoshow">
  <a href="#documentation" onClick="show('documentation');">Show&nbsp;&raquo;</a>
</div>
<div id="documentation" class="showable">
  <p>Up-to-date documentation for the Go Doctor is available online at
  <a target="_blank" href="http://gorefactor.org">http://gorefactor.org</a>.
  See the installation instructions at that link for how to install a local
  copy of the documentation.</p>
</div>
<a name="mailing-list"></a>
<h3>Joining the Mailing List</h3>
<div id="mailing-list-click" class="clicktoshow">
  <a href="#mailing-list" onClick="show('mailing-list');">Show&nbsp;&raquo;</a>
</div>
<div id="mailing-list" class="showable">
  <p>If you get stuck, please join the <a target="_blank"
  href="http://mailman.eng.auburn.edu/cgi-bin/mailman/listinfo/go-refactoring">go-refactoring
  mailing list</a>; that is a good place to ask for help.  Updates to the Go
  Doctor will also be announced on the mailing list.</li>
</div>
<a name="refactorings"></a>
<h2>Refactorings</h2>
<div id="refactorings"></div>
{{range .Refactorings}}
<a name="refactoring-{{.Key}}"></a>
<h3>Refactoring: {{.Description.Name}}</h3>
<div id="refactoring-{{.Key}}-click" class="clicktoshow">
  <a href="#refactoring-{{.Key}}" onClick="show('refactoring-{{.Key}}');">Show&nbsp;&raquo;</a>
</div>
<div id="refactoring-{{.Key}}" class="showable">
  {{.Description.HTMLDoc}}
</div>
{{end}}
<div id="references">
  <a name="references"></a>
  <h2>References</h2>
</div>
<a name="godoctor-man"></a>
<h3><tt>godoctor</tt> Man Page</h3>
<div id="godoctor-man-click" class="clicktoshow">
  <a href="#godoctor-man" onClick="show('godoctor-man');">Show&nbsp;&raquo;</a>
</div>
<div id="godoctor-man" class="showable">
  <div class="man">
    {{.ManPageHTML}}
  </div>
</div>
<a name="godoctor-vim"></a>
<h3>Vim Plugin Reference</h3>
<div id="godoctor-vim-click" class="clicktoshow">
  <a href="#godoctor-vim" onClick="show('godoctor-vim');">Show&nbsp;&raquo;</a>
</div>
<div id="godoctor-vim" class="showable">
  <div class="vimdoc">
  {{.VimdocHTML}}
  </div>
</div>
<a name="license"></a>
<h3>License</h3>
<div id="license-click" class="clicktoshow">
  <a href="#license" onClick="show('license');">Show&nbsp;&raquo;</a>
</div>
<div id="license" class="showable">
  <p>Copyright &copy; 2015, Auburn University.  All rights reserved.</p>
  <p>Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:</p>
  <p>1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.</p>
  <p>2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.</p>
  <p>3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.</p>
  <p>THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.</p>
</div>
            <!-- END CONTENT -->
          </div>
        </div>
      </div>
    </div>
    <!-- END BODY -->
</body>
</html>
`
