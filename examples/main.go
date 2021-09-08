package main

import (
	"fmt"
	"github.com/tak0kada/jma"
	"time"
)

func main() {
	now, _ := time.Parse(time.RFC3339, "2021-09-05T13:32:38Z")

	g := jma.GeoCoordinate{33.737131, 137.226929}
	var zoom uint = 8

	if !g.IsValid() {
		fmt.Println("invalid geocoordinate %s. normalizing...")
		g = g.Normalize()
	}
	fmt.Printf("GeoCoordinate: %s\n", g)

	t := g.ToTile(zoom)
	fmt.Printf("Tile: %s\n", t)

	err := jma.DownloadImageTile(t, now, 0, "./tile.png")
	if err != nil {
		fmt.Println(err)
	}
	err = jma.DownloadImage(g, zoom, jma.Rect{600, 800}, now, 0, "./example.png")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
