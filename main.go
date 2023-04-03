package main

import (
	"bytes"
	"dy-live/protobuf/protobuf"
	myProxy "dy-live/proxy"
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
)

var CaCertPath string
var CaKeyPath string

func main() {

	if runtime.GOOS == "windows" {
		CaCertPath = "\\proxy\\ca\\cert.pem"
		CaKeyPath = "\\proxy\\ca\\key.pem"
	} else {
		CaCertPath = "/proxy/ca/cert.pem"
		CaKeyPath = "/proxy/ca/key.pem"
	}

	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	pwd, _ := os.Getwd()
	caCert, err := ioutil.ReadFile(pwd + CaCertPath)
	if err != nil {
		log.Fatal(err)
	}
	caKey, err := ioutil.ReadFile(pwd + CaKeyPath)
	myProxy.SetCA(caCert, caKey)
	proxy.Verbose = *verbose

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	if err := os.MkdirAll("db", 0755); err != nil {
		log.Fatal("Can't create dir", err)
	}
	logger, err := myProxy.NewLogger("db")
	if err != nil {
		log.Fatal("can't open log file", err)
	}
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp.Request.Host == "live.douyin.com" {
			parseData(resp)
		}
		return resp
	})
	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal("listen:", err)
	}
	sl := myProxy.NewStoppableListener(l)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		log.Println("Got SIGINT exiting")
		sl.Add(1)
		sl.Close()
		logger.Close()
		sl.Done()
	}()
	log.Println("Starting Proxy")
	http.Serve(sl, proxy)
	sl.Wait()
	log.Println("All connections closed - exit")
}

func parseData(resp *http.Response) {
	/**
	body 会在读取完数据后被关闭回收 所以需要再造一个readCloser的body
	*/
	if resp != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		buf.Bytes()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))
		respBody := buf.Bytes()

		newResp := &protobuf.Response{}
		err := proto.Unmarshal(respBody, newResp)
		if err != nil {
			log.Println("proto unmarshal err", err)
		}
		messages := newResp.MessagesList
		// 正则匹配直播间的url id
		str := resp.Request.Referer()
		urlID := ""
		rege, err := regexp.Compile(`\d+`)
		if err != nil {
			log.Println("regexp compile error : ", err.Error())
		} else {
			urlID = rege.FindString(str)
		}

		for _, message := range messages {
			fmt.Println(message.Method)
			switch message.Method {
			case "WebcastMemberMessage":
				fmt.Println(message.Method, urlID)
			case "WebcastChatMessage":
				fmt.Println(message.Method)
			case "WebcastLikeMessage":
				fmt.Println(message.Method)
			case "WebcastGiftMessage":
				fmt.Println(message.Method)
				giftMessage := &protobuf.GiftMessage{}
				err := proto.Unmarshal(message.Payload, giftMessage)
				if err != nil {
					log.Println("Unmarshal", err.Error())
					return
				}
				fmt.Printf("收到%s价值%.2f元", giftMessage.Gift.Name, float32(giftMessage.Gift.DiamondCount)/10)
			case "WebcastRoomUserSeqMessage":
				fmt.Println(message.Method)
			case "WebcastGiftVoteMessage":
				fmt.Print(message.Method + " ---- ")
				fmt.Println(string(message.Payload))
			case "WebcastSocialMessage": // 有人分享流或关注主播
				socialMessage := &protobuf.SocialMessage{}
				proto.Unmarshal(message.Payload, socialMessage)
				fmt.Printf("socialMessage ---- %s \n", socialMessage)
				fmt.Printf("%s触发的消息 \n", socialMessage.User.NickName)
			default:
				fmt.Print(message.Method + " | ")
			}
		}
	}
}
