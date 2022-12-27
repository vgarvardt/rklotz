{{ define "partial/pagination.tpl" }}
    <nav aria-label="Pagination">
        <ul class="pagination text-center">
            {{ if ne .page 0 }}
                <li class="pagination-previous"><a href="/?page={{ .page | add -1 }}">Newer</a></li>
            {{ else }}
                <li class="pagination-previous disabled">Newer</li>
            {{ end }}
            <li>Page {{ .page | add 1 }} of {{ .meta.Pages }}</li>
            {{ if lt (.page | add 1) $.meta.Pages }}
                <li class="pagination-next"><a href="/?page={{ .page | add 1 }}">Older</a></li>
            {{ else }}
                <li class="pagination-next disabled">Older</li>
            {{ end }}
        </ul>
    </nav>
{{ end }}
