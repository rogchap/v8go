#include "v8go.h"
#include "v8.h"
#include "_cgo_export.h"

using namespace v8;

extern "C" {

const char* version() {
    return V8::GetVersion();
}

}

