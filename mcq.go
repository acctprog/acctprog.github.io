package main

import (
    "bytes"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func main() {
    pathwd, _ := os.Getwd()
    ext := "_mcq.txt"
	files := checkExt(ext, pathwd)
	for _, fnum := range files {
	    makeMcq(fnum + ext, "m" + fnum + ".html")
	}
}

func checkExt(ext string, pathwd string) []string {
	var files []string
	filepath.Walk(pathwd, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), ext) {
				files = append(files, strings.TrimSuffix(f.Name(), ext))
			}
		}
		return nil
	})
	return files
}

type Mcq struct {
    N int
    S int
    Q string
    A string
    B string
    C string
    D string
    X string
}

func readMcqs(fin string) []Mcq {
    content, _ := ioutil.ReadFile(fin)
	sin := strings.Replace(string(content), "\r", "", -1)
	mcqs := strings.Split(sin, "\n\n")
	var models []Mcq
	for i, mcq := range mcqs {
	    lines := strings.Split(mcq, "\n")
	    numLines := len(lines)
	    if (numLines != 7) {
	        if (numLines > 999) {
	            for _, l := range lines {
	                println(l)
	            }
	            panic("mcq " + strconv.Itoa(i) + ": " + strconv.Itoa(numLines))
	        }
	        //continue
	    }
	    var model Mcq
	    model.N = i
	    model.S = i + 1
	    model.Q = html.EscapeString(strings.TrimPrefix(lines[0], "Q. "))
	    model.A = html.EscapeString(strings.TrimPrefix(lines[1], "A. "))
	    model.B = html.EscapeString(strings.TrimPrefix(lines[2], "B. "))
	    model.C = html.EscapeString(strings.TrimPrefix(lines[3], "C. "))
	    model.D = html.EscapeString(strings.TrimPrefix(lines[4], "D. "))
	    model.X = strings.TrimPrefix(lines[5], "X. ")
	    models = append(models, model)
	}
	return models
}

func writeMcqs(fout string, models []Mcq) {
    const s0 =`<html>
<head>
</head>
<script>

var exp = [`

    const s2 = `];

function doshow_el(str, doshow) {
  if (doshow == 0) {
    document.getElementById(str).style.display = "none";
  } else {
    document.getElementById(str).style.display = "block";
  }
}

function set_el(str, txt) {
  document.getElementById(str).innerHTML = txt;
}

function sum_arr(arr) {
  var res = 0;
  for (var i = 0; i < arr.length; i++) {
    res += arr[i];
  }
  return res;
}

var box_colors = ["#ffb3af", "#b6edb6"];

function style_box(str, sco) {
    document.getElementById(str).style.backgroundColor = box_colors[sco];
}

function get_rad_value(str) {
  var radios = document.getElementsByName(str);
  for (var i = 0; i < 4; i++) {
    if (radios[i].checked) {
      return radios[i].value;
    }
  }
  return "No answer";
}

var you = [];
var sc_sco = [];

function gen_scores() {
  var qlength = exp.length;
  for (var i = 0; i < qlength; i++) {
      you[i] = get_rad_value("d_" + i);
      set_el("d_ya" + i, you[i]);
      set_el("d_ca" + i, exp[i]);
      sc_sco[i] = (exp[i] == you[i]) ? 1 : 0;
      style_box("d_b" + i, sc_sco[i]);
  }
  var tscore = sum_arr(sc_sco);
  var mscore = sc_sco.length;
  set_el("d_tscore", tscore);
  set_el("d_mscore", mscore);
}

var slide_index = -1;
var num_slides = 3;

function goNext() {
  slide_index++;
  if (slide_index == num_slides - 1) {
    gen_scores();
  }
  for (var i = 0; i < num_slides; i++) {
    doshow_el("ds_" + i, (slide_index == i ? 1 : 0));
  }
}

function goPrev() {
  slide_index--;
  for (var i = 0; i < num_slides; i++) {
    doshow_el("ds_" + i, (slide_index == i ? 1 : 0));
  }
}

</script>

<div id="ds_0">
  <div id="d_story" class="d_boxed">
    Answer the multiple choice questions on the next page.
  </div>
  <button type="button" class="c_nav" onclick="window.location.href='index.html'">Back to menu</button>
  <button type="button" class="c_nav" onclick="goNext()">Next</button>
</div>

<div id="ds_1">
`

    const s4 = `  <button type="button" class="c_nav" onclick="goPrev()">Prev</button>
  <button type="button" class="c_nav" onclick="goNext()">Next</button>
</div>

<div id="ds_2">
`  
    const s6 = `  <div class="d_boxed">
    <div><b>Your score:</b> <span id="d_tscore"></span></div>
    <div><b>Maximum score:</b> <span id="d_mscore"></span></div>
  </div>
  <button type="button" class="c_nav" onClick="goPrev()">Prev</button>
  <button type="button" class="c_nav" onClick="window.location.reload();">Retry</button>
  <button type="button" class="c_nav" onclick="window.location.href='index.html'">Return to menu</button>
</div>

<script>

goNext();

</script>

<style>
body {
  font-family: sans-serif;
  background-color: #ffdd22;
}

.c_nav {
  background-color: #4CAF50;
  border: none;
  color: white;
  padding: 15px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  cursor: pointer;
  border-radius: 8px;
  margin-top: 8px;
}

.d_boxed {
  margin-top: 8px;
  border: 1px solid green;
  width: 40em;
  padding: 15px;
  line-height: 1.5;
  background-color: white;
}
</style>
`
    const tmpl0 = `  <div class="d_boxed">
    <div>{{.S}}. {{.Q}}</div>
    <br><label><input type=radio name=d_{{.N}} value=A />A. {{.A}}</label>
    <br><label><input type=radio name=d_{{.N}} value=B />B. {{.B}}</label>
    <br><label><input type=radio name=d_{{.N}} value=C />C. {{.C}}</label>
    <br><label><input type=radio name=d_{{.N}} value=D />D. {{.D}}</label>
  </div>
`
    const tmpl1 = `  <div class="d_boxed" id="d_b{{.N}}">
    <div>{{.S}}. {{.Q}}</div>
    <br><label>A. {{.A}}</label>
    <br><label>B. {{.B}}</label>
    <br><label>C. {{.C}}</label>
    <br><label>D. {{.D}}</label>
    <br><br>
    <div><b>Your answer:</b> <span id=d_ya{{.N}}></span></div>
    <div><b>Correct answer:</b> <span id=d_ca{{.N}}></span></div>
  </div>
`
    var buf0 bytes.Buffer
    var buf1 bytes.Buffer
    var buf2 bytes.Buffer
    t0 := template.Must(template.New("tmpl0").Parse(tmpl0))
    t1 := template.Must(template.New("tmpl1").Parse(tmpl1))
    for _, m := range models {
    	t0.Execute(&buf0, m)
    	t1.Execute(&buf1, m)
    	buf2.WriteString("\"")
    	buf2.WriteString(m.X)
    	buf2.WriteString("\", ")
    }
    s3 := buf0.String()
    s5 := buf1.String()
    s1 := buf2.String()
    
    sout := s0 + s1 + s2 + s3 + s4 + s5 + s6
    
    cout := []byte(sout)
	ioutil.WriteFile(fout, cout, 0644)
}

func makeMcq(fin string, fout string) {
    models := readMcqs(fin)
    writeMcqs(fout, models)
}
