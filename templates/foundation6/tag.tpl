{{ define "title" }}Posts on "{{ .tag }}"{{ end }}

{{ define "content" }}
    {{ template "partial/heading.tpl" . }}
    {{ template "partial/about.tpl" . }}

    <h2>Posts on <span class="label secondary" style="font-size: inherit;">{{ .tag }}</span></h2>

    {{ if lt (len .posts)  1 }}
        <h2>Nothing yet =(</h2>
    {{ else }}
        {{ template "partial/posts.tpl" . }}

        {{ template "partial/pagination.tpl" . }}
    {{ end }}
{{ end }}
