// This file is generated by generate-std.clj script. Do not edit manually!


package math

import (
	. "github.com/lab47/lace/core"
	"fmt"
	"os"
)

func InternsOrThunks(env *Env, ns *Namespace) {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of math.InternsOrThunks().")
	}
	ns.ResetMeta(MakeMeta(nil, `Provides basic constants and mathematical functions.`, "1.0"))

	ns.InternVar(env, "e", e_,
		MakeMeta(
			nil,
			`e`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "ln-of-10", ln_of_10_,
		MakeMeta(
			nil,
			`Natural logarithm of 10`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "ln-of-2", ln_of_2_,
		MakeMeta(
			nil,
			`Natural logarithm of 2`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "log-10-of-e", log_10_of_e_,
		MakeMeta(
			nil,
			`Base-10 logarithm of e`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "log-2-of-e", log_2_of_e_,
		MakeMeta(
			nil,
			`Base-2 logarithm of e`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "max-double", max_double_,
		MakeMeta(
			nil,
			`Largest finite value representable by Double`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "phi", phi_,
		MakeMeta(
			nil,
			`Phi`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "pi", pi_,
		MakeMeta(
			nil,
			`pi`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "smallest-nonzero-double", smallest_nonzero_double_,
		MakeMeta(
			nil,
			`Smallest positive, non-zero value representable by Double`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "sqrt-of-2", sqrt_of_2_,
		MakeMeta(
			nil,
			`Square root of 2`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "sqrt-of-e", sqrt_of_e_,
		MakeMeta(
			nil,
			`Square root of e`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "sqrt-of-phi", sqrt_of_phi_,
		MakeMeta(
			nil,
			`Square root of phi`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "sqrt-of-pi", sqrt_of_pi_,
		MakeMeta(
			nil,
			`Square root of pi`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "abs", abs_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the absolute value of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "ceil", ceil_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the least integer value greater than or equal to x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "copy-sign", copy_sign_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"), MakeSymbol("y"))),
			`Returns value with the magnitude of x and the sign of y.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "cos", cos_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the cosine of the radian argument x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "cube-root", cube_root_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the cube root of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "dim", dim_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"), MakeSymbol("y"))),
			`Returns the maximum of x-y and 0.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "exp", exp_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns e**x, the base-e exponential of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "exp-2", exp_2_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns 2**x, the base-2 exponential of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "exp-minus-1", exp_minus_1_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns e**x - 1, the base-e exponential of x minus 1.

  This is more accurate than (- (exp x) 1.) when x is near zero.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "floor", floor_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the greatest integer value greater than or equal to x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "hypot", hypot_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("p"), MakeSymbol("q"))),
			`Returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "inf", inf_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("sign"))),
			`Returns positive infinity if sign >= 0, negative infinity if sign < 0.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "inf?", isinf_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"), MakeSymbol("sign"))),
			`Returns whether x is an infinity.

  If sign > 0, returns whether x is positive infinity; if < 0, whether negative infinity; if == 0, whether either infinity.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	ns.InternVar(env, "log", log_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the natural logarithm of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "log-10", log_10_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the decimal logarithm of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "log-2", log_2_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the binary logarithm of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "log-binary", log_binary_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the binary exponent of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "log-plus-1", log_plus_1_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the natural logarithm of 1 plus x.

  This is more accurate than (log (+ 1 x)) when x is near zero.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "modf", modf_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns a vector with the integer and fractional floating-point numbers that sum to x.

  Both values have the same sign as x.`, "1.0"))

	ns.InternVar(env, "nan", nan_,
		MakeMeta(
			NewListFrom(NewVectorFrom()),
			`Returns an IEEE 754 "not-a-number" value.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "nan?", isnan_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns whether x is an IEEE 754 "not-a-number" value.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	ns.InternVar(env, "next-after", next_after_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"), MakeSymbol("y"))),
			`Returns the next representable Double value after x towards y.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "pow", pow_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"), MakeSymbol("y"))),
			`Returns x**y, the base-x exponential of y.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "pow-10", pow_10_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns 10**x, the base-10 exponential of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "round", round_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the integer nearest to x, rounding half away from zero.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "round-to-even", round_to_even_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the integer nearest to x, rounding ties to the nearest even integer.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "sign-bit", sign_bit_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns whether x is negative or negative zero.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Boolean"}))

	ns.InternVar(env, "sin", sin_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the sine of the radian argument x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "sqrt", sqrt_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the square root of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

	ns.InternVar(env, "trunc", trunc_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("x"))),
			`Returns the integer value of x.`, "1.0").Plus(MakeKeyword("tag"), String{S: "Double"}))

}
