#!/bin/bash
set -e

# Build script for whisper.cpp static library
# Optimized for macOS (Apple Silicon) and Linux

WHISPER_DIR="internal/transcription/whisper.cpp"
BUILD_DIR="$WHISPER_DIR/build_static"
mkdir -p "$BUILD_DIR"

echo "Building whisper.cpp static library..."

CFLAGS="-O3 -DNDEBUG -std=c11 -fPIC -D_XOPEN_SOURCE=600 -DGGML_VERSION=\"1.0.0\" -DGGML_COMMIT=\"unknown\""
CXXFLAGS="-O3 -DNDEBUG -std=c++11 -fPIC -D_XOPEN_SOURCE=600 -DGGML_VERSION=\"1.0.0\" -DGGML_COMMIT=\"unknown\""

# Include paths
INCLUDES="-I$WHISPER_DIR/include -I$WHISPER_DIR/ggml/include"

# Platform specific flags
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS (Apple Silicon)
    CFLAGS="$CFLAGS -DGGML_USE_ACCELERATE -DGGML_USE_METAL -framework Accelerate"
    CXXFLAGS="$CXXFLAGS -DGGML_USE_ACCELERATE -DGGML_USE_METAL -framework Accelerate"
    LDFLAGS="-framework Accelerate -framework Metal -framework Foundation -framework CoreGraphics"
else
    # Linux (CPU only for now)
    CFLAGS="$CFLAGS -fopenmp"
    CXXFLAGS="$CXXFLAGS -fopenmp"
    LDFLAGS="-fopenmp"
fi

# Compile GGML
echo "Compiling GGML..."
cc $CFLAGS $INCLUDES -c "$WHISPER_DIR/ggml/src/ggml.c" -o "$BUILD_DIR/ggml.o"
cc $CFLAGS $INCLUDES -c "$WHISPER_DIR/ggml/src/ggml-alloc.c" -o "$BUILD_DIR/ggml-alloc.o"
cc $CFLAGS $INCLUDES -c "$WHISPER_DIR/ggml/src/ggml-backend.c" -o "$BUILD_DIR/ggml-backend.o"
cc $CFLAGS $INCLUDES -c "$WHISPER_DIR/ggml/src/ggml-quants.c" -o "$BUILD_DIR/ggml-quants.o"

# Compile Whisper
echo "Compiling Whisper..."
c++ $CXXFLAGS $INCLUDES -c "$WHISPER_DIR/src/whisper.cpp" -o "$BUILD_DIR/whisper.o"

# Create static library
echo "Creating libwhisper.a..."
ar rcs "$BUILD_DIR/libwhisper.a" "$BUILD_DIR/ggml.o" "$BUILD_DIR/ggml-alloc.o" "$BUILD_DIR/ggml-backend.o" "$BUILD_DIR/ggml-quants.o" "$BUILD_DIR/whisper.o"

echo "Build complete: $BUILD_DIR/libwhisper.a"
