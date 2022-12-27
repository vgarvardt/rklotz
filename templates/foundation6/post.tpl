{{ define "title" }}{{ .post.Title }}{{ end }}

{{ define "head_extra" }}
    {{ if .plugins.highlightjs }}{{ template "plugins/highlightjs-css.tpl" . }}{{ end }}
{{ end }}

{{ define "content" }}
    <h1 class="text-center">{{ .post.Title }}</h1>

    {{ template "partial/info.tpl" .post }}

    <div>
        {{ .post.BodyHTML | noescape }}
    </div>

    {{ if .plugins.yasha }}<p>{{ template "plugins/yasha.tpl" . }}</p>{{ end }}

    {{ template "partial/about.tpl" . }}

    {{ if .plugins.disqus }}{{ template "plugins/disqus.tpl" . }}{{ end }}

{{ end }}

{{ define "foot_extra" }}
    {{ if .plugins.highlightjs }}{{ template "plugins/highlightjs-js.tpl" . }}{{ end }}
{{ end }}
