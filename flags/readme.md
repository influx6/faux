Flags
--------
Flags provides a simple multicommand package built ontop of the internal `flag` Go/Golang package.


See [Flags Docs](https://golang.org/pkg/github.com/influx6/faux/flags) for more.

If only one command is provided then Flags treats that command like the default command to be runned when argument is giving with no command name.

## Example


```go
import (
	"errors"

	"github.com/influx6/faux/flags"
)

func main() {
	flags.Run("sitecrawler", flags.Command{
		Name:      "crawl",
		ShortDesc: "Crawls provided website URL returning json sitemap.",
		Desc:      "Crawl is the entry command to crawl a website, it runs through all pages of giving host, ignoring externals links. It prints status and link connection as json on a per link basis.",
		Usages:    []string{"sitecrawler crawl https://monzo.com"},
		Flags: []flags.Flag{
			&flags.IntFlag{
				Name:    "depth",
				Default: -1,
				Desc:    "Sets the depth to crawl through giving site",
			},
		},
		Action: func(ctx flags.Context) error {
			if len(ctx.Args()) == 0 {
				return errors.New("Must provide website url for crawling. See examples section")
			}

			return nil
		},
	})
}

```