{{ define "partial/posts.tpl" }}

    {{ range .posts }}
        <h2><a href="{{ .Path }}">{{ .Title }}</a></h2>

        <div class="clearfix">
            {{ if $.plugins.disqus }}
                <p class="float-right">
                    <a href="{{ .Path }}#disqus_thread" data-disqus-identifier="{{ .ID }}">Comments</a>
                </p>
            {{ end }}

            {{ template "partial/info.tpl" . }}
        </div>

        <p>
            {{- if .TeaserHTML -}}
                {{ .TeaserHTML | noescape }} <a href="{{ .Path }}">[Read more]</a>
            {{- else -}}
                {{ .BodyHTML | striptags | truncatechars 255 }} <a href="{{ .Path }}">[Read more]</a>
            {{- end -}}
        </p>

        <hr>
    {{ end }}

    {{ if $.plugins.disqus }}
        <script id="dsq-count-scr" src="//{{ .plugin.disqus.shortname }}.disqus.com/count.js" async></script>
    {{ end }}

{{ end }}
