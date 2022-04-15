#include <opencv2/opencv.hpp>
#include <opencv2/highgui.hpp>
#include <opencv2/imgproc/imgproc.hpp>
#include <opencv2/core/cuda.hpp>
#include <opencv2/cudaimgproc.hpp>
#include <opencv2/cudaarithm.hpp>
#include <stdio.h>
#include <string.h>

using namespace cv;
using namespace std;

string sourceBackground = "./Source/Background/BACK4.png";
string sourceBody = "./Source/girl/body/GIRL.png";
string sourceEyes = "./Source/girl/eyes/EYES1.png";
string sourceHair = "./Source/girl/hair/HAIR1.png";
string sourceDress = "./Source/girl/clothing/empty.png";
string sourceExtra = "./Source/girl/extra/ROSE.png";
string sourceCorner = "./Source/Corner/CORNER1.png";

string imagesForTesting[] = {sourceBackground, sourceBody, sourceEyes, sourceHair, sourceDress, sourceExtra, sourceCorner};

Mat overlayTwoImagesAtZeroUsingCUDA(Mat imageArray[2]) { 
	// images must have same size
	Mat result;
	Mat mask = Mat::zeros(4096, 4096, CV_8UC4);

	// alpha-channel for transperancy using GPU
	cuda::GpuMat tempImg, tempMask;
	cuda::GpuMat tempImageWithAlpha(imageArray[0].rows, imageArray[0].cols, imageArray[0].type());
    vector<cuda::GpuMat> channels;

	// initialize image in GPU
	cuda::GpuMat NewImg(imageArray[0].rows, imageArray[0].cols, imageArray[0].type()); // create new image

	// process alpha-channel
	for (int i = 0; i < 2; i++) {
		// upload image and mask layer
		tempImg.upload(imageArray[i]);
    	tempMask.upload(mask);
		
		// break image into channels
		cuda::split(tempImg, channels); 
		cout << channels.size() << endl; // 3

		// append alpha channel
		if (channels.size() == 3) { channels.push_back(tempMask); };
		cout << channels.size() << endl; // 4

		// combine channels
		cuda::merge(channels, tempImageWithAlpha); 
		tempImageWithAlpha.download(imageArray[i]); // download from GPU memory

		// overlay two images in GPU -> must have alpha channel as well
		imageArray[i].copyTo(NewImg(Rect(0, 0, imageArray[i].cols, imageArray[i].rows)));
	}

	// download image from GPU memory
	NewImg.download(result);

	return result;
}

int main(int argc, char *argv[]) {
	Mat img[2];
	Mat completeImage;

	for(int i = 0; i < 1; i++) {
		img[0] = imread(imagesForTesting[0], IMREAD_COLOR);
		img[1] = imread(imagesForTesting[1], IMREAD_COLOR);

		completeImage = overlayTwoImagesAtZeroUsingCUDA(img);
	}

	imwrite("./NFTs/test.png", completeImage); // file is being saved

    return 0;
}