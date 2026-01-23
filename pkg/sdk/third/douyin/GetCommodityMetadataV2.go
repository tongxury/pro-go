package douyin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/tebeka/selenium"
)

// 启动ChromeDriver服务
func startChromeDriver() (*exec.Cmd, error) {
	cmd := exec.Command("chromedriver", "--port=9515")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ChromeDriver: %w", err)
	}

	// 等待服务启动
	time.Sleep(2 * time.Second)

	return cmd, nil
}

func (t *Client) GetCommodityMetadataV2(ctx context.Context, url string) (*CommodityMetadata, error) {
	// 启动ChromeDriver服务
	cmd, err := startChromeDriver()
	if err != nil {
		return nil, fmt.Errorf("failed to start ChromeDriver: %w", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	// 配置Chrome选项
	caps := selenium.Capabilities{}
	caps["browserName"] = "chrome"
	caps["goog:chromeOptions"] = map[string]interface{}{
		"args": []string{
			"--headless",
			"--incognito",
			"--disable-blink-features=AutomationControlled",
			"--window-size=1920,1080",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-gpu",
			"--disable-web-security",
			"--disable-features=VizDisplayCompositor",
			"--disable-background-timer-throttling",
			"--disable-backgrounding-occluded-windows",
			"--disable-renderer-backgrounding",
			"--disable-field-trial-config",
			"--disable-ipc-flooding-protection",
		},
		"excludeSwitches": []string{
			"enable-automation",
		},
		"useAutomationExtension": false,
	}

	// 创建WebDriver
	wd, err := selenium.NewRemote(caps, "http://localhost:9515")
	if err != nil {
		return nil, fmt.Errorf("failed to create WebDriver: %w", err)
	}
	defer wd.Quit()

	// 设置页面加载超时
	if err := wd.SetImplicitWaitTimeout(30 * time.Second); err != nil {
		return nil, fmt.Errorf("failed to set implicit wait timeout: %w", err)
	}

	// 执行反检测脚本
	antiDetectionScript := `
		(() => {
			// 隐藏webdriver属性
			Object.defineProperty(navigator, 'webdriver', {
				get: () => undefined,
			});
			
			// 删除自动化相关属性
			delete window.cdc_adoQpoasnfa76pfcZLmcfl_Array;
			delete window.cdc_adoQpoasnfa76pfcZLmcfl_Promise;
			delete window.cdc_adoQpoasnfa76pfcZLmcfl_Symbol;
			
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
		})()
	`

	// 导航到页面
	if err := wd.Get(url); err != nil {
		return nil, fmt.Errorf("failed to navigate to URL: %w", err)
	}

	// 等待页面加载
	time.Sleep(3 * time.Second)

	// 执行反检测脚本
	_, err = wd.ExecuteScript(antiDetectionScript, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute anti-detection script: %w", err)
	}

	// 获取页面HTML
	html, err := wd.PageSource()
	if err != nil {
		return nil, fmt.Errorf("failed to get page source: %w", err)
	}

	fmt.Println(html)

	// 提取数据
	dataExtractionScript := `
		(() => {
			const titleElement = document.querySelector('.title-info__text');
			const priceElement = document.querySelector('.price-line__price-container__price__amount');
			const shopElement = document.querySelector('.shop-component__shop-content__basic-info__title-area__name');
			const brandElement = document.querySelector('.product-param__params__content__item__content__row__value__desc__text');

			// 提取详情图片
			const detailImages = [];
			const detailImgElements = document.querySelectorAll('.product-big-img-list__every-img');
			detailImgElements.forEach(img => {
				if (img.src) {
					detailImages.push(img.src);
				}
			});

			return {
				title: titleElement ? titleElement.textContent.trim() : '',
				brand: brandElement ? brandElement.textContent.trim() : '',
				images: detailImages
			};
		})()
	`

	_, err = wd.ExecuteScript(dataExtractionScript, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data: %w", err)
	}

	return nil, nil
}
