package exec

import (
	"context"
	"encoding/json"
	"reflect"
	"bytes"
	//"github.com/davecgh/go-spew/spew"

	"github.com/neelance/graphql-go/errors"
	"github.com/neelance/graphql-go/internal/exec/resolvable"
	"github.com/neelance/graphql-go/internal/exec/selected"
	"github.com/neelance/graphql-go/internal/query"
)

func (r *Request) Subscribe(ctx context.Context, s *resolvable.Schema, op *query.Operation) (<-chan json.RawMessage, <-chan []*errors.QueryError) {

	// consider to set capacity the same as resolver. result.Cap()
	ch := make(chan json.RawMessage, 1)
	chErr := make(chan []*errors.QueryError,1)
	var result reflect.Value

	// extract result should be chan
	var f *fieldToExec
	func() {
		defer r.handlePanic(ctx)

		sels := selected.ApplyOperation(&r.Request, s, op)
		var fields []*fieldToExec
		collectFieldsToResolve(sels, s.Resolver, &fields, make(map[string]*fieldToExec))
		// TODO: more subscriptions at once
		f = fields[0]

		var in []reflect.Value
		if f.field.HasContext {
			in = append(in, reflect.ValueOf(ctx))
		}
		if f.field.ArgsPacker != nil {
			in = append(in, f.field.PackedArgs)
		}
		callOut := f.resolver.Method(f.field.MethodIndex).Call(in)
		result = callOut[0]
	}()

	// TODO: check error callOut[1]
	if err := ctx.Err(); err != nil {
		chErr <- []*errors.QueryError{errors.Errorf("%s", err)}
		close(chErr)
		close(ch)
		return ch, chErr
	}
	// TODO: check if result chan is nil
 

	// XXX: new context??
	go func() {
		for {
			wasClosed := false
			func() {
				defer r.handlePanic(ctx)
				obj, ok := result.Recv()
				if !ok {
					close(ch)
					wasClosed = true
					return
				}
				var out bytes.Buffer
				r.execSelectionSet(ctx, f.sels, f.field.Type, &pathSegment{nil,f.field.Alias}, obj, &out)
				ch <- json.RawMessage(out.Bytes())
			}()
			if err := ctx.Err(); err != nil {
				chErr <- []*errors.QueryError{errors.Errorf("%s", err)}
			}
			if wasClosed {
				close(chErr)
				return
			}
		}
	}()

	return ch, chErr

}
