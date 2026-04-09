package grpc

import "context"

// HistoryRecorder wraps a MethodInvoker and automatically records each
// invocation into a History store.
type HistoryRecorder struct {
	invoker *MethodInvoker
	history *History
	address string
}

// NewHistoryRecorder creates a HistoryRecorder that decorates invoker and
// persists results to h.
func NewHistoryRecorder(invoker *MethodInvoker, h *History, address string) *HistoryRecorder {
	return &HistoryRecorder{
		invoker: invoker,
		history: h,
		address: address,
	}
}

// InvokeAndRecord calls the underlying MethodInvoker and stores the result.
// It returns the response string and any error, mirroring the invoker API.
func (r *HistoryRecorder) InvokeAndRecord(
	ctx context.Context,
	service, method, requestJSON string,
) (string, error) {
	resp, err := r.invoker.InvokeUnary(ctx, nil, requestJSON)

	entry := HistoryEntry{
		Address:  r.address,
		Service:  service,
		Method:   method,
		Request:  requestJSON,
		Response: resp,
	}
	if err != nil {
		entry.Error = err.Error()
	}
	r.history.Add(entry)
	return resp, err
}
