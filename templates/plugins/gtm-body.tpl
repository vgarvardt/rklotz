{{ define "plugins/gtm-body.tpl" -}}
    <noscript>
        <iframe src="https://www.googletagmanager.com/ns.html?id={{ .plugin.gtm.id}}"
                height="0" width="0" style="display:none;visibility:hidden"></iframe>
    </noscript>
{{- end }}
