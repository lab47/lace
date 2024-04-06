(ns
  ^{:go-imports ["math"]
    :doc "Provides basic constants and mathematical functions."}
  math)

(defn ^Double sin
  "Returns the sine of the radian argument x."
  {:added "1.0"
  :go "math.Sin(x.Double().D), nil"}
  [^Number x])

(defn ^Double cos
  "Returns the cosine of the radian argument x."
  {:added "1.0"
  :go "math.Cos(x.Double().D), nil"}
  [^Number x])

(defn ^Double hypot
  "Returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow."
  {:added "1.0"
  :go "math.Hypot(p.Double().D, q.Double().D), nil"}
  [^Number p, ^Number q])

(defn ^Double abs
  "Returns the absolute value of x."
  {:added "1.0"
   :go "math.Abs(x.Double().D), nil"}
  [^Number x])

(defn ^Double ceil
  "Returns the least integer value greater than or equal to x."
  {:added "1.0"
   :go "math.Ceil(x.Double().D), nil"}
  [^Number x])

(defn ^Double cube-root
  "Returns the cube root of x."
  {:added "1.0"
   :go "math.Cbrt(x.Double().D), nil"}
  [^Number x])

(defn ^Double copy-sign
  "Returns value with the magnitude of x and the sign of y."
  {:added "1.0"
  :go "math.Copysign(x.Double().D, y.Double().D), nil"}
  [^Number x, ^Number y])

(defn ^Double dim
  "Returns the maximum of x-y and 0."
  {:added "1.0"
  :go "math.Dim(x.Double().D, y.Double().D), nil"}
  [^Number x, ^Number y])

(defn ^Double exp
  "Returns e**x, the base-e exponential of x."
  {:added "1.0"
   :go "math.Exp(x.Double().D), nil"}
  [^Number x])

(defn ^Double exp-2
  "Returns 2**x, the base-2 exponential of x."
  {:added "1.0"
   :go "math.Exp2(x.Double().D), nil"}
  [^Number x])

(defn ^Double exp-minus-1
  "Returns e**x - 1, the base-e exponential of x minus 1.

  This is more accurate than (- (exp x) 1.) when x is near zero."
  {:added "1.0"
   :go "math.Expm1(x.Double().D), nil"}
  [^Number x])

(defn ^Double floor
  "Returns the greatest integer value greater than or equal to x."
  {:added "1.0"
   :go "math.Floor(x.Double().D), nil"}
  [^Number x])

(defn ^Double inf
  "Returns positive infinity if sign >= 0, negative infinity if sign < 0."
  {:added "1.0"
   :go "math.Inf(sign), nil"}
  [^Int sign])

(defn ^Boolean inf?
  "Returns whether x is an infinity.

  If sign > 0, returns whether x is positive infinity; if < 0, whether negative infinity; if == 0, whether either infinity."
  {:added "1.0"
   :go "math.IsInf(x.Double().D, sign), nil"}
  [^Number x, ^Int sign])

(defn ^Double log
  "Returns the natural logarithm of x."
  {:added "1.0"
   :go "math.Log(x.Double().D), nil"}
  [^Number x])

(defn ^Double log-10
  "Returns the decimal logarithm of x."
  {:added "1.0"
   :go "math.Log10(x.Double().D), nil"}
  [^Number x])

(defn ^Double log-plus-1
  "Returns the natural logarithm of 1 plus x.

  This is more accurate than (log (+ 1 x)) when x is near zero."
  {:added "1.0"
   :go "math.Log1p(x.Double().D), nil"}
  [^Number x])

(defn ^Double log-2
  "Returns the binary logarithm of x."
  {:added "1.0"
   :go "math.Log2(x.Double().D), nil"}
  [^Number x])

(defn ^Double log-binary
  "Returns the binary exponent of x."
  {:added "1.0"
   :go "math.Logb(x.Double().D), nil"}
  [^Number x])

(defn modf
  "Returns a vector with the integer and fractional floating-point numbers that sum to x.

  Both values have the same sign as x."
  {:added "1.0"
   :go "modf(x.Double().D)"}
  [^Number x])

(defn ^Double nan
  "Returns an IEEE 754 \"not-a-number\" value."
  {:added "1.0"
   :go "math.NaN(), nil"}
  [])

(defn ^Boolean nan?
  "Returns whether x is an IEEE 754 \"not-a-number\" value."
  {:added "1.0"
   :go "math.IsNaN(x.Double().D), nil"}
  [^Number x])

(defn ^Double next-after
  "Returns the next representable Double value after x towards y."
  {:added "1.0"
   :go "math.Nextafter(x.Double().D, y.Double().D), nil"}
  [^Number x, ^Number y])

(defn ^Double pow
  "Returns x**y, the base-x exponential of y."
  {:added "1.0"
   :go "math.Pow(x.Double().D, y.Double().D), nil"}
  [^Number x, ^Number y])

(defn ^Double pow-10
  "Returns 10**x, the base-10 exponential of x."
  {:added "1.0"
   :go "math.Pow10(x), nil"}
  [^Int x])

(defn ^Double round
  "Returns the integer nearest to x, rounding half away from zero."
  {:added "1.0"
   :go "math.Round(x.Double().D), nil"}
  [^Number x])

(defn ^Double round-to-even
  "Returns the integer nearest to x, rounding ties to the nearest even integer."
  {:added "1.0"
   :go "math.RoundToEven(x.Double().D), nil"}
  [^Number x])

(defn ^Boolean sign-bit
  "Returns whether x is negative or negative zero."
  {:added "1.0"
   :go "math.Signbit(x.Double().D), nil"}
  [^Number x])

(defn ^Double sqrt
  "Returns the square root of x."
  {:added "1.0"
   :go "math.Sqrt(x.Double().D), nil"}
  [^Number x])

(defn ^Double trunc
  "Returns the integer value of x."
  {:added "1.0"
   :go "math.Trunc(x.Double().D), nil"}
  [^Number x])

(def
  ^{:doc "pi"
    :added "1.0"
    :tag Double
    :go "math.Pi"}
  pi)

(def
  ^{:doc "e"
    :added "1.0"
    :tag Double
    :go "math.E"}
  e)

(def
  ^{:doc "Phi"
    :added "1.0"
    :tag Double
    :go "math.Phi"}
  phi)

(def
  ^{:doc "Square root of 2"
    :added "1.0"
    :tag Double
    :go "math.Sqrt2"}
  sqrt-of-2)

(def
  ^{:doc "Square root of e"
    :added "1.0"
    :tag Double
    :go "math.SqrtE"}
  sqrt-of-e)

(def
  ^{:doc "Square root of pi"
    :added "1.0"
    :tag Double
    :go "math.SqrtPi"}
  sqrt-of-pi)

(def
  ^{:doc "Square root of phi"
    :added "1.0"
    :tag Double
    :go "math.SqrtPhi"}
  sqrt-of-phi)

(def
  ^{:doc "Natural logarithm of 2"
    :added "1.0"
    :tag Double
    :go "math.Ln2"}
  ln-of-2)

(def
  ^{:doc "Base-2 logarithm of e"
    :added "1.0"
    :tag Double
    :go "math.Log2E"}
  log-2-of-e)

(def
  ^{:doc "Natural logarithm of 10"
    :added "1.0"
    :tag Double
    :go "math.Ln10"}
  ln-of-10)

(def
  ^{:doc "Base-10 logarithm of e"
    :added "1.0"
    :tag Double
    :go "math.Log10E"}
  log-10-of-e)

(def
  ^{:doc "Natural logarithm of 2"
    :added "1.0"
    :tag Double
    :go "math.Ln2"}
  ln-of-2)

(def
  ^{:doc "Largest finite value representable by Double"
    :added "1.0"
    :tag Double
    :go "math.MaxFloat64"}
  max-double)

(def
  ^{:doc "Smallest positive, non-zero value representable by Double"
    :added "1.0"
    :tag Double
    :go "math.SmallestNonzeroFloat64"}
  smallest-nonzero-double)
