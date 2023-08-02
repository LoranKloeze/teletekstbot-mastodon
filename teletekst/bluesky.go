package teletekst

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/karalabe/go-bluesky"
)

func Post101Bluesky(p Page) {
	log.Printf(">>> Posting a 101 bluesky for %s... ", p.Nr)
	ctx := context.Background()

	blueskyHandle := os.Getenv("BLUESKY_HANDLE")
	blueskyAppkey := os.Getenv("BLUESKY_PASSWORD")

	client, err := bluesky.Dial(ctx, bluesky.ServerBskySocial)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	err = client.Login(ctx, blueskyHandle, blueskyAppkey)
	switch {
	case errors.Is(err, bluesky.ErrMasterCredentials):
		panic("You're not allowed to use your full-access credentials, please create an appkey")
	case errors.Is(err, bluesky.ErrLoginUnauthorized):
		panic("Username of application password seems incorrect, please double check")
	case err != nil:
		panic("Something else went wrong, please look at the returned error")
	}

	url := "https://nos.nl/teletekst#" + p.Nr
	text := fmt.Sprintf("[%s] %s\n%s", p.Nr, p.Title, url)

	client.CustomCall(func(api *xrpc.Client) error {
		post := &bsky.FeedPost{
			Text:      text,
			CreatedAt: time.Now().Local().Format(time.RFC3339),
		}

		path := fmt.Sprintf("/tmp/gowitness/screenshots/%s_%s_cropped.png", "101", p.Nr)

		b, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}

		resp, err := atproto.RepoUploadBlob(ctx, api, bytes.NewReader(b))
		if err != nil {
			fmt.Println(err)
		}

		var images []*bsky.EmbedImages_Image
		images = append(images, &bsky.EmbedImages_Image{
			Alt: "Teletekst screenshot",
			Image: &util.LexBlob{
				Ref:      resp.Blob.Ref,
				MimeType: http.DetectContentType(b),
				Size:     resp.Blob.Size,
			},
		})

		post.Embed = &bsky.FeedPost_Embed{}
		post.Embed.EmbedImages = &bsky.EmbedImages{
			Images: images,
		}

		input := atproto.RepoCreateRecord_Input{
			Collection: "app.bsky.feed.post",
			Repo:       api.Auth.Did,
			Record: &util.LexiconTypeDecoder{
				Val: post,
			},
		}

		_, err = atproto.RepoCreateRecord(ctx, api, &input)
		if err != nil {
			fmt.Println(err)
		}
		return err
	})

	log.Printf("Done!\n")
}
