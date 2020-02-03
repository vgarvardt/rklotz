package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	f := New()

	_, err := f("", "unknown")
	require.Error(t, err)
	assert.Equal(t, ErrorUnknownFormat, err)

	formattedFromMD, err := f(rawMD, "md")
	require.NoError(t, err)
	assert.Equal(t, html, formattedFromMD)
}

const (
	rawMD = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eget risus in lorem [convallis semper](https://example.com).
Sed posuere vehicula feugiat. Maecenas facilisis nunc nisl, sit amet ornare quam scelerisque vel.
Vestibulum non nunc justo. Donec vitae justo ipsum. Cras tempor nec tortor vitae suscipit.
In vulputate lorem id quam tincidunt, non pulvinar dui varius. Sed a imperdiet orci.
Aliquam et sem in tellus dapibus lobortis. Quisque auctor laoreet massa, in tincidunt lectus rutrum vitae.

` + "```go" +
		`
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello Markdown")
}

` + "```" +
		`
Vestibulum hendrerit massa libero, et sagittis felis luctus ut. Nunc condimentum aliquet lectus,
id posuere risus rhoncus et. Vivamus sed diam aliquam, gravida neque ut, luctus purus.
Mauris fringilla sagittis pretium. In egestas urna lectus, semper vehicula libero eleifend vitae.
Duis vitae dolor sit amet purus eleifend venenatis in vitae ligula. In quis est libero.
`
	html = `<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eget risus in lorem <a href="https://example.com" target="_blank">convallis semper</a>.
Sed posuere vehicula feugiat. Maecenas facilisis nunc nisl, sit amet ornare quam scelerisque vel.
Vestibulum non nunc justo. Donec vitae justo ipsum. Cras tempor nec tortor vitae suscipit.
In vulputate lorem id quam tincidunt, non pulvinar dui varius. Sed a imperdiet orci.
Aliquam et sem in tellus dapibus lobortis. Quisque auctor laoreet massa, in tincidunt lectus rutrum vitae.</p>

<pre><code class="go">package main

import (
	&quot;fmt&quot;
)

func main() {
	fmt.Println(&quot;Hello Markdown&quot;)
}

</code></pre>

<p>Vestibulum hendrerit massa libero, et sagittis felis luctus ut. Nunc condimentum aliquet lectus,
id posuere risus rhoncus et. Vivamus sed diam aliquam, gravida neque ut, luctus purus.
Mauris fringilla sagittis pretium. In egestas urna lectus, semper vehicula libero eleifend vitae.
Duis vitae dolor sit amet purus eleifend venenatis in vitae ligula. In quis est libero.</p>
`
)
