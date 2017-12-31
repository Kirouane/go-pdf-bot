package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	chromedp "github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
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
func NewHeadless(url string) Headless {
	var err error
	// create context
	ctxt, cancel := context.WithCancel(context.Background())

	// create chrome
	//cpd, err := chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)), chromedp.WithLog(log.Printf))
	c := client.New()
	client.URL(url)(c)

	cpd, err := chromedp.New(ctxt, chromedp.WithTargets(c.WatchPageTargets(ctxt)))
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
			chromedp.ActionFunc(func(ctxt context.Context, e cdp.Executor) error {
				var err error
				buf, err = page.PrintToPDF().Do(ctxt, e)

				if err != nil {
					return err
				}

				return err
			}),
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
