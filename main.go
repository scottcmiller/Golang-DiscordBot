package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pdepip/go-binance/binance"
)


func init() {

	flag.StringVar(&Token, "t", "NDIyMTczOTczMjkyNTE1MzI5.DYX7rw.ShX5NdWRvny_7MX_crfKsN49Blc",
		"BOT TOKEN")
	flag.Parse()
}

var Token string


func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type Price struct{
	usd float64
	str string
}

func (b *Price) usdToStr() {
	b.str = "$"
	b.str += fmt.Sprintf("%.2f", b.usd)
	b.str += " USD"
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	client := binance.New("", "")
	var nanoQuery = binance.SymbolQuery{
		Symbol: "NANOBTC",
	}
	var btcQuery = binance.SymbolQuery{
		Symbol: "BTCUSDT",
	}

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})

	// If the message is "ping" reply with "Pong!"
	for {
		select {
		case <-ticker.C:
			{
				btcUSD, err := client.GetLastPrice(btcQuery)
				if err != nil {
					panic(err)
				}

				nanoBTC, err := client.GetLastPrice(nanoQuery)
				if err != nil {
					panic(err)
				}
				fmt.Println(btcUSD.Price)
				fmt.Println(nanoBTC.Price)
				nanoPrice := Price{btcUSD.Price * nanoBTC.Price, ""}
				nanoPrice.usdToStr()
				s.UpdateStreamingStatus(0, nanoPrice.str, "")
				fmt.Println(nanoPrice.str)
			}
			// do stuff
		case <-quit:
			ticker.Stop()
			return
		}
	}

}
