package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
	"store/pkg/clients"
	"store/pkg/sdk/helper"
	"time"
)

type UserDeviceRepo struct {
	db *clients.ClickHouseClient
	lc *cache.Cache
}

func NewUserDeviceRepo(db *clients.ClickHouseClient, lc *cache.Cache) *UserDeviceRepo {
	return &UserDeviceRepo{db: db, lc: lc}
}

func (t *UserDeviceRepo) GetDeviceIDByUserID(ctx context.Context, userID int64) (string, error) {

	ds, err := t.FindByUserIDs(ctx, []int64{userID})
	if err != nil {
		return "00000000-0000-0000-0000-000000000000", err
	}

	if len(ds) > 0 {
		return ds[0].DeviceID, nil
	}
	return "00000000-0000-0000-0000-000000000000", nil
}

func (t *UserDeviceRepo) FindByUserIDs(ctx context.Context, userIDs []int64) (UserDevices, error) {

	var devices UserDevices

	err := t.db.Select(ctx, &devices,
		`select  
					device_id,
					ip,
					app_version,
					created_at,
					chrome_version,
					platform,
					user_id,
					channel
    			from user_devices 
    			where user_id in (?)
    			`,
		userIDs,
	)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (t *UserDeviceRepo) FindByIDs(ctx context.Context, ids []string) (UserDevices, error) {

	if len(ids) == 0 {
		return nil, nil
	}

	var devices UserDevices

	err := t.db.Select(ctx, &devices,
		`select  
					device_id,
					ip,
					app_version,
					created_at,
					chrome_version,
					platform,
					user_id,
					channel
    			from user_devices 
    			where device_id in (?)
    			`,
		ids,
	)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (t *UserDeviceRepo) FindIpDetail(ctx context.Context, ip string) (*IpDetail, error) {

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

func (t *UserDeviceRepo) FindByID(ctx context.Context, id string) (*UserDevice, error) {

	cacheKey := "device.cache:" + id

	cacheValue, found := t.lc.Get(cacheKey)
	if found {
		return cacheValue.(*UserDevice), nil
	}

	var oldDevices UserDevices

	err := t.db.Select(ctx, &oldDevices,
		`select  
					device_id,
					ip,
					app_version,
					created_at,
					chrome_version,
					platform,
					user_id,
					channel
    			from user_devices 
    			where device_id = ? 
    			order by created_at 
    			desc limit 1`,
		id,
	)
	if err != nil {
		log.Errorw("Select err", err)
		return nil, err
	}

	if len(oldDevices) == 0 {
		return nil, nil
	}

	t.lc.Set(cacheKey, &oldDevices[0], -1)

	return &oldDevices[0], nil
}

func (t *UserDeviceRepo) Upgrade(ctx context.Context, d *UserDevice) error {

	var oldDevices UserDevices

	err := t.db.Select(ctx, &oldDevices,
		`select  
					device_id,
					ip,
					app_version,
					created_at,
					chrome_version,
					platform,
					user_id,
					channel
    			from user_devices 
    			where device_id = ? 
    			order by created_at 
    			desc limit 1`,
		d.DeviceID,
	)
	if err != nil {
		log.Errorw("Select err", err)
	}

	ip := d.Ip
	appVersion := d.AppVersion
	chromeVersion := d.ChromeVersion
	platform := d.Platform
	userID := d.UserID
	channel := d.Channel
	createdAt := d.CreatedAt
	country := d.Country
	if len(oldDevices) > 0 {

		if ip != "" {
			detail, err := t.FindIpDetail(ctx, ip)
			if err != nil {
				log.Errorw("FindIpDetail err")
			}

			if detail != nil {
				country = detail.CountryCode
			} else {
				country = ""
			}
		} else {
			ip = oldDevices[0].Ip
		}

		if ip != "" && ip != oldDevices[0].Ip {

		}

		//ip = helper.OrString(ip, oldDevices[0].Ip)
		appVersion = helper.OrString(appVersion, oldDevices[0].AppVersion)
		chromeVersion = helper.OrString(chromeVersion, oldDevices[0].ChromeVersion)
		platform = helper.OrString(platform, oldDevices[0].Platform)
		userID = helper.OrInt64(userID, oldDevices[0].UserID)
		channel = helper.OrString(channel, oldDevices[0].Channel)
		createdAt = oldDevices[0].CreatedAt
	}

	err = t.Insert(ctx, &UserDevice{
		DeviceID:      d.DeviceID,
		Ip:            ip,
		AppVersion:    appVersion,
		CreatedAt:     createdAt,
		ChromeVersion: chromeVersion,
		Platform:      platform,
		UserID:        userID,
		Channel:       channel,
		Country:       country,
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *UserDeviceRepo) Insert(ctx context.Context, d *UserDevice) error {

	err := t.db.Exec(ctx, `
				insert into user_devices (
					device_id,
					ip,
					app_version,
					created_at,
					chrome_version,
					platform,
					user_id,
					channel,
					country
				        ) 
				values (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.DeviceID, d.Ip, d.AppVersion, d.CreatedAt.Format(time.DateTime), d.ChromeVersion,
		d.Platform, d.UserID, d.Channel, d.Country,
	)
	if err != nil {
		return err
	}

	return nil
}
