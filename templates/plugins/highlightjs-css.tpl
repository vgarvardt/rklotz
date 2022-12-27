{{ define "plugins/highlightjs-css.tpl" -}}
<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/{{ .plugin.highlightjs.version}}/styles/{{ .plugin.highlightjs.theme}}.min.css">
{{- end }}
