<!DOCTYPE html>
<html lang="{{ .lang }}">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{ .description }}">
    <meta name="author" content="{{ .author }}">
    {{ block "head_meta" . }}{{ end }}

    <title>{{ .title }} | {{ block "title" . }}{{ .intro }}{{ end }}</title>

    {{ if .plugins.gtm }}{{ template "plugins/gtm-head.tpl" . }}{{ end }}

    <link href="/feed/atom" type="application/atom+xml" rel="alternate" title="Posts (Atom feed)">
    <link href="/feed/rss" type="application/rss+xml" rel="alternate" title="Posts (RSS feed)">

    <link rel="shortcut icon" href="/static/{{ .theme }}/favicon.ico?{{ .instance_id }}">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/foundation-sites@6.6.1/dist/css/foundation.min.css"
          crossorigin="anonymous">
    <link href="https://fonts.googleapis.com/css?family=Play&amp;subset=latin,cyrillic" rel="stylesheet"
          type="text/css">
    <link href='https://fonts.googleapis.com/css?family=Roboto+Mono' rel='stylesheet' type='text/css'>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.3.0/css/all.min.css">
    <link href="/static/{{ .theme }}/rklotz.css?{{ .instance_id }}" rel="stylesheet" type="text/css">
    {{ block "head_extra" . }}{{ end }}
</head>

<body>
{{ if .plugins.gtm }}{{ template "plugins/gtm-body.tpl" . }}{{ end }}

<div class="grid-container grid-x">
    <div class="cell small-12">

        <div class="top-bar">
            <div class="top-bar-left">
                <ul class="menu">
                    <li class="menu-text">
                        <a href="/">
                            <img src="/static/{{ .theme }}/favicon.png?{{ .instance_id }}"
                                 style="max-height: 2.8125rem; width: auto; -webkit-filter: invert(100%); filter: invert(100%);">
                            {{ .heading }}
                        </a>
                    </li>
                </ul>
            </div>

            <div class="top-bar-right">
                <ul class="menu">
                    <li><a href="/feed/rss">rss</a></li>
                </ul>
            </div>
        </div>

        {{ block "content" . }}{{ end }}

        <hr>
        <p class="text-center">
            <small>
                built on top of <a href="https://github.com/vgarvardt/rklotz" target="_blank">rKlotz</a>
                by <a href="http://itskrig.com" target="_blank">Vladimir Garvardt</a>
            </small>
        </p>
    </div>
</div>

<script
    src="https://code.jquery.com/jquery-3.4.1.min.js"
    integrity="sha256-CSXorXvZcTkaix6Yvo6HppcZGetbYMGWSFlBw8HfCJo="
    crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/foundation-sites@6.6.1/dist/js/foundation.min.js"
        crossorigin="anonymous"></script>
{{ block "foot_extra" . }}{{ end }}

{{ if .plugins.ga }}{{ template "plugins/ga.tpl" . }}{{ end }}
{{ if .plugins.yamka }}{{ template "plugins/yamka.tpl" . }}{{ end }}

</body>

</html>
