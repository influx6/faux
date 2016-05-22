package middleware

import (
	"github.com/ardanlabs/kit/log"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/web/app"
)

// LogExcludeKey defines a key if found within the incoming context, excludes
// that call from getting logged.
const LogExcludeKey = "ExcludeFromLog"

// Log defines a log middleware which helps us log the incoming requests comming
// into the app.
func Log(h app.Handler) app.Handler {
	return func(ctx context.Context, w *app.ResponseRequest) error {

		// Do we excluse this call from the logs?
		_, ok := ctx.Get(LogExcludeKey)
		if ok {
			return h(ctx, w)
		}

		log.Dev("middleware.Log", "Log", "Started : Method[%s] : From[%s] : Path[%s] : Query[%s]: Server[%s]", w.R.Method, w.R.URL.Host, w.R.URL.Path, w.R.URL.RawQuery, w.R.RemoteAddr)

		err := h(ctx, w)
		log.Dev("middleware.Log", "Log", "Info : Status[%s] : Method[%s] : From[%s] : Path[%s] : Server[%s]", w.Status(), w.R.Method, w.R.URL.Host, w.R.URL.Path, w.R.RemoteAddr)
		if err != nil {
			log.Error("middleware.Log", "Log", err, "Completed")
			return err
		}

		log.Dev("middleware.Log", "Log", "Completed")
		return nil
	}
}
