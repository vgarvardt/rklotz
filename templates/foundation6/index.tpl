{{ define "content" }}

    {{ template "partial/heading.tpl" . }}
    {{ template "partial/about.tpl" . }}

    {{ if lt .meta.Posts 1 }}
        <h2>Nothing yet =(</h2>
    {{ else }}
        {{ template "partial/posts.tpl" . }}

        {{ template "partial/pagination.tpl" . }}
    {{ end }}

{{ end }}
