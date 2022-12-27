{{ define "title"}}500{{ end }}

{{ define "content" }}

    {{ template "partial/heading.tpl" . }}

    <div class="callout alert">
        <h5>500</h5>
        <p>{{ .error }}</p>
    </div>

{{ end }}
