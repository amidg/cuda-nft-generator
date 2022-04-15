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
string sourceHair = "./Source/girl/hair/HAIR2.png";
string sourceDress = "./Source/girl/clothing/empty.png";
string sourceExtra = "./Source/girl/extra/ROSE.png";
string sourceCorner = "./Source/Corner/CORNER1.png";

string imagesForTesting[] = {sourceBackground, sourceBody, sourceEyes, sourceHair, sourceDress, sourceExtra, sourceCorner};

void makeBlackExtraBlack(Mat& input) {
	for (int r = 0; r < input.rows; ++r)
    {
        for (int c = 0; c < input.cols; ++c)
        {
            Vec4b& pixel = input.at<Vec4b>(r,c);
            if (pixel[0] < 20 && pixel[1] < 20 && pixel[2] < 20) {
				pixel[0] = 0;
				pixel[1] = 0;
				pixel[2] = 0;
				pixel[3] = 255;
			}

			input.at<Vec4b>(r,c) = pixel;
        }
    }
}

void addAlphaChannelTo(Mat& input) {
	cuda::GpuMat alphaMask(input.rows, input.cols, CV_8UC1, Scalar(255)); //8bit 1 channel alpha mask, should add white

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

	tempImg.release();
	tempImageWithAlpha.release(); // deallocate memory
}

Mat overlayImagesUsingCUDA(Mat imageArray[], int imageArraySize) { 
	// images must have same size
	Mat result;

	// initialize image in GPU
	cuda::GpuMat gpuresult(imageArray[0].rows, imageArray[0].cols, CV_8UC4);

	for (int iter = 0; iter < imageArraySize; iter++) {
		if (imageArray[iter].channels() == 1) {
			cout << "empty image detected: " << imagesForTesting[iter] << endl;
			break;
		}

		switch (iter) {
		case 0: // 0 must be some sort of background image
			addAlphaChannelTo(imageArray[iter]);
			gpuresult.upload(imageArray[iter]);
			//imageArray[iter].copyTo(gpuresult(Rect(0, 0, imageArray[iter].cols, imageArray[iter].rows)));
		default:
			// overlay two images in GPU -> must have alpha channel as well
			//gpuresult.download(imageArray[iter-1]);
			//cuda::bitwise_and(imageArray[iter-1], imageArray[iter], gpuresult);
			// cuda::addWeighted(gpuresult, 1, imageArray[iter], alpha, 0, gpuresult);
			cuda::add(gpuresult, imageArray[iter], gpuresult);
			//imageArray[iter].copyTo(gpuresult(Rect(0, 0, imageArray[iter].cols, imageArray[iter].rows))); // direct overlay, no blending
		}
	}

	// download image from GPU memory
	gpuresult.download(result);

	return result;
}

/*
	build: g++ cuda_image_gen.cpp -o testcuda -g `pkg-config --libs --cflags opencv4`
*/

int main(int argc, char *argv[]) {
	Mat characterImageArray[MAX_NUMBER_OF_NFT_ELEMENTS];

	Mat completeImage;

	// read images from source
	for (int i = 0; i < MAX_NUMBER_OF_NFT_ELEMENTS; i++) {
		characterImageArray[i] = imread(imagesForTesting[i], IMREAD_UNCHANGED);
		makeBlackExtraBlack(characterImageArray[i]);
	}

	// build image
	completeImage = overlayImagesUsingCUDA(characterImageArray, sizeof(characterImageArray)/sizeof(Mat));	

	imwrite("./NFTs/test.png", completeImage); // file is being saved

	cout << "done" << endl;

    return 0;
}