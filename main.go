package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/cheggaaa/pb/v3"

	"github.com/sqweek/dialog"
)

const infoHash = "f59269930d5d91dc4d65782d71399082cc355ce4"

var dhtNodes = [...]string{
	"router.bittorrent.com:6881",
	"dht.transmissionbt.com:6881",
	"router.utorrent.com:6881",
	"dht.libtorrent.org:25401",
	"dht.aelitis.com:6881"}

func main() {
	log.SetOutput(os.Stdout)

	log.Println("Select a directory")

	dir, err := dialog.Directory().Title("Select a directory").Browse()

	if err != nil {
		log.Fatalf("Failed to pick directory: %v", err)
	}

	log.Printf("Selected directory: %s\n", dir)

	log.Println("Creating torrent client...")

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = dir
	clientConfig.Seed = false
	client, err := torrent.NewClient(clientConfig)

	if err != nil {
		log.Fatalf("Failed to create torrent client: %v", err)
	}

	defer client.Close()

	log.Println("Resolving DHT nodes...")

	var dhtIps []string

	for _, node := range dhtNodes {
		host, port, err := net.SplitHostPort(node)

		if err != nil {
			log.Printf("Failed to split host:port: %v, skipping", err)

			continue
		}

		ips, err := net.LookupIP(host)

		if err != nil {
			log.Printf("Failed to lookup host %s: %v", host, err)

			continue
		}

		for _, ip := range ips {
			// library can't handle ipv6 - skip
			if ip.To4() == nil {
				continue
			}

			log.Printf("  Got IP %s, port %s\n", ip.String(), port)
			dhtIps = append(dhtIps, ip.String()+":"+port)
		}
	}

	client.AddDhtNodes(dhtIps)

	log.Println("Initializing torrent...")

	ih := metainfo.NewHashFromHex(infoHash)

	t, _ := client.AddTorrentInfoHash(ih)

	log.Println("Getting info...")
	<-t.GotInfo()
	log.Printf("Downloading %s...\n", t.Name())

	progressBar := pb.Full.Start64(t.Info().TotalLength())
	progressBar.Set(pb.Bytes, true)
	progressBar.SetRefreshRate(time.Second)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	interrupted := false

	log.Println("Ctrl-C to cancel")

	go func() {
		for {
			select {
			case <-quit:
				interrupted = true
				client.Close()
			default:
				progressBar.SetCurrent(t.BytesCompleted())
				time.Sleep(250 * time.Millisecond)
			}
		}
	}()

	t.DownloadAll()
	client.WaitAll()

	progressBar.Finish()

	if interrupted {
		log.Println("Download interrupted.")
	} else {
		log.Println("Download completed.")
	}
}
