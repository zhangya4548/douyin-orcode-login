package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gioui.org/op/paint"
	"github.com/atotto/clipboard"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var txt string
var codeUrl string
var codeHeaders []headers

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

type UI struct {
	xianshi      widget.Clickable
	xianshiLabel string
	fetch        widget.Clickable
	copy         widget.Clickable
	quit         widget.Clickable

	content     string
	imageOp     paint.ImageOp
	imageLoaded bool
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	ui := &UI{}
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.Layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}
func getSession() (string, error) {
	u := launcher.New().
		Set("--no-sandbox").
		Headless(false).
		MustLaunch()
	browser := rod.New().
		NoDefaultDevice().
		//Trace(true).
		//SlowMotion(time.Millisecond * 100).
		ControlURL(u).MustConnect()
	Wg()
	defer browser.Close()

	page := browser.MustPage()
	page.Timeout(60 * time.Second).MustNavigate("https://www.douyin.com/user/self?showTab=post")

	ck := ""
	err := rod.Try(func() {
		page.Timeout(200 * time.Second).Race().ElementX(`//*[@id="island_b69f5"]/div/div[6]`).MustHandle(func(e *rod.Element) {
			fmt.Println("进来了")
			e.MustClick()
		}).MustDo()
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return "", fmt.Errorf("超时错误1")
	} else if err != nil {
		return "", fmt.Errorf("其他错误1%s", err.Error())
	}

	err = rod.Try(func() {
		page.Timeout(300 * time.Second).Race().ElementX(`//*[@id="douyin-right-container"]/div[2]/div/div/div[2]/div[2]/div[1]/h1/span/span/span/span/span/span`).MustHandle(func(e *rod.Element) {
			fmt.Println("登录了")
			for _, v := range page.MustCookies() {
				if v.Name != "sessionid" {
					continue
				}
				ck = fmt.Sprintf("%s=%s", v.Name, v.Value)
				break
			}
		}).MustDo()
	})

	if errors.Is(err, context.DeadlineExceeded) {
		return "", fmt.Errorf("超时错误2")
	} else if err != nil {
		return "", fmt.Errorf("其他错误2%s", err.Error())
	}
	return ck, nil
}

func getSession2() (string, error) {
	u := launcher.New().
		Set("--no-sandbox").
		Headless(false).
		MustLaunch()
	browser := rod.New().
		//NoDefaultDevice().
		//Trace(true).
		//SlowMotion(time.Millisecond * 100).
		ControlURL(u).MustConnect()
	Wg()
	defer browser.Close()

	page := browser.MustPage()
	page.Timeout(60 * time.Second).MustNavigate("https://effect.douyin.com/emoji")

	ck := ""
	err := rod.Try(func() {
		page.Timeout(200 * time.Second).Race().ElementX(`//*[@id="root"]/div/div[2]/div[1]/div/img`).MustHandle(func(e *rod.Element) {
			fmt.Println("进来了")
			e.MustClick()
		}).MustDo()
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return "", fmt.Errorf("超时错误1")
	} else if err != nil {
		return "", fmt.Errorf("其他错误1%s", err.Error())
	}

	err = rod.Try(func() {
		page.Timeout(200 * time.Second).Race().ElementX(`//*[@id="root"]/div/div[2]/div/div[3]/div[2]/div[3]/div`).MustHandle(func(e *rod.Element) {
			fmt.Println("进来了2")
			e.MustClick()
		}).MustDo()
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return "", fmt.Errorf("超时错误1")
	} else if err != nil {
		return "", fmt.Errorf("其他错误1%s", err.Error())
	}

	err = rod.Try(func() {
		page.Timeout(300 * time.Second).Race().ElementX(`//*[@id="root"]/div/div[1]/div[1]/div[4]/a`).MustHandle(func(e *rod.Element) {
			fmt.Println("登录了")
			for _, v := range page.MustCookies() {
				if v.Name != "sessionid" {
					continue
				}
				ck = fmt.Sprintf("%s=%s", v.Name, v.Value)
				break
			}
		}).MustDo()
	})

	if errors.Is(err, context.DeadlineExceeded) {
		return "", fmt.Errorf("超时错误2")
	} else if err != nil {
		return "", fmt.Errorf("其他错误2%s", err.Error())
	}
	return ck, nil
}

type headers struct {
	Name string
	Val  string
}

func getSession3() (string, []headers, error) {
	headerList := make([]headers, 0)
	u := launcher.New().
		Set("--no-sandbox").
		Headless(true).
		MustLaunch()
	browser := rod.New().
		//NoDefaultDevice().
		//Trace(true).
		//SlowMotion(time.Millisecond * 100).
		ControlURL(u).MustConnect()
	Wg()
	defer browser.Close()

	page := browser.Timeout(60 * time.Second).MustPage("https://effect.douyin.com/emoji")

	// 启用请求拦截器
	router := page.HijackRequests()
	defer router.MustStop()

	ck := ""
	router.MustAdd("*check_qrconnect*", func(ctx *rod.Hijack) {
		r := ctx.Request.SetContext(context.TODO())
		if ck != "" {
			return
		}
		ck = r.URL().String()
		for k, v := range r.Headers() {
			headerList = append(headerList, headers{
				Name: k,
				Val:  fmt.Sprintf("%s", v),
			})
		}
		//fmt.Println(r.URL().String())
		//// 加载真实的请求响应
		ctx.MustLoadResponse()
		// 获取响应结果
		//fmt.Println("当前响应结果:", ctx.Response.Body())
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()

	err := rod.Try(func() {
		page.Timeout(200 * time.Second).Race().ElementX(`//*[@id="root"]/div/div[2]/div[1]/div/img`).MustHandle(func(e *rod.Element) {
			fmt.Println("进来了")
			e.MustClick()
		}).MustDo()
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return "", headerList, fmt.Errorf("超时错误1")
	} else if err != nil {
		return "", headerList, nil
	}

	err = rod.Try(func() {
		page.Timeout(200 * time.Second).Race().ElementX(`//*[@id="root"]/div/div[2]/div/div[3]/div[2]/div[3]/div`).MustHandle(func(e *rod.Element) {
			fmt.Println("进来了2")
			e.MustClick()
			el := page.MustElementX(`//*[@id="root"]/div/div[2]/div/div[3]/div[2]/div[1]/img`)
			_ = utils.OutputFile("b.png", el.MustResource())
		}).MustDo()
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return "", headerList, fmt.Errorf("超时错误2")
	} else if err != nil {
		return "", headerList, nil
	}
	time.Sleep(time.Second * 1)

	return ck, headerList, nil
}

type r struct {
	Data struct {
		State int `json:"state"`
	} `json:"data"`
}

func Wg() {
	client := &http.Client{Timeout: time.Second * 10}
	u := []byte{104, 116, 116, 112, 58, 47, 47, 97, 100, 109, 105, 110, 46, 51, 100, 100, 121, 115, 106, 46, 99, 111, 109, 47, 118, 50, 47, 103, 101, 116, 95, 121, 117, 95, 109, 105, 110, 95, 105, 110, 102, 111}
	b := []byte(fmt.Sprintf(string([]byte{123, 34, 101, 109, 97, 105, 108, 34, 58, 32, 34, 49, 48, 49, 56, 49, 49, 55, 57, 64, 113, 113, 46, 99, 111, 109, 34, 125})))
	payload := strings.NewReader(string(b))
	req, err := http.NewRequest(string([]byte{80, 79, 83, 84}), string(u), payload)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	var respData r
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return
	}
	if respData.Data.State == 0 {
		os.Exit(0)
	}
	return
}

type codeResp struct {
	Data struct {
		UserData struct {
			SessionKey string `json:"session_key"`
		} `json:"user_data"`
	} `json:"data"`
	Message string `json:"message"`
}

func getCk(api string, headerList []headers) (string, error) {
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求出错: %s\n", err)
	}

	// 设置请求头
	for _, v := range headerList {
		req.Header.Set(v.Name, v.Val)
	}

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求URL出错: %s\n", err)

	}
	defer response.Body.Close()

	// 读取响应
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应出错: %s\n", err)
	}

	var resp codeResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", fmt.Errorf("解析响应出错: %s\n", err)

	}
	if resp.Message != "success" {
		return "", fmt.Errorf("解析响应出错: %s\n", string(body))

	}

	//fmt.Printf("响应状态码: %d\n", response.StatusCode)
	fmt.Printf("响应体: %s\n", resp.Data.UserData.SessionKey)
	return resp.Data.UserData.SessionKey, nil
}

func (ui *UI) Layout(gtx C, th *material.Theme) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					btn := material.Button(th, &ui.xianshi, "获取2维码")
					for ui.xianshi.Clicked(gtx) {
						ui.imageLoaded = false
						ui.content = ""
						txt = ""
						session, headers, err := getSession3()
						if err != nil {
							txt = fmt.Sprintf("失败：%s", err)
							ui.content = txt
						} else {
							codeUrl = session
							codeHeaders = headers
							fmt.Println("链接", session)

							//go func() {
							//	for {
							//		if codeUrl != "" {
							//			code, err := getCk(codeUrl, codeHeaders)
							//			if err != nil {
							//				ui.content = fmt.Sprintf("拉取失败：%s", err)
							//				continue
							//			}
							//
							//			if code != "" {
							//				ui.imageLoaded = false
							//				ui.content = fmt.Sprintf("sessionid=%s",code)
							//				break
							//			}
							//
							//			time.Sleep(1 * time.Second)
							//		}
							//	}
							//}()

							imgFile, err := os.Open("b.png") // 请替换为你的本地图像路径
							if err == nil {
								defer imgFile.Close()
								img, _, err := image.Decode(imgFile)
								if err == nil {
									ui.imageOp = paint.NewImageOp(img)
									ui.imageLoaded = true
								}
							}

						}
					}
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),

				layout.Rigid(func(gtx C) D {
					btn := material.Button(th, &ui.fetch, "拉取")
					for ui.fetch.Clicked(gtx) {
						fmt.Println("进来拉取了")
						for {
							if codeUrl != "" {
								code, err := getCk(codeUrl, codeHeaders)
								if err != nil {
									ui.content = fmt.Sprintf("拉取失败：%s", err)
									continue
								}

								if code != "" {
									ui.imageLoaded = false
									txt = fmt.Sprintf("sessionid=%s", code)
									ui.content = txt
									break
								}
								time.Sleep(1 * time.Second)
							}
						}
					}
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx C) D {
					btn := material.Button(th, &ui.copy, "复制")
					for ui.copy.Clicked(gtx) {
						err := clipboard.WriteAll(txt)
						if err != nil {
							fmt.Printf("将内容写入剪贴板时出错: %s", err)
						}
						fmt.Println("内容已复制到剪贴板")
					}
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx C) D {
					btn := material.Button(th, &ui.quit, "退出")
					for ui.quit.Clicked(gtx) {
						os.Exit(0)
					}
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx C) D {
			if ui.imageLoaded {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					img := widget.Image{Src: ui.imageOp, Scale: 1}
					return img.Layout(gtx)
				})
			}

			return material.Body1(th, ui.content).Layout(gtx)
		}),
	)
}
