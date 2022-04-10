#include <jetson-utils/cudaOverlay.h>
#include <jetson-utils/cudaMappedMemory.h>
#include <jetson-utils/imageIO.h>

#include <algorithm>  // for std::max()

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

uchar3* bodyImage = NULL;
uchar3* eyesImage = NULL;
uchar3* hairImage = NULL;
uchar3* dressImage = NULL;
uchar3* extraImage = NULL;
uchar3* cornerImage = NULL;
uchar3* backImage = NULL;

int2 dimsA = make_int2(0,0);
int2 dimsB = make_int2(0,0);

int main(int argc, char *argv[]) {
    // load the input images
    if( !loadImage("my_image_a.jpg", &imgInputA, &dimsA.x, &dimsA.y) )
        return false;

    if( !loadImage("my_image_b.jpg", &imgInputB, &dimsB.x, &dimsB.y) )
        return false;

    // allocate the output image, with dimensions to fit both inputs side-by-side
    const int2 dimsOutput = make_int2(dimsA.x, dimsA.y);

    if( !cudaAllocMapped(&imgOutput, dimsOutput.x, dimsOutput.y) )
        return false;

    // compost the two images (the last two arguments are x,y coordinates in the output image)
    CUDA(cudaOverlay(bodyImage, dimsA, imgOutput, dimsOutput, 0, 0));
    CUDA(cudaOverlay(eyesImage, dimsA, imgOutput, dimsOutput, 0, 0));
    CUDA(cudaOverlay(hairImage, dimsA, imgOutput, dimsOutput, 0, 0));
    CUDA(cudaOverlay(dressImage, dimsA, imgOutput, dimsOutput, 0, 0));
    CUDA(cudaOverlay(extraImage, dimsA, imgOutput, dimsOutput, 0, 0));
    CUDA(cudaOverlay(cornerImage, dimsA, imgOutput, dimsOutput, 0, 0));
    CUDA(cudaOverlay(backImage, dimsA, imgOutput, dimsOutput, 0, 0));
    return 0;
}