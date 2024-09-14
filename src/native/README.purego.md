
*** Type Conversions (Go <=> C)

string <=> char*
bool <=> _Bool
uintptr <=> uintptr_t
uint <=> uint32_t or uint64_t
uint8 <=> uint8_t
uint16 <=> uint16_t
uint32 <=> uint32_t
uint64 <=> uint64_t
int <=> int32_t or int64_t
int8 <=> int8_t
int16 <=> int16_t
int32 <=> int32_t
int64 <=> int64_t
float32 <=> float
float64 <=> double
struct <=> struct (WIP - darwin only)
func <=> C function
unsafe.Pointer, *T <=> void*
[]T => void*

*** https://pkg.go.dev/github.com/ebitengine/purego?GOOS=windows

func NewCallback(fn interface{}) uintptr
func RegisterFunc(fptr interface{}, cfn uintptr)
func RegisterLibFunc(fptr interface{}, handle uintptr, name string)
func SyscallN(fn uintptr, args ...uintptr) (r1, r2, err uintptr)

*** https://pkg.go.dev/github.com/ebitengine/purego?GOOS=linux

const (
	RTLD_DEFAULT = 0x00000 // Pseudo-handle for dlsym so search for any loaded symbol
	RTLD_LAZY    = 0x00001 // Relocations are performed at an implementation-dependent time.
	RTLD_NOW     = 0x00002 // Relocations are performed when the object is loaded.
	RTLD_LOCAL   = 0x00000 // All symbols are not made available for relocation processing by other modules.
	RTLD_GLOBAL  = 0x00100 // All symbols are available for relocation processing of other modules.
)

func Dlclose(handle uintptr) error
func Dlopen(path string, mode int) (uintptr, error)
func Dlsym(handle uintptr, name string) (uintptr, error)
func NewCallback(fn interface{}) uintptr
func RegisterFunc(fptr interface{}, cfn uintptr)
func RegisterLibFunc(fptr interface{}, handle uintptr, name string)
func SyscallN(fn uintptr, args ...uintptr) (r1, r2, err uintptr)
type Dlerror func (e Dlerror) Error() string

*** https://pkg.go.dev/github.com/ebitengine/purego?GOOS=darwin

const (
	RTLD_DEFAULT = ^uintptr(0) - 1 // Pseudo-handle for dlsym so search for any loaded symbol
	RTLD_LAZY    = 0x1             // Relocations are performed at an implementation-dependent time.
	RTLD_NOW     = 0x2             // Relocations are performed when the object is loaded.
	RTLD_LOCAL   = 0x4             // All symbols are not made available for relocation processing by other modules.
	RTLD_GLOBAL  = 0x8             // All symbols are available for relocation processing of other modules.
)

func Dlclose(handle uintptr) error
func Dlopen(path string, mode int) (uintptr, error)
func Dlsym(handle uintptr, name string) (uintptr, error)
func NewCallback(fn interface{}) uintptr
func RegisterFunc(fptr interface{}, cfn uintptr)
func RegisterLibFunc(fptr interface{}, handle uintptr, name string)
func SyscallN(fn uintptr, args ...uintptr) (r1, r2, err uintptr)
type Dlerror func (e Dlerror) Error() string
