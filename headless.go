package main

import (
	"context"
	"log"
	"time"

	chromedp "github.com/knq/chromedp"
	"github.com/knq/chromedp/cdp"
	"github.com/knq/chromedp/cdp/page"
	"github.com/knq/chromedp/client"
)

// Headless service
type Headless struct {
	cdp    *chromedp.CDP
	ctx    context.Context
	cancel context.CancelFunc
}

//Pdf model
type Pdf struct {
	Filename string
	Content  []byte
}

//NewHeadless service
func NewHeadless() Headless {
	var err error
	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// create chrome
	//cpd, err := chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)), chromedp.WithLog(log.Printf))
	cpd, err := chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)))

	if err != nil {
		log.Fatal(err)
	}

	return Headless{cpd, ctxt, cancel}
}

//PrintPdf generate pdf
func (h Headless) PrintPdf(name string, source string) Pdf {
	var buf []byte

	err := h.cdp.Run(
		h.ctx,
		chromedp.Tasks{
			chromedp.Navigate(source),
			chromedp.Sleep(500 * time.Millisecond),
			printPdfTask(&buf),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return Pdf{
		name + ".pdf",
		buf,
	}
}

func printPdfTask(buf *[]byte) chromedp.Action {
	return chromedp.ActionFunc(func(ctxt context.Context, h cdp.Handler) error {
		var err error
		*buf, err = page.PrintToPDF().Do(ctxt, h)

		if err != nil {
			return err
		}

		return err
	})
}
