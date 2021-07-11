package providers

import (
	"context"
	"fmt"

	"code.gitea.io/sdk/gitea"
)

type GiteaCredential struct {
	URL   string
	Token string
	Debug bool
}

type GiteaClient struct {
	delegate  *gitea.Client
	namespace string
	name      string
	commitRef string
}

func NewGiteaClient(namespace string, name string, commitRef string, cred *GiteaCredential) (*GiteaClient, error) {
	client, err := gitea.NewClient(cred.URL, gitea.SetToken(cred.Token))
	if err != nil {
		return nil, err
	}
	if cred.Debug {
		gitea.SetDebugMode()(client)
	}

	giteaClient := &GiteaClient{
		delegate:  client,
		namespace: namespace,
		name:      name,
		commitRef: commitRef,
	}
	return giteaClient, nil
}

func (c *GiteaClient) GetFileListing(ctx context.Context, path string) ([]FileListingEntry, error) {
	c.delegate.SetContext(ctx)

	contents, _, err := c.delegate.ListContents(c.namespace, c.name, c.commitRef, path)
	if err != nil {
		return nil, err
	}

	fileListing := make([]FileListingEntry, 0, len(contents))
	for _, content := range contents {
		entry := FileListingEntry{Type: content.Type, Name: content.Name, Path: content.Path}
		fileListing = append(fileListing, entry)
	}
	return fileListing, nil
}

func (c *GiteaClient) GetFileContent(ctx context.Context, path string) (fileContent string, err error) {
	c.delegate.SetContext(ctx)

	data, _, err := c.delegate.GetFile(c.namespace, c.name, c.commitRef, path)
	return fmt.Sprintf("%s", data), err
}