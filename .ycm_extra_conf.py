def Settings( **kwargs ):
  return {
    'flags': [ '-fno-rtti', '-fpic', '-std=c++11', '-Ideps/include', '-pthread', '-lv8', '-Ldeps/darwin-x86_64' ],
  }
