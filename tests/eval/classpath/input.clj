(require 'a.b.c)

(binding [lace.core/*classpath* ["."]]
  (require 'd.e.f))

(binding [lace.core/*classpath* ["." "a"]]
  (require 'b.c))

(binding [lace.core/*classpath* ["x/y"]]
  (require 'z))
