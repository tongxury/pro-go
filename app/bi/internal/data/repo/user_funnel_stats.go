package repo

import (
	"context"
	"store/pkg/clients"
	"strings"
)

type UserFunnelRecordRepo struct {
	db *clients.ClickHouseClient
}

func NewUserFunnelRecordRepo(db *clients.ClickHouseClient) *UserFunnelRecordRepo {
	return &UserFunnelRecordRepo{db: db}
}

func (t *UserFunnelRecordRepo) Optimize(ctx context.Context) error {
	err := t.db.Exec(ctx, `optimize table user_funnel_records final`)
	if err != nil {
		return err
	}

	return nil
}

func (t *UserFunnelRecordRepo) GetUserFunnelStats(ctx context.Context, dateString string) ([]UserFunnelRecord, error) {

	date := strings.ReplaceAll(dateString, "-", "")

	// 官方访问
	var visits []UserFunnelRecord
	err := t.db.Select(ctx, &visits,
		`SELECT
    					'visit' as event,
    					1 as level,
    					channel, 
       					count(DISTINCT device_id) as count
			   FROM event_logs
			   WHERE event_id = 'home_page_view'
			     AND user_id = 0
			     AND toYYYYMMDD(create_time) = ?
			   group by channel
				`,
		date,
	)
	if err != nil {
		return nil, err
	}

	// 点击Add Chrome按钮
	var clicks []UserFunnelRecord
	err = t.db.Select(ctx, &clicks,
		`SELECT 
					'addToChrome' as event,
					2 as level,
    				channel, 
       				count(DISTINCT device_id) as count
                 FROM event_logs
                 WHERE event_id like 'click_add_chrome%'
                   AND user_id = 0
                   AND device_id IN
                       (SELECT device_id
                        FROM event_logs
                        WHERE event_id = 'home_page_view'
                          AND user_id = 0
                          AND toYYYYMMDD(create_time) = ?)
				group by channel
				`,
		date,
	)
	if err != nil {
		return nil, err
	}

	// 到达登录页面
	var installs []UserFunnelRecord
	err = t.db.Select(ctx, &installs,
		`SELECT 
    				'install' as event,
					3 as level,
    				channel, 
    				count(DISTINCT device_id) as count
                 FROM event_logs
                 WHERE event_id = 'onboarding_login_view'
                   AND referrer = 'plugin_installed'
                   AND user_id = 0
                   AND device_id IN
                       (SELECT device_id
                        FROM event_logs
                        WHERE event_id like 'click_add_chrome%'
                          AND device_id IN
                              (SELECT device_id
                               FROM event_logs
                               WHERE event_id = 'home_page_view'
                                 AND user_id = 0
                                 AND toYYYYMMDD(create_time) = ?))
				group by channel
				`,
		date,
	)
	if err != nil {
		return nil, err
	}

	// 到达登录页面
	var registers []UserFunnelRecord
	err = t.db.Select(ctx, &registers,
		`SELECT 
    				'register' as event,
					4 as level,
    				channel, 
    				count(DISTINCT device_id) as count
                 FROM event_logs
                 WHERE event_id = 'register'
                   AND device_id IN
                       (SELECT device_id
                        FROM event_logs
                        WHERE event_id = 'onboarding_login_view'
                          AND referrer = 'plugin_installed'
                          AND user_id = 0
                          AND device_id IN
                              (SELECT DISTINCT device_id
                               FROM event_logs
                               WHERE event_id like 'click_add_chrome%'
                                 AND device_id IN
                                     (SELECT DISTINCT device_id
                                      FROM event_logs
                                      WHERE event_id = 'home_page_view'
                                        AND user_id = 0
                                        AND toYYYYMMDD(create_time) = ?)))
				group by channel
				`,
		date,
	)
	if err != nil {
		return nil, err
	}

	// 订阅
	var subscribes []UserFunnelRecord
	err = t.db.Select(ctx, &subscribes,
		`SELECT 
    				'subscribe' as event,
					5 as level,
    				channel,
    				count(DISTINCT device_id) as count
                FROM event_logs
                WHERE event_id = 'subscribe'
				and device_id in
					  (SELECT device_id
					   FROM event_logs
					   WHERE event_id = 'register'
						 AND device_id IN
							 (SELECT device_id
							  FROM event_logs
							  WHERE event_id = 'onboarding_login_view'
								AND referrer = 'plugin_installed'
								AND user_id = 0
								AND device_id IN
									(SELECT DISTINCT device_id
									 FROM event_logs
									 WHERE event_id like 'click_add_chrome%'
									   AND device_id IN
										   (SELECT DISTINCT device_id
											FROM event_logs
											WHERE event_id = 'home_page_view'
											  AND user_id = 0
								  			 AND toYYYYMMDD(create_time) = ?)))) 
				group by channel
				`,
		date,
	)
	if err != nil {
		return nil, err
	}

	var records []UserFunnelRecord

	for _, xx := range [][]UserFunnelRecord{visits, clicks, installs, registers, subscribes} {
		for _, x := range xx {
			records = append(records, UserFunnelRecord{
				Date:    dateString,
				Channel: x.Channel,
				Event:   x.Event,
				Level:   x.Level,
				Count:   x.Count,
			})
		}
	}

	return records, nil
}

func (t *UserFunnelRecordRepo) Inserts(ctx context.Context, records []UserFunnelRecord) error {

	if len(records) == 0 {
		return nil
	}

	batch, err := t.db.PrepareBatch(ctx, `insert into user_funnel_records`)
	if err != nil {
		return err
	}

	for _, x := range records {
		err := batch.Append(
			x.Date, x.Channel, x.Event, x.Level, x.Count,
		)
		if err != nil {
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		return err
	}

	// 表不大
	t.Optimize(ctx)

	return nil
}

func (t *UserFunnelRecordRepo) FindByDate(ctx context.Context, date string) ([]UserFunnelRecord, error) {

	var records []UserFunnelRecord

	err := t.db.Select(ctx, &records, `
	select event, sum(count) as count
	from user_funnel_records
	where date = ?
	group by date, event, level
	order by level
	`, date)
	if err != nil {
		return nil, err
	}

	return records, nil
}
