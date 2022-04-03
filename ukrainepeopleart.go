package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/schollz/progressbar/v3"
)

/////////////////////////////////////////////
//GLOBAL VARIABLES
var (
	pathToLogFile            = "./NFTs/log.txt"
	nftpath                  = "./NFTs/"
	pathToSourceImagesFolder = "./Source/"
	pathToCorner             = "./Source/Corner/"
	pathToBackground         = "./Source/Background/"
	basepathgirl             = "./Source/girl/"
	pathtoGirlBody           = "./Source/girl/body/"
	pathToGirlEyes           = "./Source/girl/eyes/"
	pathToGirlHair           = "./Source/girl/hair/"
	pathToGirlClothing       = "./Source/girl/clothing/"
	pathToGirlExtra          = "./Source/girl/extra/"
	basepathboy              = "./Source/boy/"
	pathtoBoyBody            = "./Source/boy/body/"
	pathToBoyEyes            = "./Source/boy/eyes/"
	pathToBoyHair            = "./Source/boy/hair/"
	pathToBoyClothing        = "./Source/boy/clothing/"
	pathToBoyExtra           = "./Source/boy/extra/"
)

// image struct
const maxNumberOfImagesPerComponent = 100
const imageGenerationRetries = 100

type componentdata struct {
	componenttype string
	filename      [maxNumberOfImagesPerComponent]string
	filecounter   int
}

type sourcelibrary struct {
	body       componentdata
	eyes       componentdata
	hair       componentdata
	clothing   componentdata
	extra      componentdata
	corner     componentdata
	background componentdata
}

const imagecomponents = 7 // background, body, eyes, hair, clothing, extra, corner

const (
	Background int = 0
	Body       int = 1
	Eyes       int = 2
	Hair       int = 3
	Clothing   int = 4
	Extra      int = 5
	Corner     int = 6
)

var library sourcelibrary

type ukranian struct {
	name       string
	bodytype   string
	eyestype   string
	hairtype   string
	dresstype  string
	extra      string
	corner     string
	background string
}

func createukranian(name string, body string, eyes string, hair string, dress string, extra string, corner string, background string) *ukranian {

	person := ukranian{name: name, bodytype: body, eyestype: eyes, hairtype: hair, dresstype: dress, extra: extra, corner: corner, background: background}

	return &person
}

// BASIC FUNCTIONS ////////////////////////////////////////////////////////////
func readfile(filename string) (string, []string) {
	var contentAsSingleString string
	var contentLineByLine []string

	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		contentLineByLine = append(contentLineByLine, scanner.Text())
		contentAsSingleString = contentAsSingleString + (scanner.Text() + "\n")
	}

	return contentAsSingleString, contentLineByLine
}

func eliminateNewLineCrap(text string) string {
	os := runtime.GOOS
	if os == "windows" {
		return (strings.Replace(text, "\r\n", "", -1))
	} else {
		return (strings.Replace(text, "\n", "", -1))
	}
}

func getNumberOfFilesAtFolder(pathToFolder string, component *componentdata) {
	files, err := ioutil.ReadDir(pathToFolder)

	if err != nil {
		log.Fatal(err)
	}

	// component.filename[0] = "empty"
	for iter, f := range files {
		component.filename[iter] = f.Name() // 0 -> reserved for nothing
		component.filecounter = component.filecounter + 1
	}
}

func checkNumberOfAvailableImages(basepath string, category string) { //purely void function
	switch category {
	case "body":
		getNumberOfFilesAtFolder(basepath+category+"/", &library.body)
	case "eyes":
		getNumberOfFilesAtFolder(basepath+category+"/", &library.eyes)
	case "hair":
		getNumberOfFilesAtFolder(basepath+category+"/", &library.hair)
	case "clothing":
		getNumberOfFilesAtFolder(basepath+category+"/", &library.clothing)
	case "extra":
		getNumberOfFilesAtFolder(basepath+category+"/", &library.extra)
	case "corner":
		getNumberOfFilesAtFolder(pathToCorner, &library.corner)
	case "background":
		getNumberOfFilesAtFolder(pathToBackground, &library.background)
	}
}

func generatelibrary(gender string) {
	switch gender {
	case "girl":
		checkNumberOfAvailableImages(basepathgirl, "body")
		checkNumberOfAvailableImages(basepathgirl, "eyes")
		checkNumberOfAvailableImages(basepathgirl, "hair")
		checkNumberOfAvailableImages(basepathgirl, "clothing")
		checkNumberOfAvailableImages(basepathgirl, "extra")
	case "boy":
		checkNumberOfAvailableImages(basepathboy, "body")
		checkNumberOfAvailableImages(basepathboy, "eyes")
		checkNumberOfAvailableImages(basepathboy, "hair")
		checkNumberOfAvailableImages(basepathboy, "clothing")
		checkNumberOfAvailableImages(basepathboy, "extra")
	}

	checkNumberOfAvailableImages(pathToCorner, "corner")
	checkNumberOfAvailableImages(pathToBackground, "background")
}

func checkEntireLogFile(input string) (nomatch bool) {
	_, logSlice := readfile(pathToLogFile)
	currentlog := make([]string, len(logSlice))
	copy(currentlog, logSlice)
	nomatch = true

	for i := 0; i < len(currentlog); i++ {
		nomatch = nomatch && !(currentlog[i] == input)
	}

	return nomatch
}

/////////////////////////////////////////////////////////////
// IMAGE FUNCTIONS:
func recordImageID(imageID string) {
	logfile, openerr := os.OpenFile(pathToLogFile, os.O_APPEND|os.O_WRONLY, 0644)

	_, logerr := logfile.WriteString(imageID + "\n")

	if openerr != nil || logerr != nil {
		println("file error")
		os.Exit(3)
	}

	logfile.Close()
}

func generateImageID(gender string) (imageID, body, eyes, hair, clothing, extra, corner, background string) {
	rand.Seed(time.Now().UnixNano())
	for iter := 0; iter < imageGenerationRetries; iter++ {
		//n := a + rand.Intn(b-a+1)
		body = library.body.filename[rand.Intn(library.body.filecounter)]             // cannot be 0
		eyes = library.eyes.filename[rand.Intn(library.eyes.filecounter)]             // cannot be 0
		hair = library.hair.filename[rand.Intn(library.hair.filecounter)]             // cannot be 0
		clothing = library.clothing.filename[rand.Intn(library.clothing.filecounter)] // cannot be 0
		extra = library.extra.filename[rand.Intn(library.extra.filecounter)]
		corner = library.corner.filename[rand.Intn(library.corner.filecounter)]
		background = library.background.filename[rand.Intn(library.background.filecounter)] // cannot be 0

		body = body[:len(body)-len(".png")]
		eyes = eyes[:len(eyes)-len(".png")]
		hair = hair[:len(hair)-len(".png")]
		clothing = clothing[:len(clothing)-len(".png")]
		extra = extra[:len(extra)-len(".png")]
		corner = corner[:len(corner)-len(".png")]
		background = background[:len(background)-len(".png")]

		imageID = gender + "_" + body + eyes + hair + clothing + extra + corner + background

		if checkEntireLogFile(imageID) {
			break
		}

		if iter == imageGenerationRetries-1 {
			os.Exit(3) // stop program, error
		}
	}

	return imageID, body, eyes, hair, clothing, extra, corner, background
}

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	image, err := png.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer file.Close()

	return image, err
}

func createNFTimage(person *ukranian) (success bool, err error) { //returns regenerated jpegs and its sizes
	// read images from the person data
	imageID := person.name
	body := person.bodytype
	eyes := person.eyestype
	hair := person.hairtype
	clothing := person.dresstype
	extra := person.extra
	corner := person.corner
	background := person.background
	imageGender := imageID[:4]

	if imageGender == "boy_" {
		imageGender = imageGender[:3]
	}

	imageData := [imagecomponents]string{background, body, eyes, hair, clothing, extra, corner} // set in order of image merging

	// create default image
	ukrainenftimg := gg.NewContext(4096, 4096)
	var openedImage image.Image

	// add new layers from the person data
	for imageGenOrder := 0; imageGenOrder < imagecomponents; imageGenOrder++ {
		// choose layer
		if imageData[imageGenOrder] != "empty" {
			switch imageGenOrder {
			case Background:
				openedImage, err = openImage(pathToBackground + imageData[Background] + ".png")
			case Body:
				openedImage, err = openImage(pathToSourceImagesFolder + "/" + imageGender + "/" + "body/" + imageData[Body] + ".png")
			case Eyes:
				openedImage, err = openImage(pathToSourceImagesFolder + "/" + imageGender + "/" + "eyes/" + imageData[Eyes] + ".png")
			case Hair:
				openedImage, err = openImage(pathToSourceImagesFolder + "/" + imageGender + "/" + "hair/" + imageData[Hair] + ".png")
			case Clothing:
				openedImage, err = openImage(pathToSourceImagesFolder + "/" + imageGender + "/" + "clothing/" + imageData[Clothing] + ".png")
			case Extra:
				openedImage, err = openImage(pathToSourceImagesFolder + "/" + imageGender + "/" + "extra/" + imageData[Extra] + ".png")
			case Corner:
				openedImage, err = openImage(pathToCorner + imageData[Corner] + ".png")
			}
		}

		// check for errors
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}

		// add layer
		ukrainenftimg.DrawImage(openedImage, 0, 0)
	}

	// save image and log it
	ukrainenftimg.SavePNG(nftpath + "/" + imageGender + "/" + imageID)
	recordImageID(imageID)

	if err == nil {
		success = true
	}

	return success, err
}

// DEBUG FUNCTIONS ///////////////////////////////
func showListOfFiles(filetype string) {
	switch filetype {
	case "body":
		checkNumberOfAvailableImages(basepathgirl, "body")
		fmt.Println(library.body.filename)
		fmt.Println(library.body.filecounter)
	case "eyes":
		checkNumberOfAvailableImages(basepathgirl, "eyes")
		fmt.Println(library.eyes.filename)
		fmt.Println(library.eyes.filecounter)
	case "hair":
		checkNumberOfAvailableImages(basepathgirl, "hair")
		fmt.Println(library.hair.filename)
		fmt.Println(library.hair.filecounter)
	case "clothing":
		checkNumberOfAvailableImages(basepathgirl, "clothing")
		fmt.Println(library.clothing.filename)
		fmt.Println(library.clothing.filecounter)
	case "extra":
		checkNumberOfAvailableImages(basepathgirl, "extra")
		fmt.Println(library.extra.filename)
		fmt.Println(library.extra.filecounter)
	case "corner":
		checkNumberOfAvailableImages(pathToCorner, "corner")
		fmt.Println(library.corner.filename)
		fmt.Println(library.corner.filecounter)
	case "background":
		checkNumberOfAvailableImages(pathToBackground, "background")
		fmt.Println(library.background.filename)
		fmt.Println(library.background.filecounter)
	}
}

// MAIN
func main() {
	numberOfImagesFlag := flag.Int("nftnumber", 0, "command number of images for generation")
	genderFlag := flag.String("gender", "", "specify boy or girl")
	showListOfFileNames := flag.String("showfilelist", "", "show list of files of the specified type, e.g. body")
	randomIDs := flag.Int("randomids", 0, "generate random image id of the speicified count")
	flag.Parse()

	*genderFlag = eliminateNewLineCrap(*genderFlag)
	*showListOfFileNames = eliminateNewLineCrap(*showListOfFileNames)

	if *showListOfFileNames != "" {
		showListOfFiles(*showListOfFileNames)
		os.Exit(1) // end of program
	}

	generatelibrary(*genderFlag)
	if *randomIDs != 0 {
		for i := 0; i < *randomIDs; i++ {
			imageID, _, _, _, _, _, _, _ := generateImageID(*genderFlag)
			fmt.Println(imageID)
		}
	} else if *randomIDs == 0 {
		generation := progressbar.DefaultBytes(
			(int64(*numberOfImagesFlag)),
			"Generating NFTs",
		)

		for i := 0; i < *numberOfImagesFlag; i++ {
			imageid, body, eyes, hair, clothing, extra, corner, background := generateImageID(*genderFlag)
			person := createukranian(imageid, body, eyes, hair, clothing, extra, corner, background)
			success, err := createNFTimage(person)
			if !success {
				fmt.Println(err)
				os.Exit(3)
			}

			generation.Add(1)
		}
	}
}
