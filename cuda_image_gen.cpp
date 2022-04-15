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

const int MAX_NUMBER_OF_NFT_ELEMENTS = 7;
const double alpha = 1;

string sourceBackground = "./Source/Background/BACK4.png";
string sourceBody = "./Source/girl/body/GIRL.png";
string sourceEyes = "./Source/girl/eyes/EYES1.png";
string sourceHair = "./Source/girl/hair/HAIR1.png";
string sourceDress = "./Source/girl/clothing/empty.png";
string sourceExtra = "./Source/girl/extra/ROSE.png";
string sourceCorner = "./Source/Corner/CORNER1.png";

string imagesForTesting[] = {sourceBackground, sourceBody, sourceEyes, sourceHair, sourceDress, sourceExtra, sourceCorner};

void addAlphaChannelTo(Mat& input) {
	cuda::GpuMat alphaMask(input.rows, input.cols, CV_8UC1, Scalar(255)); //8bit 1 channel alpha mask

	// alpha-channel for transperancy using GPU
	cuda::GpuMat tempImg; //, tempMask;
	cuda::GpuMat tempImageWithAlpha(input.rows, input.cols, input.type());
    vector<cuda::GpuMat> channels;

	if ( input.channels() == 3) {
		tempImg.upload(input);
		//tempMask.upload(mask);
					
		// break image into channels
		cuda::split(tempImg, channels); 

		// append alpha channel
		if (channels.size() == 3) { channels.push_back(alphaMask); };

		// combine channels
		cuda::merge(channels, tempImageWithAlpha); 
		tempImageWithAlpha.download(input); // download from GPU memory

		//cout << "Alpha Channel added: " << (input.channels() == 4) << endl;
	}
}

Mat overlayImagesUsingCUDA(Mat imageArray[], int imageArraySize) { 
	// images must have same size
	Mat result;

	// initialize image in GPU
	cuda::GpuMat gpuresult(imageArray[0].rows, imageArray[0].cols, CV_8UC4, Scalar(255)); //(imageArray[0].rows, imageArray[0].cols, imageArray[0].type()); // create new image

	for (int iter = 0; iter < imageArraySize; iter++) {
		if (imageArray[iter].channels() == 1) {
			cout << "empty image detected" << endl;
			break;
		}

		switch (iter) {
		case 0: // 0 must be some sort of background image
			addAlphaChannelTo(imageArray[iter]);
			gpuresult.upload(imageArray[iter]);
			//imageArray[iter].copyTo(gpuresult(Rect(0, 0, imageArray[iter].cols, imageArray[iter].rows)));
		default:
			// overlay two images in GPU -> must have alpha channel as well
			gpuresult.download(imageArray[iter-1]);
			cuda::addWeighted(imageArray[iter-1], 1, imageArray[iter], alpha, 0, gpuresult);
			//imageArray[iter].copyTo(NewImg(Rect(0, 0, imageArray[iter].cols, imageArray[iter].rows)));
		}
	}

	// download image from GPU memory
	gpuresult.download(result);

	return result;
}

int main(int argc, char *argv[]) {
	Mat img[MAX_NUMBER_OF_NFT_ELEMENTS];
	Mat completeImage;

	img[0] = imread(imagesForTesting[0], IMREAD_UNCHANGED);
	img[1] = imread(imagesForTesting[1], IMREAD_UNCHANGED);
	img[2] = imread(imagesForTesting[2], IMREAD_UNCHANGED);
	img[3] = imread(imagesForTesting[3], IMREAD_UNCHANGED);
	img[4] = imread(imagesForTesting[4], IMREAD_UNCHANGED);
	img[5] = imread(imagesForTesting[5], IMREAD_UNCHANGED);
	img[6] = imread(imagesForTesting[6], IMREAD_UNCHANGED);

	completeImage = overlayImagesUsingCUDA(img, sizeof(img)/sizeof(Mat));	

	imwrite("./NFTs/test.png", completeImage); // file is being saved

    return 0;
}