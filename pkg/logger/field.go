package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Any(key string, value interface{}) zap.Field              { return zap.Any(key, value) }
func Array(key string, val zapcore.ArrayMarshaler) zap.Field   { return zap.Array(key, val) }
func Binary(key string, val []byte) zap.Field                  { return zap.Binary(key, val) }
func Bool(key string, val bool) zap.Field                      { return zap.Bool(key, val) }
func Boolp(key string, val *bool) zap.Field                    { return zap.Boolp(key, val) }
func Bools(key string, bs []bool) zap.Field                    { return zap.Bools(key, bs) }
func ByteString(key string, val []byte) zap.Field              { return zap.ByteString(key, val) }
func ByteStrings(key string, bss [][]byte) zap.Field           { return zap.ByteStrings(key, bss) }
func Complex128(key string, val complex128) zap.Field          { return zap.Complex128(key, val) }
func Complex128p(key string, val *complex128) zap.Field        { return zap.Complex128p(key, val) }
func Complex128s(key string, nums []complex128) zap.Field      { return zap.Complex128s(key, nums) }
func Complex64(key string, val complex64) zap.Field            { return zap.Complex64(key, val) }
func Complex64p(key string, val *complex64) zap.Field          { return zap.Complex64p(key, val) }
func Complex64s(key string, nums []complex64) zap.Field        { return zap.Complex64s(key, nums) }
func Duration(key string, val time.Duration) zap.Field         { return zap.Duration(key, val) }
func Durationp(key string, val *time.Duration) zap.Field       { return zap.Durationp(key, val) }
func Durations(key string, ds []time.Duration) zap.Field       { return zap.Durations(key, ds) }
func Error(err error) zap.Field                                { return zap.Error(err) }
func Errors(key string, errs []error) zap.Field                { return zap.Errors(key, errs) }
func Float32(key string, val float32) zap.Field                { return zap.Float32(key, val) }
func Float32p(key string, val *float32) zap.Field              { return zap.Float32p(key, val) }
func Float32s(key string, nums []float32) zap.Field            { return zap.Float32s(key, nums) }
func Float64(key string, val float64) zap.Field                { return zap.Float64(key, val) }
func Float64p(key string, val *float64) zap.Field              { return zap.Float64p(key, val) }
func Float64s(key string, nums []float64) zap.Field            { return zap.Float64s(key, nums) }
func Inline(val zapcore.ObjectMarshaler) zap.Field             { return zap.Inline(val) }
func Int(key string, val int) zap.Field                        { return zap.Int(key, val) }
func Int16(key string, val int16) zap.Field                    { return zap.Int16(key, val) }
func Int16p(key string, val *int16) zap.Field                  { return zap.Int16p(key, val) }
func Int16s(key string, nums []int16) zap.Field                { return zap.Int16s(key, nums) }
func Int32(key string, val int32) zap.Field                    { return zap.Int32(key, val) }
func Int32p(key string, val *int32) zap.Field                  { return zap.Int32p(key, val) }
func Int32s(key string, nums []int32) zap.Field                { return zap.Int32s(key, nums) }
func Int64(key string, val int64) zap.Field                    { return zap.Int64(key, val) }
func Int64p(key string, val *int64) zap.Field                  { return zap.Int64p(key, val) }
func Int64s(key string, nums []int64) zap.Field                { return zap.Int64s(key, nums) }
func Int8(key string, val int8) zap.Field                      { return zap.Int8(key, val) }
func Int8p(key string, val *int8) zap.Field                    { return zap.Int8p(key, val) }
func Int8s(key string, nums []int8) zap.Field                  { return zap.Int8s(key, nums) }
func Intp(key string, val *int) zap.Field                      { return zap.Intp(key, val) }
func Ints(key string, nums []int) zap.Field                    { return zap.Ints(key, nums) }
func NamedError(key string, err error) zap.Field               { return zap.NamedError(key, err) }
func Namespace(key string) zap.Field                           { return zap.Namespace(key) }
func Object(key string, val zapcore.ObjectMarshaler) zap.Field { return zap.Object(key, val) }
func Reflect(key string, val interface{}) zap.Field            { return zap.Reflect(key, val) }
func Skip() zap.Field                                          { return zap.Skip() }
func Stack(key string) zap.Field                               { return zap.Stack(key) }
func StackSkip(key string, skip int) zap.Field                 { return zap.StackSkip(key, skip) }
func String(key string, val string) zap.Field                  { return zap.String(key, val) }
func Stringer(key string, val fmt.Stringer) zap.Field          { return zap.Stringer(key, val) }
func Stringp(key string, val *string) zap.Field                { return zap.Stringp(key, val) }
func Strings(key string, ss []string) zap.Field                { return zap.Strings(key, ss) }
func Time(key string, val time.Time) zap.Field                 { return zap.Time(key, val) }
func Timep(key string, val *time.Time) zap.Field               { return zap.Timep(key, val) }
func Times(key string, ts []time.Time) zap.Field               { return zap.Times(key, ts) }
func Uint(key string, val uint) zap.Field                      { return zap.Uint(key, val) }
func Uint16(key string, val uint16) zap.Field                  { return zap.Uint16(key, val) }
func Uint16p(key string, val *uint16) zap.Field                { return zap.Uint16p(key, val) }
func Uint16s(key string, nums []uint16) zap.Field              { return zap.Uint16s(key, nums) }
func Uint32(key string, val uint32) zap.Field                  { return zap.Uint32(key, val) }
func Uint32p(key string, val *uint32) zap.Field                { return zap.Uint32p(key, val) }
func Uint32s(key string, nums []uint32) zap.Field              { return zap.Uint32s(key, nums) }
func Uint64(key string, val uint64) zap.Field                  { return zap.Uint64(key, val) }
func Uint64p(key string, val *uint64) zap.Field                { return zap.Uint64p(key, val) }
func Uint64s(key string, nums []uint64) zap.Field              { return zap.Uint64s(key, nums) }
func Uint8(key string, val uint8) zap.Field                    { return zap.Uint8(key, val) }
func Uint8p(key string, val *uint8) zap.Field                  { return zap.Uint8p(key, val) }
func Uint8s(key string, nums []uint8) zap.Field                { return zap.Uint8s(key, nums) }
func Uintp(key string, val *uint) zap.Field                    { return zap.Uintp(key, val) }
func Uintptr(key string, val uintptr) zap.Field                { return zap.Uintptr(key, val) }
func Uintptrp(key string, val *uintptr) zap.Field              { return zap.Uintptrp(key, val) }
func Uintptrs(key string, us []uintptr) zap.Field              { return zap.Uintptrs(key, us) }
func Uints(key string, nums []uint) zap.Field                  { return zap.Uints(key, nums) }
