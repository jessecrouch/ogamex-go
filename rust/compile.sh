#!/bin/sh

# Compile the rust workspace
echo "Compiling Rust workspace..."
if ! cargo build --release; then
    echo "ERROR: Rust compilation failed!"
    exit 1
fi

# Copy the compiled rust libraries to the storage/rust-libs directory.
if [ -f target/release/libbattle_engine_ffi.so ]; then
    mkdir -p ../storage/rust-libs
    cp target/release/libbattle_engine_ffi.so ../storage/rust-libs/
    echo "Copied libbattle_engine_ffi.so"
else
    echo "ERROR: libbattle_engine_ffi.so not found after compilation!"
    exit 1
fi

if [ -f target/release/libtest_ffi.so ]; then
    cp target/release/libtest_ffi.so ../storage/rust-libs/
    echo "Copied libtest_ffi.so"
fi

echo "Rust compilation completed successfully!"
