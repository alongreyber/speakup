#!/bin/bash

# Prepare Kaldi
cd kaldi/tools
make -j4
make atlas openfst OPENFST_VERSION=1.4.1
cd ../src
make -j4
./configure --static --static-math=yes --static-fst=yes --use-cuda=no
make depend
cd ../../
