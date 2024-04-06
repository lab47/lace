(ns lace.crypto-test
  (:require [lace.test :refer [deftest is]]
            [lace.crypto]))

(deftest hmac-sha1
  (is (= "+9sdGxiqbAgyS31ktx+3Y3BpDh0="
         (lace.base64/encode-string (lace.crypto/hmac :sha1 "" ""))))
  (is (= "IpLOxioNEjv8uRFDMyVok9csrXs="
         (lace.base64/encode-string (lace.crypto/hmac :sha1 "asdf" "sdf"))))
  (is (= "93z3g3y8BWMq7Q1pjthYFmDy32U="
         (lace.base64/encode-string (lace.crypto/hmac
                                      :sha1
                                      "asdflkjhsdfiu239471204hsdkjhbf9128734yhfay49812734y1ihsfdkjhf2913874y"
                                      "sdfaskdjhfk1289374yrfh981273y4rihf18274yrfuh9r8hoojffh13iur0")))))

(deftest hmac-sha224
  (is (= "XOFPcolGYiE+J0jSprojS3QmORDO3eL1qScVJA=="
         (lace.base64/encode-string (lace.crypto/hmac :sha224 "" ""))))
  (is (= "Ql+wxqlZho2DIqRqAS0YmfUu2hu0od7CcNPIJw=="
         (lace.base64/encode-string (lace.crypto/hmac :sha224 "asdf" "sdf"))))
  (is (= "Fwg9+WAc0H14tLzuwOcyarWyJA6qHT9pTHYTDw=="
         (lace.base64/encode-string (lace.crypto/hmac
                                      :sha224
                                      "asdflkjhsdfiu239471204hsdkjhbf9128734yhfay49812734y1ihsfdkjhf2913874y"
                                      "sdfaskdjhfk1289374yrfh981273y4rihf18274yrfuh9r8hoojffh13iur0")))))

(deftest hmac-sha256
  (is (= "thNnmggU2ex3L5XXeMNfxf8Wl8STcVZTxscSFEKSxa0="
         (lace.base64/encode-string (lace.crypto/hmac :sha256 "" ""))))
  (is (= "bIe8a07h5VPL2Qm3LdTUKM/edjouFyaRH+R1Hix4dVI="
         (lace.base64/encode-string (lace.crypto/hmac :sha256 "asdf" "sdf"))))
  (is (= "rrGSbGQtL8ryX5dUrkBuxzzbchavIuGUjcRTgaVjetI="
         (lace.base64/encode-string (lace.crypto/hmac
                                      :sha256
                                      "asdflkjhsdfiu239471204hsdkjhbf9128734yhfay49812734y1ihsfdkjhf2913874y"
                                      "sdfaskdjhfk1289374yrfh981273y4rihf18274yrfuh9r8hoojffh13iur0")))))

(deftest hmac-sha384
  (is (= "bB8u6Tj60uJL2RKYR0OCyiGMdds9g+EUs9Q2d3bRTTVRKJ516CCc1LeSMChAI0rc"
         (lace.base64/encode-string (lace.crypto/hmac :sha384 "" ""))))
  (is (= "RjiHB1h1C7sauXeTDzzAzCSe1jGV4CT40ec66q74TMhe646SzlgUaRbjWVlzx6P+"
         (lace.base64/encode-string (lace.crypto/hmac :sha384 "asdf" "sdf"))))
  (is (= "4Llv2ryHDcSnGFaYAsdiD1439FUjGYsUsjZ8UgasloDQvVUGxVieZ8/svEvmzODo"
         (lace.base64/encode-string (lace.crypto/hmac
                                      :sha384
                                      "asdflkjhsdfiu239471204hsdkjhbf9128734yhfay49812734y1ihsfdkjhf2913874y"
                                      "sdfaskdjhfk1289374yrfh981273y4rihf18274yrfuh9r8hoojffh13iur0")))))


(deftest hmac-sha512
  (is (= "uTbO6Gyfh6pdPG8uhMtaQjml/lBICm7Ga3CrWx9KxnMMbFFUIbMn7B1pQC5T37Sa1zgesGezOP17DLIiRyJdRw=="
         (lace.base64/encode-string (lace.crypto/hmac :sha512 "" ""))))
  (is (= "3doRHNyXxgzniRzRNR6mggdgzysr9QX/PUYi5wWQE0QD3P8c3/VDQlMx9D2XvIoxORHQYtjH+dLzaHm00wNV9Q=="
         (lace.base64/encode-string (lace.crypto/hmac :sha512 "asdf" "sdf"))))
  (is (= "/9bBOJzjUA0ZJnm6a+JEFcGfI6cjB46j0l7o8GKTX0qYLi6dEWzxJOm+jY4jJcztNExSNlDO83j0LqxcYxhgow=="
         (lace.base64/encode-string (lace.crypto/hmac
                                      :sha512
                                      "asdflkjhsdfiu239471204hsdkjhbf9128734yhfay49812734y1ihsfdkjhf2913874y"
                                      "sdfaskdjhfk1289374yrfh981273y4rihf18274yrfuh9r8hoojffh13iur0")))))

(deftest sha256
  (is (= "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4"
         (lace.hex/encode-string (lace.crypto/sha256 "1234")))))

(deftest sha224
  (is (= "99fb2f48c6af4761f904fc85f95eb56190e5d40b1f44ec3a9c1fa319"
         (lace.hex/encode-string (lace.crypto/sha224 "1234")))))

(deftest sha384
  (is (= "504f008c8fcf8b2ed5dfcde752fc5464ab8ba064215d9c5b5fc486af3d9ab8c81b14785180d2ad7cee1ab792ad44798c"
         (lace.hex/encode-string (lace.crypto/sha384 "1234")))))

(deftest sha512
  (is (= "d404559f602eab6fd602ac7680dacbfaadd13630335e951f097af3900e9de176b6db28512f2e000b9d04fba5133e8b1c6e8df59db3a8ab9d60be4b97cc9e81db"
         (lace.hex/encode-string (lace.crypto/sha512 "1234")))))

(deftest sha512-224
  (is (= "90ce2868ea834f4a9e0008d504bfdcd144cf7de5913b97ecbae7c58c"
         (lace.hex/encode-string (lace.crypto/sha512-224 "1234")))))

(deftest sha512-256
  (is (= "3bb20e367b1a988158e3dbbed06bf6415beea23f56bbb15dbf346505edc7e7df"
         (lace.hex/encode-string (lace.crypto/sha512-256 "1234")))))

(deftest md5
  (is (= "81dc9bdb52d04dc20036dbd8313ed055"
         (lace.hex/encode-string (lace.crypto/md5 "1234")))))

(deftest sha1
  (is (= "7110eda4d09e062aa5e4a390b0a572ac0d2c0220"
         (lace.hex/encode-string (lace.crypto/sha1 "1234")))))
