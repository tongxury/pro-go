package repo

import (
	"context"
	"github.com/patrickmn/go-cache"
	"store/pkg/clients"
)

type IpDetailRepo struct {
	db *clients.ClickHouseClient
	lc *cache.Cache
}

func NewIpDetailRepo(db *clients.ClickHouseClient, lc *cache.Cache) *IpDetailRepo {
	return &IpDetailRepo{db: db, lc: lc}
}

func (t *IpDetailRepo) FindByIp(ctx context.Context, ip string) (*IpDetail, error) {

	if ip == "" {
		return nil, nil
	}

	cacheKey := "ip_detail.cache:" + ip

	cacheValue, found := t.lc.Get(cacheKey)
	if found {
		return cacheValue.(*IpDetail), nil
	}

	var details IpDetails

	err := t.db.Select(ctx, &details,
		`select  
					ip_range_start,
					ip_range_end,
					ip_range_start_num,
					ip_range_end_num,
					country_code,
					state1,
					state2,
					city,
					postcode,
					latitude,
					longitude,
					timezone
    			from ip_details 
    			where IPv4StringToNum(?) between ip_range_start_num and ip_range_end_num
    			`,
		ip,
	)

	if err != nil {
		return nil, err
	}

	if len(details) == 0 {
		return nil, nil
	}

	t.lc.Set(cacheKey, &details[0], -1)

	return &details[0], nil
}
