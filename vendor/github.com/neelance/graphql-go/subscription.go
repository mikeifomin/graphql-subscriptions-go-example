package graphql

import (
	"context"
	//"encoding/json"
	"sync"

	"github.com/neelance/graphql-go/errors"
	"github.com/neelance/graphql-go/internal/common"
	"github.com/neelance/graphql-go/internal/exec/resolvable"
	"github.com/neelance/graphql-go/internal/exec/selected"
	"github.com/neelance/graphql-go/internal/query"

	"github.com/neelance/graphql-go/internal/validation"
	"github.com/neelance/graphql-go/internal/exec"
	"github.com/neelance/graphql-go/introspection"
)

func (s *Schema) Subscribe(ctx context.Context, queryString string, operationName string, variables map[string]interface{}) (<-chan *Response) {
		return s.subscribe(ctx,queryString,operationName,variables,s.res)
}

func closeErrs(errs []*errors.QueryError) <-chan *Response {
	ch := make(chan *Response)
	ch <- &Response{Errors: errs}
	close(ch)
	return ch
}
func closeOneErr(err *errors.QueryError) <-chan *Response {
	return closeErrs([]*errors.QueryError{err})
}

func (s *Schema) subscribe(ctx context.Context, queryString string, operationName string, variables map[string]interface{}, res *resolvable.Schema) <-chan *Response {
	doc, qErr := query.Parse(queryString)
	if qErr != nil {
		return closeOneErr(qErr)
	}

	errs := validation.Validate(s.schema, doc)
	if len(errs) != 0 {
		return closeErrs(errs)
	}

	op, err := getOperation(doc, operationName)
	if err != nil {
		return closeOneErr(errors.Errorf("%s",err))
	}

	r := &exec.Request{
		Request: selected.Request{
			Doc:    doc,
			Vars:   variables,
			Schema: s.schema,
		},
		Limiter: make(chan struct{}, s.maxParallelism),
		Tracer:  s.tracer,
		Logger:  s.logger,
	}

	varTypes := make(map[string]*introspection.Type)

	for _, v := range op.Vars {
		t, err := common.ResolveType(v.Type, s.schema.Resolve)
		if err != nil {
			return closeOneErr(err)
		}
		varTypes[v.Name.Name] = introspection.WrapType(t)
	}
	traceCtx, finish := s.tracer.TraceQuery(ctx, queryString, operationName, variables, varTypes)

	out := make(chan *Response)
	if op.Type != query.Subscription {
		b, errs := r.Execute(traceCtx, res, op)
	  out <- &Response{	
			Data:   b,
			Errors: errs,
		}
		close(out)
		return out
	}

	chData, chErr := r.Subscribe(traceCtx, res, op)

  // fan-in
	var wg sync.WaitGroup
	wg.Add(2)
	go func(){
		for d := range chData {
       out <- &Response{
				 Data: d,
			 }
		}
		wg.Done()
	}()
	go func(){
		for e := range chErr {
			out <- &Response{
				Errors: e,
			}
		}
		wg.Done()
	}()
	go func(){
		wg.Wait()
		close(out)
		finish([]*errors.QueryError{})
	}()
	return out
}

