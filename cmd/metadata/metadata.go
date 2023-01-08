package main

import (
	"flag"
	"fmt"
	"github.com/psyton/audiotags"
	"os"
)

var usage = func() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: %s [optional flags] filename\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		return
	}

	f, err := audiotags.Open(flag.Arg(0))
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return
	}

	if !f.HasMedia() {
		fmt.Printf("No supported media in file\n")
		return
	}

	fmt.Printf("\nTags: \n")
	for tag, value := range f.ReadTags() {
		fmt.Printf("%s: %v\n", tag, value)
	}
	cover, err := f.ReadImage()
	if err != nil {
		fmt.Printf("error reading cover: %v\n", err)
	}
	if cover != nil {
		fmt.Printf("Cover: %dx%d\n", cover.Bounds().Size().X, cover.Bounds().Size().Y)
	}

	fmt.Printf("\nProps: \n")
	props := f.ReadAudioProperties()
	fmt.Printf("Bitrate: %v\n", props.Bitrate)
	fmt.Printf("Length: %v\n", props.Length)
	fmt.Printf("LengthMs: %v\n", props.LengthMs)
	fmt.Printf("Samplerate: %d\n", props.Samplerate)
	fmt.Printf("Channels: %d\n", props.Channels)
}
