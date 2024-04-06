(create-ns 'test.test)
(require 'lace.os)
(alias 'os 'lace.os)
(alias 't 'test.test)

(os/sh "ls")
(test.test/foo)
(t/bar)
