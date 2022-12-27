{{ define "title"}}404{{ end }}

{{ define "content" }}

    {{ template "partial/heading.tpl" . }}

    <div class="callout alert">
        <h5>4-oh-4</h5>
        <p>{{ .error }}</p>
    </div>

{{ end }}
