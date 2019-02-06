package appregistry

import (
	"net/url"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apprclient "github.com/operator-framework/go-appr/appregistry"
)

// NewClientFactory return a factory which can be used to instantiate a new appregistry client
func NewClientFactory() ClientFactory {
	return &factory{}
}

// ClientFactory is an interface that wraps the New method.
//
// New returns a new instance of appregistry Client from the specified source.
type ClientFactory interface {
	New(source string) (Client, error)
}

type factory struct{}

func (f *factory) New(source string) (Client, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	transport := httptransport.New(u.Host, u.Path, []string{u.Scheme})
	transport.Consumers["application/x-gzip"] = runtime.ByteStreamConsumer()
	c := apprclient.New(transport, strfmt.Default)

	return &client{
		adapter: &apprApiAdapterImpl{client: c},
		decoder: &blobDecoderImpl{},
	}, nil
}
