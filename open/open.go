// Open a file, directory, or URI using the OS's default
// application for that object type.  Optionally, you can
// specify an application to use.
//
// This is a proxy for the following commands:
//
//	        OSX: "open"
//	    Windows: "start"
//	Linux/Other: "xdg-open"
//
// This is a golang port of the node.js module: https://github.com/pwnall/node-open
package open

import (
	"context"

	"github.com/skratchdot/open-golang/opts"
)

// Opener opens a file, directory, or URI using the OS's default application,
// optionally with the supplied options.
type Opener interface {
	Open(ctx context.Context, target string, options ...opts.Option) error
}

// OpenerFunc is an adapter that allows ordinary functions to be used as
// [Opener] values.
type OpenerFunc func(ctx context.Context, target string, options ...opts.Option) error

// Open calls fn with ctx, target, and options.
func (fn OpenerFunc) Open(ctx context.Context, target string, options ...opts.Option) error {
	return fn(ctx, target, options...)
}

// Interface compile-time compliance check
var _ Opener = OpenerFunc(Open)

// Open will open a file, directory, or URI using the OS's default
// application for that obect type or URI mime.
func Open(ctx context.Context, target string, options ...opts.Option) error {
	cfg, cmd := createCmd(ctx, target, options...)
	if cfg.Wait() {
		return cmd.Run()
	} else {
		return cmd.Start()
	}
}

// Deprecated: use [Open] with [Wait]
//
// Open a file, directory, or URI using the OS's default
// application for that object type. Wait for the open
// command to complete.
//
//go:fix inline
func Run(target string, options ...opts.Option) error {
	return Open(context.TODO(), target, opts.Wait())
}

// Deprecated: use [Open]
//
// Open a file, directory, or URI using the OS's default
// application for that object type. Don't wait for the
// open command to complete.
//
//go:fix inline
func Start(target string, options ...opts.Option) error {
	return Open(context.TODO(), target)
}

// Deprecated: use [Open] with [Wait] and [WithApp]
//
// Open a file, directory, or URI using the specified application.
// Wait for the open command to complete.
//
//go:fix inline
func RunWith(target string, appName string) error {
	return Open(context.TODO(), target, opts.WithApp(appName), opts.Wait())
}

// Deprecated: use [Open] with [WithApp]
//
// Open a file, directory, or URI using the specified application.
// Don't wait for the open command to complete.
//
//go:fix inline
func StartWith(target string, appName string) error {
	return Open(context.TODO(), target, opts.WithApp(appName))
}
