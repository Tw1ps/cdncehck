package check

import (
	"encoding/json"
	io "io/ioutil"
	"net"
	"strings"

	"github.com/yl2chen/cidranger"
)

type CdnCache struct {
	Name  string   `json:"name"`
	Cname []string `json:"cname"`
	Cidr  []string `json:"cidr"`
}

type CdnClient struct {
	CidRanges   map[string][]string
	CnameRanges map[string][]string
	rangers     map[string]cidranger.Ranger
}

type Result struct {
	Provider string `json:"provider"`
	Match    string `json:"match"`
	Found    bool   `json:"found"`
}

func InitCdnClient(filepath string) (*CdnClient, error) {
	var cli CdnClient
	cli.CidRanges = make(map[string][]string)
	cli.CnameRanges = make(map[string][]string)
	cli.rangers = make(map[string]cidranger.Ranger)

	cache, err := loadJson(filepath)
	if err != nil {
		return nil, err
	}

	for _, cc := range cache {
		cli.CnameRanges[cc.Name] = append(cli.CnameRanges[cc.Name], cc.Cname...)
		cli.CidRanges[cc.Name] = append(cli.CidRanges[cc.Name], cc.Cidr...)
	}

	for k, v := range cli.CidRanges {
		ranger := cidranger.NewPCTrieRanger()

		for _, cidr := range v {
			_, network, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			_ = ranger.Insert(cidranger.NewBasicRangerEntry(*network))
		}
		cli.rangers[k] = ranger
	}

	return &cli, nil
}

func loadJson(filepath string) ([]CdnCache, error) {
	var ret []CdnCache
	data, err := io.ReadFile(filepath)
	if err != nil {
		return ret, err
	}
	dataJson := []byte(data)
	err = json.Unmarshal(dataJson, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (cli *CdnClient) Check(target string) Result {
	var result Result
	ip := net.ParseIP(target)
	i, _, err := net.ParseCIDR(target)
	if err == nil {
		ip = i
	}
	if ip == nil {
		for provider, cnames := range cli.CnameRanges {
			for _, v := range cnames {
				flag := strings.Contains(target, v)
				if flag {
					result.Provider = provider
					result.Match = v
					result.Found = true
					return result
				}
			}
		}
	} else {
		for provider, ranger := range cli.rangers {
			contains, _ := ranger.ContainingNetworks(ip)
			if len(contains) > 0 {
				nw := contains[0].Network()
				result.Provider = provider
				result.Match = nw.String()
				result.Found = true
				return result
			}
		}
	}

	return result
}

func (cli *CdnClient) RangeCheck(targets []string) Result {
	var result Result
	for _, v := range targets {
		result = cli.Check(v)
		if result.Found {
			return result
		}
	}
	return result
}
