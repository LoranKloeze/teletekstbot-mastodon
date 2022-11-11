// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package teletekst

import (
	"context"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/oliamb/cutter"
)

func PersistScreenshot(p Page) {
	log.Printf(">>> Persisting a screenshot for %s... ", p.Nr)
	createScreenshot(p)
	cropScreenshot(p)
	log.Printf("Done!\n")
}

func cropScreenshot(p Page) {
	img, err := getImageFromFilePath("/tmp/gowitness/screenshots/" + p.Nr + ".png")
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

	f, err := os.Create("/tmp/gowitness/screenshots/" + p.Nr + "_cropped.png")
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

func createScreenshot(p Page) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	resp := createContainer(ctx, cli, p)

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

func createContainer(ctx context.Context, cli *client.Client, p Page) container.ContainerCreateCreatedBody {
	os.MkdirAll("/tmp/gowitness/screenshots", 0755)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "leonjza/gowitness",
		Cmd:   []string{"gowitness", "single", "https://nos.nl/teletekst#" + p.Nr, "--delay", "1", "-o", p.Nr},
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
