package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	"store/pkg/clients"
)

func main() {

	//config, err := confcenter.GetConfig[conf.BizConfig]()
	//if err != nil {
	//	log.Fatalf("failed get config: %v", err)
	//}

	client := clients.NewLokiClient("http://kuboard-loki.kuboard.svc.cluster.local:3100")

	//lcClient := cache.New(10*time.Minute, 10*time.Minute)

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("@every 1s", func() {
		result, err := client.Query("sum by(container) (rate({namespace=\"staging\", container=\"sg-llm\"} |~ `(?i)error` [1m]))") //now.Add(-time.Minute), now,
		if err != nil {
			log.Debugw("Query err", err)
			return
		}

		//target := float64(0.8)

		log.Debugw("Query result", result)

		if len(result) > 0 {
			if len(result[0].Value) > 1 {
				log.Debugw("Query llm value", result[0].Value[1])
				//if conv.Float64(result[0].Value[1]) > target {
				//
				//	if _, found := lcClient.Get("lock"); found {
				//		return
				//	}
				//
				//	fc := feishuo.NewClient("https://open.feishu.cn/open-apis/bot/v2/hook/d05367c5-3a8a-44a7-b2b5-5b421a6f805a")
				//
				//	params := feishuo.RichText{
				//		Title: "监控报警",
				//		Content: []feishuo.RichTextParagraph{
				//			[]feishuo.RichTextContent{
				//				{
				//					Tag:  "text",
				//					Text: fmt.Sprintf("【%s】服务错误日志每秒频次超过 %v 次", "sg-llm", target),
				//				},
				//				{
				//					Tag:  "a",
				//					Text: "请尽快查看",
				//					Href: "http://18.215.179.67:18060/k8s-proxy/studygpt/api/v1/namespaces/kuboard/services/http:kuboard-loki-grafana:3000/proxy/explore?orgId=1&kiosk=tv&left=%7B%22datasource%22:%221_42aydIz%22,%22queries%22:%5B%7B%22expr%22:%22%7Bnamespace%3D%5C%22prod%5C%22,%20container%3D%5C%22sg-llm%5C%22%7D%20%7C~%20%60%28%3Fi%29error%60%22,%22refId%22:%22A%22,%22editorMode%22:%22builder%22,%22queryType%22:%22range%22%7D%5D,%22range%22:%7B%22from%22:%22now-5m%22,%22to%22:%22now%22%7D%7D\n",
				//				},
				//			},
				//		},
				//	}
				//
				//	log.Debugw("Sending ", params)
				//
				//	err := fc.SendRichText(context.Background(), params)
				//	if err != nil {
				//		log.Errorw("SendRichText err", err)
				//		return
				//	}
				//
				//	lcClient.Set("lock", "1", 5*time.Minute)
				//}
			}
		}

	})
	if err != nil {
		log.Fatalw("AddFunc err", err)
	}

	_, err = c.AddFunc("@every 1s", func() {
		result, err := client.Query("sum by(container) (rate({namespace=\"staging\", container=\"sg-user\"} |~ `(?i)error` [1m]))") //now.Add(-time.Minute), now,
		if err != nil {
			log.Debugw("Query err", err)
			return
		}

		if len(result) > 0 {
			if len(result[0].Value) > 1 {
				log.Debugw("Query staging user", result[0].Value[1])
			}
		}

	})
	if err != nil {
		log.Fatalw("AddFunc err", err)
	}

	c.Run()
}

//http://18.215.179.67:18060/k8s-proxy/studygpt/api/v1/namespaces/kuboard/services/http:kuboard-loki-grafana:3000/proxy/explore?orgId=1&kiosk=tv&left=%7B%22datasource%22:%221_42aydIz%22,%22queries%22:%5B%7B%22expr%22:%22%7Bnamespace%3D%5C%22prod%5C%22,%20container%3D%5C%22sg-llm%5C%22%7D%20%7C~%20%60%28%3Fi%29error%60%22,%22refId%22:%22A%22,%22editorMode%22:%22builder%22,%22queryType%22:%22range%22%7D%5D,%22range%22:%7B%22from%22:%22now-5m%22,%22to%22:%22now%22%7D%7D
