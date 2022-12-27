{{ define "partial/heading.tpl" }}
    <h1 class="text-center">{{ .heading }}</h1>
    <h2 class="text-center">
        <small>{{ .intro }}</small>
    </h2>
    <hr>
{{ end }}
