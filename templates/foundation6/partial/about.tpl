{{ define "partial/about.tpl" }}
    <div class="callout primary">
        <div class="grid-x">
            <div class="cell small-12 medium-2 text-center">
                <img src="/static/{{ .theme }}/favicon.png?{{ .instance_id }}" alt="Vladimir Garvardt"
                     style="max-height: 150px; width: auto; border-radius: 50%;">
            </div>
            <div class="cell small-12 medium-8">
                <h3 class="text-center">rKlotz by Vladimir Garvardt</h3>
                <p>
                    Yet another simple single-user file-based golang-driven blog engine.
                    Source code available <a href="https://github.com/vgarvardt/rklotz" target="_blank">on Github</a>.
                    Usage example: <a href="https://itskrig.com/" target="_blank">itskrig.com</a>
                    (<a href="https://github.com/vgarvardt/itskrig.com" target="_blank">source code</a>).
                </p>
            </div>
            <div class="cell small-12 medium-2">
                <ul class="no-bullet">
                    <li>
                        <a href="https://github.com/vgarvardt/rklotz" class="btn btn-default btn-lg btn-block"
                           target="_blank">
                            <i class="fa-brands fa-github" aria-hidden="true"></i> <span
                                class="network-name">Github</span>
                        </a>
                    </li>
                    <li>
                        <a href="https://linkedin.com/in/vgarvardt/" class="btn btn-default btn-lg btn-block"
                           target="_blank">
                            <i class="fa-brands fa-linkedin"></i> <span class="network-name">LinkedIn</span>
                        </a>
                    </li>
                    <li>
                        <a href="https://mas.to/@vgarvardt" class="btn btn-default btn-lg btn-block"
                           target="_blank">
                            <i class="fa-brands fa-mastodon" aria-hidden="true"></i> <span
                                class="network-name">Mastodon</span>
                        </a>
                    </li>
                    <li>
                        <a href="https://twitter.com/vgarvardt" class="btn btn-default btn-lg btn-block"
                           target="_blank">
                            <i class="fa-brands fa-twitter" aria-hidden="true"></i> <span
                                class="network-name">Twitter</span>
                        </a>
                    </li>
                    <li>
                        <a href="https://itskrig.com" class="btn btn-default btn-lg btn-block" target="_blank">
                            <i class="fa-solid fa-globe" aria-hidden="true"></i> <span
                                class="network-name">Website</span>
                        </a>
                    </li>
                </ul>
            </div>
        </div>
    </div>
{{ end }}
