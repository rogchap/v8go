// Copyright 2020 the V8 project authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef INCLUDE_CPPGC_INTERNAL_API_CONSTANTS_H_
#define INCLUDE_CPPGC_INTERNAL_API_CONSTANTS_H_

#include <cstddef>
#include <cstdint>

#include "v8config.h"  // NOLINT(build/include_directory)

namespace cppgc {
namespace internal {

// Embedders should not rely on this code!

// Internal constants to avoid exposing internal types on the API surface.
namespace api_constants {

constexpr size_t kKB = 1024;
constexpr size_t kMB = kKB * 1024;
constexpr size_t kGB = kMB * 1024;

// Offset of the uint16_t bitfield from the payload contaning the
// in-construction bit. This is subtracted from the payload pointer to get
// to the right bitfield.
static constexpr size_t kFullyConstructedBitFieldOffsetFromPayload =
    2 * sizeof(uint16_t);
// Mask for in-construction bit.
static constexpr uint16_t kFullyConstructedBitMask = uint16_t{1};

static constexpr size_t kPageSize = size_t{1} << 17;

#if defined(V8_TARGET_ARCH_ARM64) && defined(V8_OS_MACOS)
constexpr size_t kGuardPageSize = 0;
#else
constexpr size_t kGuardPageSize = 4096;
#endif

static constexpr size_t kLargeObjectSizeThreshold = kPageSize / 2;

#if defined(CPPGC_CAGED_HEAP)
#if defined(CPPGC_2GB_CAGE)
constexpr size_t kCagedHeapReservationSize = static_cast<size_t>(2) * kGB;
#else   // !defined(CPPGC_2GB_CAGE)
constexpr size_t kCagedHeapReservationSize = static_cast<size_t>(4) * kGB;
#endif  // !defined(CPPGC_2GB_CAGE)
constexpr size_t kCagedHeapReservationAlignment = kCagedHeapReservationSize;
#endif  // defined(CPPGC_CAGED_HEAP)

static constexpr size_t kDefaultAlignment = sizeof(void*);

// Maximum support alignment for a type as in `alignof(T)`.
static constexpr size_t kMaxSupportedAlignment = 2 * kDefaultAlignment;

// Granularity of heap allocations.
constexpr size_t kAllocationGranularity = sizeof(void*);

}  // namespace api_constants

}  // namespace internal
}  // namespace cppgc

#endif  // INCLUDE_CPPGC_INTERNAL_API_CONSTANTS_H_
