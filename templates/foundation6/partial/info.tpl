{{ define "partial/info.tpl" }}
    <p>
        {{ .PublishedAt | format_date }}
        {{ if gt (len .Tags) 0 }}
            on
            {{ range .Tags }}
                <a href="/tag/{{ . }}" class="label secondary">{{ . }}</a>
            {{ end }}
        {{ end }}
    </p>
{{ end }}
