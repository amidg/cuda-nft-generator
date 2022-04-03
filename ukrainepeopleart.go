package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

/////////////////////////////////////////////
//GLOBAL VARIABLES
var (
	pathToLogFile            = "./NFTs/log.txt"
	pathToSourceImagesFolder = "./Source/"
	pathToCornerStuff        = "./Source/Corner/"
	pathToBackground         = "./Source/Background/"
	basepathgirl             = "./Source/Girl/"
	nftpathGirl              = "./NFTs/Girl/"
	pathtoGirlBody           = "./Source/Girl/body/"
	pathToGirlEyes           = "./Source/Girl/eyes/"
	pathToGirlHair           = "./Source/Girl/hair/"
	pathToGirlClothing       = "./Source/Girl/clothing/"
	pathToGirlExtra          = "./Source/Girl/extra/"
	basepathboy              = "./Source/Boy/"
	nftpathBoy               = "./NFTs/Boy/"
	pathtoBoyBody            = "./Source/Boy/body/"
	pathToBoyEyes            = "./Source/Boy/eyes/"
	pathToBoyHair            = "./Source/Boy/hair/"
	pathToBoyClothing        = "./Source/Boy/clothing/"
	pathToBoyExtra           = "./Source/Boy/extra/"
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

var library sourcelibrary

type ukranian struct {
	name      string
	bodytype  string
	eyestype  string
	hairtype  string
	dresstype string
	extra     string
	corner    string
	backgroud string
}

func createukranian(name string, body string, eyes string, hair string, dress string, extra string, corner string, background string) *ukranian {

	person := ukranian{name: name, bodytype: body, eyestype: eyes, hairtype: hair, dresstype: dress, extra: extra, corner: corner, backgroud: background}

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
		component.filename[iter+1] = f.Name() // 0 -> reserved for nothing
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
		getNumberOfFilesAtFolder(pathToCornerStuff, &library.corner)
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

	checkNumberOfAvailableImages(pathToCornerStuff, "corner")
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

/*
	type sourcelibrary struct {
	body       componentdata
	eyes       componentdata
	hair       componentdata
	clothing   componentdata
	extra      componentdata
	corner     componentdata
	background componentdata
}

*/

/////////////////////////////////////////////////////////////
// IMAGE FUNCTIONS:
func generateImageID(gender string) (imageID, body, eyes, hair, clothing, extra, corner, background string) {
	rand.Seed(time.Now().UnixNano())
	for iter := 0; iter < imageGenerationRetries; iter++ {
		body := library.body.filename[rand.Intn(library.body.filecounter)+1]
		eyes := library.eyes.filename[rand.Intn(library.eyes.filecounter)+1]
		hair := library.hair.filename[rand.Intn(library.hair.filecounter)+1]
		clothing := library.clothing.filename[rand.Intn(library.clothing.filecounter)+1]
		extra := library.extra.filename[rand.Intn(library.extra.filecounter)]
		corner := library.corner.filename[rand.Intn(library.corner.filecounter)]
		background := library.background.filename[rand.Intn(library.background.filecounter)+1]

		body = body[:len(body)-4]
		eyes = eyes[:len(eyes)-4]
		hair = hair[:len(hair)-4]
		clothing = clothing[:len(clothing)-4]

		if len(extra) == 0 {
			extra = "empty"
		} else if len(extra) > 4 {
			extra = extra[:len(extra)-4]
		}

		if len(corner) == 0 {
			corner = "empty"
		} else if len(corner) > 4 {
			corner = corner[:len(corner)-4]
		}

		background = background[:len(background)-4]

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

// func readAndDecodeImage() (image.Image, int, int) { //returns regenerated jpegs and its sizes
// 	rand.Seed(time.Now().UnixNano())
// 	var imageChosen = rand.Intn(numberOfAvailableImages) //choose random image between 0 and max image
// 	var imagePath = "ImgSource/" + wolfTemplateNames[imageChosen]
// 	//fmt.Println(imagePath)

// 	// Read image from file that already exists
// 	wolfImage, wolfImageErr := os.Open(imagePath) //example of image path "Source/wolf1.jpeg"
// 	if wolfImageErr != nil {
// 		// Handle error
// 	}
// 	defer wolfImage.Close()

// 	// directly decode image
// 	loadedWolfImage, loadedWolfImageErr := jpeg.Decode(wolfImage)
// 	if loadedWolfImageErr != nil {
// 		// Handle error
// 	}

// 	//determine image size
// 	imageRect := loadedWolfImage.Bounds()
// 	wolfImgWidth := imageRect.Dx()
// 	wolfImgHeight := imageRect.Dy()
// 	fmt.Print(wolfImgWidth)
// 	fmt.Print(" by ")
// 	fmt.Println(wolfImgHeight)

// 	return loadedWolfImage, wolfImgWidth, wolfImgHeight
// }

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
		checkNumberOfAvailableImages(pathToCornerStuff, "corner")
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
		switch *genderFlag {
		case "girl":
			for i := 0; i < *numberOfImagesFlag; i++ {

			}
		case "boy":
			for i := 0; i < *numberOfImagesFlag; i++ {

			}
		default:
			fmt.Println("please provide gender")
			os.Exit(3)
		}
	}
}
