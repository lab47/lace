(ns lace.test-lace.hiccup
  (:require
   [lace.test :refer [deftest is testing]]
   [lace.hiccup :refer [html raw-string]]))

;; These tests are largely copied from the Hiccup source

(deftest tag-names
  (testing "basic tags"
    (is (= (html [:div]) "<div></div>"))
    (is (= (html ["div"]) "<div></div>"))
    (is (= (html ['div]) "<div></div>")))
  (testing "tag syntax sugar"
    (is (= (html [:div#foo]) "<div id=\"foo\"></div>"))
    (is (= (html [:div.foo]) "<div class=\"foo\"></div>"))
    (is (= (html [:div.foo (str "bar" "baz")])
           "<div class=\"foo\">barbaz</div>"))
    (is (= (html [:div.a.b]) "<div class=\"a b\"></div>"))
    (is (= (html [:div.a.b.c]) "<div class=\"a b c\"></div>"))
    (is (= (html [:div#foo.bar.baz])
           "<div class=\"bar baz\" id=\"foo\"></div>"))))

(deftest tag-contents
  (testing "empty tags"
    (is (= (html [:div]) "<div></div>"))
    (is (= (html [:h1]) "<h1></h1>"))
    (is (= (html [:script]) "<script></script>"))
    (is (= (html [:text]) "<text></text>"))
    (is (= (html [:a]) "<a></a>"))
    (is (= (html [:iframe]) "<iframe></iframe>"))
    (is (= (html [:title]) "<title></title>"))
    (is (= (html [:section]) "<section></section>"))
    (is (= (html [:select]) "<select></select>"))
    (is (= (html [:object]) "<object></object>"))
    (is (= (html [:video]) "<video></video>")))
  (testing "void tags"
    (is (= (html [:br]) "<br />"))
    (is (= (html [:link]) "<link />"))
    (is (= (html [:colgroup {:span 2}] "<colgroup span=\"2\" />")))
    (is (= (html [:colgroup [:col]] "<colgroup><col /></colgroup>"))))
  (testing "tags containing text"
    (is (= (html [:text "Lorem Ipsum"]) "<text>Lorem Ipsum</text>")))
  (testing "contents are concatenated"
    (is (= (html [:body "foo" "bar"]) "<body>foobar</body>"))
    (is (= (html [:body [:p] [:br]]) "<body><p></p><br /></body>")))
  (testing "seqs are expanded"
    (is (= (html [:body (list "foo" "bar")]) "<body>foobar</body>"))
    (is (= (html (list [:p "a"] [:p "b"])) "<p>a</p><p>b</p>")))
  (testing "keywords are turned into strings"
    (is (= (html [:div :foo]) "<div>foo</div>")))
  (testing "vecs don't expand - error if vec doesn't have tag name"
    (is (thrown? Error
                 (html (vector [:p "a"] [:p "b"])))))
  (testing "tags can contain tags"
    (is (= (html [:div [:p]]) "<div><p></p></div>"))
    (is (= (html [:div [:b]]) "<div><b></b></div>"))
    (is (= (html [:p [:span [:a "foo"]]])
           "<p><span><a>foo</a></span></p>"))))

(deftest tag-attributes
  (testing "tag with blank attribute map"
    (is (= (html [:xml {}]) "<xml></xml>")))
  (testing "tag with populated attribute map"
    (is (= (html [:xml {:a "1", :b "2"}]) "<xml a=\"1\" b=\"2\"></xml>"))
    (is (= (html [:img {"id" "foo"}]) "<img id=\"foo\" />"))
    (is (= (html [:img {'id "foo"}]) "<img id=\"foo\" />"))
    (is (= (html [:xml {:a "1", 'b "2", "c" "3"}])
           "<xml a=\"1\" b=\"2\" c=\"3\"></xml>")))
  ;; This logic differs from Clojure Hiccup, due to use of lace.html/escape
  #_ (testing "attribute values are escaped"
       (is (= (html [:div {:id "\""}]) "<div id=\"&quot;\"></div>")))
  (testing "attributes are escaped via lace.html/escape"
    (is (= (html [:div {:id "<_&_>_\"_'"}])
           "<div id=\"&lt;_&amp;_&gt;_&#34;_&#39;\"></div>")))
  (testing "boolean attributes"
    (is (= (html [:input {:type "checkbox" :checked true}])
           "<input checked=\"checked\" type=\"checkbox\" />"))
    (is (= (html [:input {:type "checkbox" :checked false}])
           "<input type=\"checkbox\" />")))
  (testing "nil attributes"
    (is (= (html [:span {:class nil} "foo"])
           "<span>foo</span>")))
  (testing "resolving conflicts between attributes in the map and tag"
    (is (= (html [:div.foo {:class "bar"} "baz"])
           "<div class=\"foo bar\">baz</div>"))
    (is (= (html [:div#bar.foo {:id "baq"} "baz"])
           "<div class=\"foo\" id=\"baq\">baz</div>"))))

(deftest render-modes
  (testing "closed tag"
    (is (= (html [:p] [:br]) "<p></p><br />"))
    (is (= (html {:mode :xhtml} [:p] [:br]) "<p></p><br />"))
    (is (= (html {:mode :html} [:p] [:br]) "<p></p><br>"))
    (is (= (html {:mode :xml} [:p] [:br]) "<p /><br />"))
    (is (= (html {:mode :sgml} [:p] [:br]) "<p><br>")))
  (testing "boolean attributes"
    (is (= (html {:mode :xml} [:input {:type "checkbox" :checked true}])
           "<input checked=\"checked\" type=\"checkbox\" />"))
    (is (= (html {:mode :sgml} [:input {:type "checkbox" :checked true}])
           "<input checked type=\"checkbox\">")))
  (testing "laziness and binding scope"
    (is (= (html {:mode :sgml} [:html [:link] (list [:link])])
           "<html><link><link></html>"))))

(deftest raw-and-escaped-content
  (let [text "The <blink> tag should not be used"]
    (is (= (html [:div text])
           "<div>The &lt;blink&gt; tag should not be used</div>"))
    (is (= (html [:div (raw-string text)])
           "<div>The <blink> tag should not be used</div>"))))
