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

type Icon struct {
	Name string
	Key  string
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

func (c Client) List(prefix, marker string, search bool) ([]Icon, error) {
	// prefix
	// ""
	// svgs/brands/
	// svgs/regular/
	// svgs/solid/
	// svgs/brands/apple
	var opt *cos.BucketGetOptions
	if !search {
		if marker == "" {
			opt = &cos.BucketGetOptions{
				Prefix:  prefix,
				MaxKeys: 24,
			}
		} else {
			opt = &cos.BucketGetOptions{
				Prefix:  prefix,
				MaxKeys: 24,
				Marker:  "svgs/" + marker + ".svg",
			}
		}
	} else {
		opt = &cos.BucketGetOptions{
			Prefix: prefix,
		}
	}

	r, _, err := c.client.Bucket.Get(context.Background(), opt)
	if err != nil {
		return nil, err
	}
	var iconList []Icon
	for _, o := range r.Contents {
		if o.Size != 0 {
			iconList = append(iconList, Icon{
				Name: utils.ParseName(o.Key),
				Key:  utils.ParsePrefix(o.Key),
			})
		}
	}
	if len(iconList) == 0 {
		return []Icon{}, nil
	}
	return iconList, nil
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
