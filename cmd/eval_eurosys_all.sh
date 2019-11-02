#!/usr/bin/env bash


DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

declare -a model_list=(
  MLPerf_ResNet50_v1.5
  VGG16 VGG19
  MLPerf_Mobilenet_v1
  ResNet_v1_50 ResNet_v1_101 ResNet_v1_152 ResNet_v2_50 ResNet_v2_101 ResNet_v2_152
  Inception_ResNet_v2 Inception_v1 Inception_v2 Inception_v3 Inception_v4
  AI_Matrix_ResNet152 AI_Matrix_Densenet121 AI_Matrix_GoogleNet AI_Matrix_ResNet50
  BVLC_AlexNet_Caffe BVLC_GoogLeNet_Caffe
  MobileNet_v1_0.5_128
  MobileNet_v1_0.5_160
  MobileNet_v1_0.5_192 MobileNet_v1_0.5_224
  MobileNet_v1_0.25_128 MobileNet_v1_0.25_160 MobileNet_v1_0.25_192 MobileNet_v1_0.25_224
  MobileNet_v1_0.75_128 MobileNet_v1_0.75_160 MobileNet_v1_0.75_192 MobileNet_v1_0.75_224
  MobileNet_v1_1.0_128 MobileNet_v1_1.0_160 MobileNet_v1_1.0_192
  MobileNet_v1_1.0_224

  # Faster_RCNN_Inception_v2_COCO Faster_RCNN_NAS_COCO Faster_RCNN_ResNet101_COCO
  Faster_RCNN_ResNet50_COCO
  MLPerf_SSD_MobileNet_v1_300x300
  MLPerf_SSD_ResNet34_1200x1200
  Mask_RCNN_ResNet50_v2_Atrous_COCO
  Mask_RCNN_Inception_v2_COCO
  Mask_RCNN_Inception_ResNet_v2_Atrous_COCO Mask_RCNN_ResNet101_v2_Atrous_COCO
  SSD_Inception_v2_COCO
  SSD_MobileNet_v1_COCO
  SSD_MobileNet_v2_COCO
  SSD_MobileNet_v1_FPN_Shared_Box_Predictor_640x640_COCO14_Sync SSD_MobileNet_v1_PPN_Shared_Box_Predictor_300x300_COCO14_Sync

  DeepLabv3_Xception_65_PASCAL_VOC_Train_Val DeepLabv3_MobileNet_v2_PASCAL_VOC_Train_Val DeepLabv3_MobileNet_v2_DM_05_PASCAL_VOC_Train_Val
  SRGAN
)



if [ -f main ]; then
  rm main
fi


go build main.go

for i in "${model_list[@]}"; do
  ./eval_eurosys.sh 3.89.83.65 1 $i ${DIR}/../../eurosys_20_results/p2/gpu eurosys_gpu
  ./eval_eurosys.sh 3.89.83.65 1 $i ${DIR}/../../eurosys_20_results/p2/cpu eurosys_cpu

  ./eval_eurosys.sh 54.144.25.35 1 $i ${DIR}/../../eurosys_20_results/g3/gpu eurosys_gpu
  ./eval_eurosys.sh 54.144.25.35 1 $i ${DIR}/../../eurosys_20_results/g3/cpu eurosys_cpu

  ./eval_eurosys.sh 3.88.52.25 1 $i ${DIR}/../../eurosys_20_results/p3/gpu eurosys_gpu
  ./eval_eurosys.sh 3.88.52.25 1 $i ${DIR}/../../eurosys_20_results/p3/cpu eurosys_cpu
done

