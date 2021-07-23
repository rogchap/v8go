#ifndef V8GO_INTERNAL_H
#define V8GO_INTERNAL_H

#if defined(__MINGW32__) || defined(__MINGW64__)
// MinGW header files do not implicitly include windows.h
struct _EXCEPTION_POINTERS;
#endif

#include "libplatform/libplatform.h"
#include "v8.h"
#include "v8-inspector.h"

#include "v8go.h"

#endif
