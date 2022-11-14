// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package teletekst

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/oliamb/cutter"
)

func PersistScreenshot101(p Page) {
	log.Printf(">>> Persisting a 101 screenshot for %s... ", p.Nr)
	createScreenshot(p, "101")
	cropScreenshot(p, "101")
	log.Printf("Done!\n")
}

func PersistScreenshotReply(p Page) {
	log.Printf(">>> Persisting a reply screenshot for %s... ", p.Nr)
	createScreenshot(p, "reply")
	cropScreenshot(p, "reply")
	log.Printf("Done!\n")
}

func cropScreenshot(p Page, prefix string) {
	path := fmt.Sprintf("/tmp/gowitness/screenshots/%s_%s.png", prefix, p.Nr)
	img, err := getImageFromFilePath(path)
	if err != nil {
		panic(err)
	}
	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:  485,
		Height: 583,
		Anchor: image.Point{544, 165},
	})
	if err != nil {
		panic(err)
	}

	croppedPath := fmt.Sprintf("/tmp/gowitness/screenshots/%s_%s_cropped.png", prefix, p.Nr)
	f, err := os.Create(croppedPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, croppedImg); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, err := png.Decode(f)
	return image, err
}

func createScreenshot(p Page, prefix string) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	resp := createContainer(ctx, cli, p, prefix)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	waitforContainer(ctx, cli, resp.ID)
}

func waitforContainer(ctx context.Context, cli *client.Client, id string) {
	statusCh, errCh := cli.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

}

func createContainer(ctx context.Context, cli *client.Client, p Page, prefix string) container.ContainerCreateCreatedBody {
	os.MkdirAll("/tmp/gowitness/screenshots", 0755)
	outputFile := fmt.Sprintf("%s_%s", prefix, p.Nr)
	url := fmt.Sprintf("https://nos.nl/teletekst#%s?t=%d", p.Nr, time.Now().UnixNano())
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "leonjza/gowitness",
		Cmd:   []string{"gowitness", "single", url, "--delay", "1", "-o", outputFile},
		Tty:   false,
	}, &container.HostConfig{Mounts: []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: "/tmp/gowitness",
			Target: "/data",
		},
	}}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	return resp
}
