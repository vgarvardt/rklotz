{{ define "plugins/disqus.tpl" -}}
    <div id="disqus_thread"></div>
    <script>
        var disqus_config = function () {
            this.page.url = '{{ .current_url }}';
            this.page.identifier = '{{ .post.ID }}';
        };
        (function () { // DON'T EDIT BELOW THIS LINE
            var d = document, s = d.createElement('script');
            s.src = 'https://{{ .plugin.disqus.shortname }}.disqus.com/embed.js';
            s.setAttribute('data-timestamp', +new Date());
            (d.head || d.body).appendChild(s);
        })();
    </script>
    <noscript>
        Please enable JavaScript to view the <a href="https://disqus.com/?ref_noscript">comments powered by Disqus.</a>
    </noscript>
{{- end }}
