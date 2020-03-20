package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"remind/send"

	"github.com/sirupsen/logrus"
)

var (
	confPath        string
	userConfPath    string
	messageConfPath string

	// 输出版本信息
	version bool
)

var (
	conf  *send.Conf
	users send.UserSource
	msg   *send.Content
)

func init() {
	// flag.BoolVar(&version, "v", false, "输出版本信息")
	flag.StringVar(&confPath, "conf", "conf/conf.yaml", "配置文件路径")
	flag.StringVar(&userConfPath, "user_conf", "conf/user_conf.yml", "通知用户配置文件路径")
	flag.StringVar(&messageConfPath, "message_conf", "message/notify.yaml", "添加备忘录信息")
}

func main() {
	flag.Parse()

	// fmt.Println("Version:", Version)
	// fmt.Println("Go Version:", runtime.Version())
	// fmt.Println("Build:", Build)
	// fmt.Println("Commit:", Commit)

	if version {
		os.Exit(0)
	}

	var err error
	conf, err = send.LoadConf(confPath)
	log := logrus.WithField("confPath", confPath)
	if err != nil {
		log.WithError(err).Error("load conf error")
		os.Exit(1)
	}

	log = log.WithField("conf", conf)

	users, err = send.LoadConfUsers(userConfPath)
	if err != nil {
		log.WithError(err).Error("load users failed")
		os.Exit(1)
	}

	msg, err = send.LoadMessage(messageConfPath)
	if err != nil {
		log.WithError(err).Error("load message content failed")
		os.Exit(1)
	}

	notify := make(chan os.Signal, 2)
	signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM)

	var w sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	if err := run(ctx, &w); err != nil {
		logrus.Errorf("run with error: %s", err)
		os.Exit(1)
	}

	logrus.Info("running")

	sig := <-notify
	logrus.Infof("received signal: %s", sig)

	if err := stop(cancel, &w); err != nil {
		logrus.Errorf("stop with error: %s", err)
	}

	logrus.Info("exit")
}

// run Consumer MQ queue("sms.channel.event") and store prometheus, Consume MQ queue("paas.notify") and send email
func run(ctx context.Context, w *sync.WaitGroup) error {
	w.Add(1)

	go func() {
		defer w.Done()
		smtp := &send.Smtp{
			Host:     conf.Smtp.Addr,
			Port:     conf.Smtp.Port,
			From:     conf.Smtp.From,
			Username: conf.Smtp.Username,
			Password: conf.Smtp.Password,
		}

		n := &send.Notify{
			Name:    "reciver",
			Title:   "your worker plan",
			Content: msg, // 改成content, level直接进行发送
		}
		// send email
		tos, err := users.Get(n)
		if err != nil {
			logrus.WithError(err).Error("get tos user failed")
		}

		err = smtp.Send(n, tos)
		if err != nil {
			logrus.WithError(err).Error("send email failed")
		}
	}()

	return nil
}

func stop(cancel context.CancelFunc, w *sync.WaitGroup) error {
	cancel()
	w.Wait()
	return nil
}
