{{ define "plugins/yasha.tpl" -}}
<script src="//yastatic.net/es5-shims/0.0.2/es5-shims.min.js"></script>
<script src="//yastatic.net/share2/share.js"></script>
<div class="ya-share2" data-services="{{ .plugin.yasha.services }}" data-lang="{{ .plugin.yasha.lang }}" data-size="{{ .plugin.yasha.size }}" data-title="{{ .post.Title }}"></div>
{{- end }}
