package imitation

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	configDir  = "configs"
	configFile = "imitation"
)

type Response struct {
	Time    time.Time
	Message string
}

type Request struct {
	Channel chan Response
	Mu      sync.Mutex
}

type Durations struct {
	SenderTime   time.Duration `mapstructure:"sender_time"`
	ReceiverTime time.Duration `mapstructure:"receiver_time"`
	StopTime     time.Duration `mapstructure:"stop_time"`
}

type Imitation struct {
	Durations Durations
	Request   Request
}

func New() (*Imitation, error) {
	dur := new(Durations)

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(dur); err != nil {
		return nil, err
	}

	return &Imitation{
		Durations: *dur,
		Request: Request{
			Channel: make(chan Response, 10),
		},
	}, nil
}

func (i *Imitation) Sender(ctx context.Context, wg *sync.WaitGroup, message string) {
	defer wg.Done()

	timer := time.NewTimer(time.Second * time.Duration(i.Durations.SenderTime.Seconds()))
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("sender was stopped")
			return
		case <-timer.C:
			i.Request.Mu.Lock()
			i.Request.Channel <- Response{
				Time:    time.Now(),
				Message: message,
			}
			i.Request.Mu.Unlock()

			timer.Reset(time.Second * time.Duration(i.Durations.SenderTime.Seconds()))
		}
	}
}

func (i *Imitation) Receiver(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	timer := time.NewTimer(time.Second * time.Duration(i.Durations.ReceiverTime.Seconds()))
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("receiver was stopped")
			return
		case <-timer.C:
			// check the amount of data in the channel before printing
			channelLenght := len(i.Request.Channel)

			for j := 0; j < channelLenght; j++ {
				value := <-i.Request.Channel
				log.Printf(" full time: %v; message: %s\n", value.Time.Format(time.ANSIC), value.Message)
			}

			timer.Reset(time.Second * time.Duration(i.Durations.ReceiverTime.Seconds()))
		}
	}
}

func (i *Imitation) StopAll(cancel context.CancelFunc) {
	time.Sleep(time.Second * time.Duration(i.Durations.StopTime.Seconds()))

	cancel()

	i.Request.Mu.Lock()
	close(i.Request.Channel)
	i.Request.Mu.Unlock()
}
