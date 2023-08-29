```
dead simple command line tool for go's standard text/template

install:
$ go install go.senan.xyz/tmplc@latest

usage:
$ tmplc [FILE]...

FILE can be any number of "partials" or "leaves". where leaves have a ".tmpl" suffix, and partials are everything else.

after parsing all partials and leaves, new files will be created for the leaves without the ".tmpl" suffix.

it's even possible to use it as a static site generator:
$ cat partials
{{ define "layout" }}
<html>
  <head><title>{{ template "title" }}</title></head>
  <body>{{ template "content" }}</body>
</html>
{{ end }}

$ cat page-1.html.tmpl
{{ define "title" }}page one{{ end }}
{{ define "content" }}hello from <i>page 1</i>{{ end }}
{{ template "layout" }}

$ cat page-2.html.tmpl
{{ define "title" }}page two{{ end }}
{{ define "content" }}hello from <i>page 2</i>{{ end }}
{{ template "layout" }}

$ tmplc partials *.tmpl

$ cat page-1.html
<html>
  <head><title>page one</title></head>
  <body>hello from <i>page 1</i></body>
</html>

$ cat page-2.html
<html>
  <head><title>page two</title></head>
  <body>hello from <i>page 2</i></body>
</html>
```
