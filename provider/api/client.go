package api

import (
	"github.com/Khan/genqlient/graphql"
	"net/http"
)


func NewClient(url string, token string,version string) graphql.Client{

	

	 
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	httpClient := &http.Client{
		Transport: &transportWithAuth{
			Transport: transport,
			Token:    token,
			Version:  version,
		},
	}


	 
	client := graphql.NewClient(url, httpClient)
	return client


}







type transportWithAuth struct {
	Transport http.RoundTripper
	Token    string
	Version   string
}

func (t *transportWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization",t.Token)
	req.Header.Set("x-terraform-provider-version",t.Version)
	return t.Transport.RoundTrip(req)
}
