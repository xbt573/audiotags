/***************************************************************************
   copyright            : (C) 2014 by Nick Sellen
   email                : code@nicksellen.co.uk
***************************************************************************/

/***************************************************************************
 *   This library is free software; you can redistribute it and/or modify  *
 *   it  under the terms of the GNU Lesser General Public License version  *
 *   2.1 as published by the Free Software Foundation.                     *
 *                                                                         *
 *   This library is distributed in the hope that it will be useful, but   *
 *   WITHOUT ANY WARRANTY; without even the implied warranty of            *
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU     *
 *   Lesser General Public License for more details.                       *
 *                                                                         *
 *   You should have received a copy of the GNU Lesser General Public      *
 *   License along with this library; if not, write to the Free Software   *
 *   Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  *
 *   USA                                                                   *
 ***************************************************************************/

package main

import (
	"io/ioutil"
	"path/filepath"
	"fmt"
	"log"
	"os"
	"strings"
	"image"
	_ "image/png"
	_ "image/jpeg"

	//"github.com/nbonaparte/audiotags"
	".."
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("pass path to file")
		return
	}

	if len(os.Args) % 2 == 1 {
		fmt.Println("when modifying file every tag must have value")
		return
	}

	file, err := audiotags.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println("Current tags")
	fmt.Println("------------")
	props := file.ReadTags()
	for k, v := range props {
		fmt.Printf("%s: %s\n", k, strings.Replace(strings.Replace(v, "\r\n", "\n", -1), "\r", "\n", -1))
	}

	audioProps := file.ReadAudioProperties()
	fmt.Printf("length: %d\nbitrate: %d\nsamplerate: %d\nchannels: %d\n",
		audioProps.LengthMs, audioProps.Bitrate, audioProps.Samplerate, audioProps.Channels)

	if len(os.Args) > 2 {
		tags := make(map[string]string)
		tmp_tag := ""
		for i, arg := range os.Args[2:] {
			// escape newlines so they are correctly represented in the tags
                        arg = strings.Replace(arg, "\\n", "\n", -1)
			if i % 2 == 0 {
				tmp_tag = arg
			} else {
				tags[tmp_tag] = arg
				fmt.Printf("setting %s to %s\n", tmp_tag, arg)
			}
		}
		test := file.WriteTags(tags)
		if !test {
			log.Fatal("failed to write tags")
		}
	}

	img_name := "cover"
	if imgs, err := filepath.Glob(img_name + ".*"); err == nil && len(imgs) > 0 {
		img, err := os.Open(imgs[0])
		if err != nil {
			log.Fatalln(err)
		}
		defer img.Close()

		b, err := ioutil.ReadAll(img)
		if err != nil {
			log.Fatalln(err)
		}

		// reset offset position
		_, err = img.Seek(0, 0)
		if err != nil {
			log.Fatalln(err)
		}

		cfg, mimetype, err := image.DecodeConfig(img)
		if err != nil {
			log.Fatalln(err)
		}

		var img_fmt int
		if mimetype == "jpeg" {
			img_fmt = audiotags.JPEG
		} else {
			img_fmt = audiotags.PNG
		}
		fmt.Println("adding picture...")
		fmt.Println(len(b))
		if !file.WritePicture(b, img_fmt, cfg.Width, cfg.Height) {
			log.Fatalln("failed to write picture")
		}
	}
}
