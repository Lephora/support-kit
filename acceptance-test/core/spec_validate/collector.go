package spec_validate

import (
	"context"
	"crypto/tls"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

func GithubQuery() *githubv4.Client {
	token := os.Getenv("GITHUB_TOKEN")
	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			Base: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Source: oauth2.ReuseTokenSource(nil, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})),
		},
	}
	return githubv4.NewClient(httpClient)

}

func QuerySpec(client *githubv4.Client, repoOwner, repoName, expression string) (string, error) {
	var q struct {
		Repository struct {
			Object struct {
				Blob struct {
					Text string
				} `graphql:"... on Blob"`
			} `graphql:"object(expression:$expression)"`
		} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
	}
	variables := map[string]interface{}{
		"repoOwner":  githubv4.String(repoOwner),
		"repoName":   githubv4.String(repoName),
		"expression": githubv4.String(expression),
	}
	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		return "", nil
	}
	return q.Repository.Object.Blob.Text, nil
}

var FlowPool []*struct {
	Req  *http.Request
	Resp *http.Response
}

func CollectFlow(req *http.Request, resp *http.Response) {
	req.URL.Host = "lephora"
	FlowPool = append(FlowPool, &struct {
		Req  *http.Request
		Resp *http.Response
	}{
		Req:  req,
		Resp: resp,
	})
}
