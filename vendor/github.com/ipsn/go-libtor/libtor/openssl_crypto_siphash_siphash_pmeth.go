// go-libtor - Self-contained Tor from Go
// Copyright (c) 2018 Péter Szilágyi. All rights reserved.

package libtor

/*
#define DSO_NONE
#define OPENSSLDIR "/usr/local/ssl"
#define ENGINESDIR "/usr/local/lib/engines"

#include <../crypto/siphash/siphash_pmeth.c>
*/
import "C"
