{{ define "plugins/highlightjs-js.tpl" -}}
<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/{{ .plugin.highlightjs.version}}/highlight.min.js"></script>
<script>hljs.initHighlightingOnLoad();</script>
{{- end }}
