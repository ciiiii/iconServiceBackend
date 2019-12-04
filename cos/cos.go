package cos

import (
	"net/url"
	"net/http"
	"github.com/ciiiii/iconServiceBackend/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"fmt"
	"strings"
	"context"
	"github.com/ciiiii/iconServiceBackend/utils"
)

type Client struct {
	client *cos.Client
	prefix string
}

type Object struct {
	Name string
	Key  string
	Url  string
	Date string
	Size int
}

func Init() *Client {
	prefix := fmt.Sprintf("http://%s.cos.%s.myqcloud.com", config.Parser().Cos.BucketName, config.Parser().Cos.Region)
	u, _ := url.Parse(prefix)
	b := &cos.BaseURL{BucketURL: u}
	c := Client{cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Parser().Cos.SecretID,
			SecretKey: config.Parser().Cos.SecretKey,
		},
	}), prefix}
	return &c
}

func (c Client) List(prefix string) ([]Object, error) {
	opt := &cos.BucketGetOptions{
		Prefix: prefix,
	}
	v, _, err := c.client.Bucket.Get(context.Background(), opt)
	if err != nil {
		return nil, err
	}
	var objectList []Object
	for _, o := range v.Contents {
		if o.Size != 0 {
			objectList = append(objectList, Object{
				Name: utils.ParseKey(o.Key),
				Key:  o.Key,
				Url:  strings.Join([]string{c.prefix, o.Key}, "/"),
				Date: o.LastModified,
				Size: o.Size,
			})
		}
	}
	return objectList, nil
}

func (c Client) Upload(path string) error {
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "image/svg+xml",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}
	splitPath := strings.Split(path, "web/")
	key := splitPath[len(splitPath)-1]
	fmt.Println(key)
	_, err := c.client.Object.PutFromFile(context.Background(), key, path, opt)
	return err
}
