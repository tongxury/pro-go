package douyin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
)

type CommodityMetadata struct {
	Title  string
	Brand  string
	Images []string
}

func (t *Client) GetCommodityMetadataByAPI(ctx context.Context, url string) (*CommodityMetadata, error) {

	res, err := resty.New().R().
		SetContext(ctx).
		Get("http://118.196.62.225:9527/douyin/commodity?url=" + url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, errors.New(res.Status())
	}

	var data CommodityMetadata
	err = json.Unmarshal(res.Body(), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (t *Client) GetCommodityMetadata(ctx context.Context, url string) (*CommodityMetadata, error) {
	// 随机化浏览器参数
	userAgents := []string{
		//"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		//"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		//"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",

		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chromium/120.0.6099.109 Safari/537.36",
	}

	windowSizes := []string{
		"1920,1080",
		"1366,768",
		"1440,900",
		"1536,864",
		"1280,720",
		"1600,900",
	}

	rand.Seed(time.Now().UnixNano())
	userAgent := userAgents[rand.Intn(len(userAgents))]
	windowSize := windowSizes[rand.Intn(len(windowSizes))]

	// 创建allocator context，使用更隐蔽的配置
	allocCtx, cancel := chromedp.NewExecAllocator(ctx,
		chromedp.Flag("headless", true), // 显示浏览器窗口
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-field-trial-config", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("window-size", windowSize),
		chromedp.UserAgent(userAgent),
		// 移除明显的反检测标志
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		// 添加更多真实的浏览器参数
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-field-trial-config", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("disable-background-networking", true),
	)
	defer cancel()

	// 创建chromedp context
	chromedpCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 设置随机超时时间
	timeout := time.Duration(30+rand.Intn(30)) * time.Second
	chromedpCtx, cancel = context.WithTimeout(chromedpCtx, timeout)
	defer cancel()

	var result CommodityMetadata
	var html string

	err := chromedp.Run(chromedpCtx,
		// 更高级的反检测脚本
		chromedp.Evaluate(`
			(() => {
				// 更隐蔽的方式隐藏webdriver
				const originalGetOwnPropertyDescriptor = Object.getOwnPropertyDescriptor;
				Object.getOwnPropertyDescriptor = function(obj, prop) {
					if (prop === 'webdriver' && obj === navigator) {
						return undefined;
					}
					return originalGetOwnPropertyDescriptor.call(this, obj, prop);
				};
				
				// 随机化一些浏览器属性
				Object.defineProperty(navigator, 'plugins', {
					get: () => [1, 2, 3, 4, 5],
				});
				
				Object.defineProperty(navigator, 'languages', {
					get: () => ['zh-CN', 'zh', 'en'],
				});
				
				// 添加一些真实的鼠标事件监听器
				document.addEventListener('mousemove', () => {});
				document.addEventListener('click', () => {});
				document.addEventListener('scroll', () => {});
				
				// 随机化canvas指纹
				const originalGetContext = HTMLCanvasElement.prototype.getContext;
				HTMLCanvasElement.prototype.getContext = function(type, ...args) {
					const context = originalGetContext.call(this, type, ...args);
					if (type === '2d') {
						const originalFillText = context.fillText;
						context.fillText = function(...args) {
							args[1] += Math.random() * 0.001;
							return originalFillText.call(this, ...args);
						};
					}
					return context;
				};
				
				// 随机化音频指纹
				const originalGetChannelData = AudioBuffer.prototype.getChannelData;
				AudioBuffer.prototype.getChannelData = function(channel) {
					const data = originalGetChannelData.call(this, channel);
					const newData = new Float32Array(data.length);
					for (let i = 0; i < data.length; i++) {
						newData[i] = data[i] + (Math.random() - 0.5) * 0.0001;
					}
					return newData;
				};
				
				// 模拟真实的浏览器环境
				Object.defineProperty(navigator, 'hardwareConcurrency', {
					get: () => 8,
				});
				
				Object.defineProperty(navigator, 'deviceMemory', {
					get: () => 8,
				});
				
				// 删除自动化相关属性
				delete window.cdc_adoQpoasnfa76pfcZLmcfl_Array;
				delete window.cdc_adoQpoasnfa76pfcZLmcfl_Promise;
				delete window.cdc_adoQpoasnfa76pfcZLmcfl_Symbol;
				
				// 更高级的反检测
				Object.defineProperty(navigator, 'webdriver', {
					get: () => undefined,
				});
				
				// 模拟真实的用户行为
				Object.defineProperty(navigator, 'onLine', {
					get: () => true,
				});
				
				// 随机化屏幕分辨率
				Object.defineProperty(screen, 'width', {
					get: () => 1920 + Math.floor(Math.random() * 100),
				});
				
				Object.defineProperty(screen, 'height', {
					get: () => 1080 + Math.floor(Math.random() * 100),
				});
				
				// 模拟真实的鼠标事件
				const originalAddEventListener = EventTarget.prototype.addEventListener;
				EventTarget.prototype.addEventListener = function(type, listener, options) {
					if (type === 'mousemove' || type === 'click' || type === 'scroll') {
						// 添加一些随机延迟
						setTimeout(() => {
							originalAddEventListener.call(this, type, listener, options);
						}, Math.random() * 10);
					} else {
						originalAddEventListener.call(this, type, listener, options);
					}
				};
			})()
		`, nil),

		// 导航到页面
		chromedp.Navigate(url),

		// 随机等待时间
		chromedp.Sleep(time.Duration(2+rand.Intn(3))*time.Second),

		// 模拟鼠标移动
		chromedp.Evaluate(`
			(() => {
				// 随机鼠标移动
				const events = ['mousemove', 'mouseover', 'mouseenter'];
				const elements = document.querySelectorAll('*');
				const randomElement = elements[Math.floor(Math.random() * elements.length)];
				
				events.forEach(eventType => {
					const event = new MouseEvent(eventType, {
						view: window,
						bubbles: true,
						cancelable: true,
						clientX: Math.random() * window.innerWidth,
						clientY: Math.random() * window.innerHeight
					});
					randomElement.dispatchEvent(event);
				});
			})()
		`, nil),

		// 模拟滚动
		chromedp.Evaluate(`
			(() => {
				// 随机滚动
				const scrollAmount = Math.random() * 100 + 50;
				window.scrollTo(0, scrollAmount);
				
				// 触发滚动事件
				window.dispatchEvent(new Event('scroll'));
			})()
		`, nil),

		chromedp.Sleep(time.Duration(1+rand.Intn(2))*time.Second),

		// 再次滚动
		chromedp.Evaluate(`
			(() => {
				const scrollAmount = Math.random() * 200 + 100;
				window.scrollTo(0, scrollAmount);
				window.dispatchEvent(new Event('scroll'));
			})()
		`, nil),

		chromedp.Sleep(time.Duration(1+rand.Intn(2))*time.Second),

		// 获取页面HTML
		chromedp.OuterHTML("body", &html),

		// 提取数据
		chromedp.Evaluate(`
			(() => {
				// 使用更通用的选择器
				const titleElement = document.querySelector('.title-info__text, .product-title, h1, .title');
				const brandElement = document.querySelector('.product-param__params__content__item__content__row__value__desc__text, .brand-name, .product-brand, .brand');
				
				// 提取图片，使用多种可能的选择器
				const detailImages = [];
				
				const swiperImgs = document.querySelectorAll('.head-figure__media-view__content');
				if (swiperImgs.length > 0) {
					swiperImgs.forEach(img => {
						const styleAttr = img.getAttribute('style');
        				if (styleAttr && styleAttr.startsWith('background-image')) {

							const url = styleAttr.replace(/&quot;/g, '"')
								.replace(/&amp;/g, '&')
								// 移除所有非URL部分
								.replace(/^.*background-image:\s*/i, '')
								.replace(/url\(/i, '')
								.replace(/\)\s*;?\s*$/g, '')
								// 移除引号
								.replace(/^["']/, '')
								.replace(/["']$/, '')
								// 清理空白
								.trim();
							detailImages.push(url);
						}
					});
				}

				const imgSelectors = [
					'head-figure__media-view__video__main-pic',
					'.product-big-img-list__every-img',
					'.product-images img',
					'.detail-images img',
					'.product-gallery img',
					'.product-img img',
					'.image-list img'
				];
				
				for (const selector of imgSelectors) {
					const imgs = document.querySelectorAll(selector);
					if (imgs.length > 0) {
						imgs.forEach(img => {
							if (img.src && img.src.startsWith('http')) {
								detailImages.push(img.src);
							}
						});
						break;
					}
				}

				return {
					title: titleElement ? titleElement.textContent.trim() : '',
					brand: brandElement ? brandElement.textContent.trim() : '',
					images: detailImages
				};
			})()
		`, &result),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract data: %w", err)
	}

	fmt.Println(userAgent)

	if result.Title == "" || len(result.Images) == 0 {
		return nil, errors.New("failed to extract required data: empty title or images")
	}

	return &result, nil
}
